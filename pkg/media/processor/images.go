package processor

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math"
	"os"

	"github.com/momentum-xyz/ubercontroller/types"

	"github.com/nfnt/resize"
	_ "golang.org/x/image/webp"
)

func (p *Processor) WriteToF(img image.Image) (error, string) {
	var w bytes.Buffer
	err := png.Encode(&w, img)
	if err != nil {
		return err, ""
	}
	body := w.Bytes()
	ID := p.GetMD5HashByte(body)
	return p.SaveWriteToFile(p.ImPathF+ID, body), ID
}

func (p *Processor) WriteToScaled(base string, img image.Image, rsize string) error {
	if size, ok := types.Tsizes[rsize]; ok {
		fpath := p.ImPathS[rsize] + base
		if p.FileExists(fpath) {
			p.log.Debugf("Scaled file %s already exist, skip", fpath)
			return nil
		}
		return p.SaveWriteToPNG(fpath, DownSampleTo(img, size))
	}
	return errors.New("Not such size defined in the size map")

}

func (p *Processor) SaveWriteToPNG(fname string, img image.Image) error {
	tfilename := fname + ".tmp"

	w, err := os.Create(tfilename)
	if err != nil {
		return err
	}

	err = png.Encode(w, img)
	w.Close()
	if err != nil {
		os.Remove(tfilename)
		return err
	}
	if err = os.Rename(tfilename, fname); err != nil {
		os.Remove(tfilename)
		return err
	}
	return nil
}

func (p *Processor) SaveWriteToFile(fname string, data []byte) error {
	tfilename := fname + ".tmp"
	err := os.WriteFile(tfilename, data, 0666)
	if err != nil {
		os.Remove(tfilename)
		return err
	}
	if err = os.Rename(tfilename, fname); err != nil {
		return err
	}
	return nil
}

func (p *Processor) ProcessImage(src []byte) (string, error) {
	img, format, err := image.Decode(bytes.NewReader(src))
	if err != nil {
		return "", err
	}
	p.log.Info("Incoming image:", format)

	var ID string
	if format == "gif" {
		ID = p.GetMD5HashByte(src)
		err = p.SaveWriteToFile(p.ImPathF+ID, src)
	} else {
		err, ID = p.WriteToF(img)
	}
	if err != nil {
		return "", err
	}

	for _, v := range types.Tprecalcs {
		if err = p.WriteToScaled(ID, img, v); err != nil {
			return "", err
		}
	}

	p.log.Info("Hash:", ID)
	return ID, err

}

func DownSampleTo(img image.Image, NewPixelCount int) image.Image {
	ox := float64(img.Bounds().Max.X)
	oy := float64(img.Bounds().Max.Y)

	scl := math.Sqrt(float64(NewPixelCount) / (ox * oy))

	nx := uint(math.Round(ox * scl))
	ny := uint(math.Round(oy * scl))

	imgout := resize.Thumbnail(nx, ny, img, resize.Bilinear)
	return imgout
}
