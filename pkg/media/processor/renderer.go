package processor

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/disintegration/gift"
	"github.com/getsentry/sentry-go"
	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/render/mod"
)

type FrameDesc struct {
	Background []uint32     `json:"background"`
	BGimage    string       `json:"bgimage"`
	Color      []uint32     `json:"color"`
	Thickness  int          `json:"thickness"`
	Width      int          `json:"width"`
	Height     int          `json:"height"`
	X          int          `json:"x"`
	Y          int          `json:"y"`
	Text       *TextDesc    `json:"text"`
	Sub        []*FrameDesc `json:"sub"`
}

type TextDesc struct {
	String    string   `json:"string"`
	Fontname  string   `json:"fontfile"`
	Fontsize  float64  `json:"fontsize"`
	Fontcolor []uint32 `json:"fontcolor"`
	Wrap      bool     `json:"wrap"`
	PadX      int      `json:"padX"`
	PadY      int      `json:"padY"`
	AlignH    string   `json:"alignH"`
	AlignV    string   `json:"alignV"`
	DPI       float64  `json:"dpi"`
}

type FrameRenderRequest struct {
	ID    *string
	Frame *FrameDesc
	wg    sync.WaitGroup
}

const (
	defaultTextDPI    = 100
	oakDefaultTextDPI = 72
)

var libgdmutex sync.Mutex

// 	libgdmutex.Lock()
// 	ret := img.StringFT(fg, fontname, ptsize, angle, x, y, str)
// 	libgdmutex.Unlock()
// 	return ret

// }

func needrender(ID *string) bool {
	return true
}

func (p *Processor) getColor(components []uint32) color.Color {
	if components == nil {
		p.log.Warn("No color")
		return color.Black
	}
	lc := len(components)
	if lc < 3 {
		sentry.CaptureMessage("Wrong color elements")
		p.log.Warn("Wrong color: ", components)
		return color.Black
	}

	if lc == 4 {
		return color.RGBA{R: uint8(components[0]), G: uint8(components[1]), B: uint8(components[2]), A: uint8(components[3])}
	}
	return color.RGBA{R: uint8(components[0]), G: uint8(components[1]), B: uint8(components[2]), A: 255}
}

func (p *Processor) RenderFrame(req *FrameRenderRequest) {
	if req.Frame == nil {
		return
	}

	img := render.NewCompositeM()
	img.Append(render.NewColorBoxM(req.Frame.Width, req.Frame.Height, color.Transparent))

	p.renderSubFrame(req.Frame, req.Frame.X, req.Frame.Y, img)

	fname := p.ImPathF + *req.ID
	tfilename := fname + ".tmp"

	f, err := os.Create(fname)
	if !p.HandleError(err) {
		defer f.Close()
		res := img.ToSprite().Modify(mod.CropToSize(req.Frame.Width, req.Frame.Height, gift.TopLeftAnchor))
		p.HandleError(png.Encode(f, res.GetRGBA()))
	}

	// if fileExists(fname) {
	// 	os.Remove(fname)
	// }
	os.Rename(tfilename, fname)

	if reader, err := os.Open(fname); !p.HandleError(err) {
		if img, _, err := image.Decode(reader); !p.HandleError(err) {
			for _, v := range Tprecalcs {
				p.WriteToScaled(*req.ID, img, v)
			}
		}
	}

	p.RenderDone <- req
	p.log.Debug("Render Done")
	req.wg.Done()
}

func (p *Processor) renderSubFrame(frame *FrameDesc, xul int, yul int, img *render.CompositeM) {
	if frame == nil {
		return
	}
	p.log.Debug(frame)
	img.AppendOffset(
		render.NewColorBoxM(frame.Width, frame.Height, p.getColor(frame.Background)),
		floatgeom.Point2{float64(xul), float64(yul)},
	)
	if frame.BGimage != "" {
		bgimpath := p.ImPathF + frame.BGimage
		p.log.Debug("bgimage path:", bgimpath)
		if p.FileExists(bgimpath) {
			p.log.Debug("bgimage exists")
			img1, err := render.LoadSprite(bgimpath)
			p.HandleError(err)
			p.log.Debug("sprite loaded")
			img1.Modify(mod.Resize(frame.Width, frame.Height, gift.NearestNeighborResampling))
			img.AppendOffset(img1, floatgeom.Point2{float64(xul), float64(yul)})
		}
	}
	p.log.Debug("subframe5:", "xul:", xul, "yul:", yul)
	if frame.Thickness > 0 {
		thc := float64(frame.Thickness) / 2
		col := p.getColor(frame.Color)
		fw := float64(frame.Width)
		fh := float64(frame.Height)
		img.AppendOffset(
			render.NewThickLine(thc, thc, fw-thc, thc, col, int(thc)),
			floatgeom.Point2{float64(xul), float64(yul)},
		)
		img.AppendOffset(
			render.NewThickLine(thc, thc, thc, fh-thc, col, int(thc)),
			floatgeom.Point2{float64(xul) - 2*thc + fw, float64(yul)},
		)
		img.AppendOffset(
			render.NewThickLine(thc, thc, fw-thc, thc, col, int(thc)),
			floatgeom.Point2{float64(xul), float64(yul) - 2*thc + fh},
		)
		img.AppendOffset(
			render.NewThickLine(thc, thc, thc, fh-thc, col, int(thc)),
			floatgeom.Point2{float64(xul), float64(yul)},
		)
	}
	p.log.Debug(frame)
	if frame.Text != nil {
		if frame.Text.DPI == 0 {
			frame.Text.DPI = defaultTextDPI
		}
		p.drawText(frame, xul, yul, img)
	}

	for _, sf := range frame.Sub {
		p.log.Debug(sf)
		p.renderSubFrame(sf, xul+sf.X, yul+sf.Y, img)
	}
}

func (p *Processor) drawText(frame *FrameDesc, xbase, ybase int, img *render.CompositeM) {
	if frame.Text == nil || frame.Text.String == "" {
		return
	}
	p.log.Debug("xbase:", xbase, "ybase:", ybase)

	var fontName string
	if frame.Text.Fontname == "" {
		fontName = "IBMPlexSans-Bold" + ".ttf"
	} else {
		fontName = frame.Text.Fontname + ".ttf"
	}

	clr := image.NewUniform(p.getColor(frame.Text.Fontcolor))
	xmax := frame.Width - frame.Text.PadX*2 - 2*frame.Thickness
	ymax := frame.Height - frame.Text.PadY*2 - 2*frame.Thickness
	fontFile := filepath.Join(p.Fontpath, fontName)

	libgdmutex.Lock()
	defer libgdmutex.Unlock()

	drawstring := frame.Text.String
	p.log.Debug("drawstring:", drawstring)

	var tw float64
	var th float64
	var fnt *render.Font
	var txts []*render.Text
	if frame.Text.Wrap {
		var err error
		fnt, err = newFont(fontFile, clr, render.FontOptions{
			Size: frame.Text.Fontsize,
			DPI:  frame.Text.DPI,
		})
		p.HandleError(err)

		for _, s := range strings.Split(drawstring, "\n") {
			lines := p.getLines(strings.Split(s, " "), xmax, fontFile, render.FontOptions{
				Size: frame.Text.Fontsize,
				DPI:  frame.Text.DPI,
			})
			p.log.Debug("lines:", len(lines))
			for _, l := range lines {
				txts = append(txts, fnt.NewText(l, 0, 0))
			}
		}

		for _, t := range txts {
			sr := t.ToSprite().GetRGBA().Rect
			if float64(sr.Max.X) > tw {
				tw = float64(sr.Max.X)
			}
			th += float64(sr.Max.Y)
		}
	} else {
		drawstring = SanitizeString(drawstring)
		opts := render.FontOptions{
			Size: frame.Text.Fontsize,
			DPI:  frame.Text.DPI,
		}
		opts = p.genFontOptions(drawstring, fontFile, xmax, ymax, opts)

		var err error
		fnt, err = newFont(fontFile, clr, opts)
		p.HandleError(err)

		txt := fnt.NewText(drawstring, 0, 0)
		txts = append(txts, txt)

		fb, fa := fnt.BoundString(drawstring)
		p.log.Debug("fb:", fb, "fa:", fa)
		tw = float64(fa.Floor() + fb.Min.X.Floor())
		th = float64(fb.Max.Y.Floor() - fb.Min.Y.Floor())
	}
	p.log.Debug("tw:", tw, "th:", th)

	var xb float64
	switch frame.Text.AlignH {
	case "right":
		xb = float64(xbase) + float64(frame.Width) - float64(frame.Text.PadX) - float64(frame.Thickness) - tw
	case "left":
		xb = float64(xbase) + float64(frame.Text.PadX) + float64(frame.Thickness)
	default:
		xb = float64(xbase) + float64(frame.Width)/2 - tw/2
	}

	//// default: alignV = top
	var yb float64
	switch frame.Text.AlignV {
	case "center":
		yb = float64(ybase) + float64(frame.Height)/2 - th/2
	case "bottom":
		yb = float64(ybase) + float64(frame.Height) - float64(frame.Text.PadY) - float64(frame.Thickness) - th
	default:
		yb = float64(ybase) + float64(frame.Text.PadY) + float64(frame.Thickness)
	}

	p.log.Debug("xb:", xb, "yb:", yb)

	for i, txt := range txts {
		ts := txt.ToSprite()
		img.AppendOffset(ts, floatgeom.Point2{xb, yb + float64(ts.GetRGBA().Rect.Max.Y*i)})
	}
}

func (p *Processor) genFontOptions(str, fontFile string, w, h int, opts render.FontOptions) render.FontOptions {
	if opts.Size != 0 {
		return opts
	}

	fbeg := 0.0
	fend := float64(math.Min(float64(w), float64(h))) + 100
	fcur := (fend - fbeg) / 2
	for {
		fnt, err := newFont(fontFile, image.Black, render.FontOptions{
			Size: fcur,
			DPI:  opts.DPI,
		})
		p.HandleError(err)

		fb, fa := fnt.BoundString(str)
		p.log.Debug("fb:", fb, "fa:", fa)
		fw := fa.Floor() - fb.Min.X.Floor()
		fh := fb.Max.Y.Floor() - fb.Min.Y.Floor()

		fd := fw
		max := w
		if w-fw > h-fh {
			fd = fh
			max = h
		}
		p.log.Debug("fbeg", fbeg, "fend:", fend, "fcur:", fcur, "fd:", fd, "max:", max)

		if math.Abs(float64(fd)-float64(max)) <= 1 {
			return render.FontOptions{
				Size: fcur,
				DPI:  opts.DPI,
			}
		}

		if fd < max {
			fbeg = fcur
			fcur += (fend - fbeg) / 2
		} else {
			fend = fcur
			fcur -= (fend - fbeg) / 2
		}
	}
}

func (p *Processor) getLines(strs []string, wmax int, fontFile string, options render.FontOptions) []string {
	fnt, err := newFont(fontFile, image.Black, options)
	p.HandleError(err)

	var n int
	var lines []string
	for n < len(strs) {
		line, end := p.getLine(strs[n:], wmax, fnt)
		if end == 0 {
			break
		}
		lines = append(lines, line)
		n += end
	}
	return lines
}

func (p *Processor) getLine(strs []string, wmax int, fnt *render.Font) (string, int) {
	if len(strs) == 0 {
		return "", 0
	}
	if len(strs) == 1 {
		return strs[0], 1
	}

	for i := 0; i < len(strs)-1; i++ {
		str := strings.Join(strs[:i+1], " ")
		if fnt.MeasureString(str).Ceil() > wmax {
			return strings.Join(strs[:i], " "), i
		}
	}
	return strings.Join(strs, " "), len(strs)
}

func newFont(file string, clr image.Image, options render.FontOptions) (*render.Font, error) {
	gen := render.FontGenerator{
		File:        file,
		Color:       clr,
		FontOptions: options,
	}
	return gen.Generate()
}

func newText(str string, file string, clr image.Image, options render.FontOptions) (*render.Text, error) {
	fnt, err := newFont(file, clr, options)
	if err != nil {
		return nil, err
	}
	return fnt.NewText(str, 0, 0), nil
}

func SanitizeString(drawstring string) string {
	newstr := drawstring[:1]
	lastchar := rune(drawstring[0])
	for _, char := range drawstring[1:] {
		if !((char == ' ' && lastchar == ' ') || (lastchar == '\n' && char == ' ')) {
			if char == '\n' && lastchar == ' ' {
				newstr = newstr[:len(newstr)-1] + "\n"
				lastchar = '\n'
			} else {
				newstr = newstr + string(char)
				lastchar = char
			}
		}
	}
	return newstr
}
