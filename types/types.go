package types

import "sync"

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
	Wg    sync.WaitGroup
}

type MetaDef struct {
	H, W int
	Mime string
}

var Tsizes = map[string]int{
	"s0": 1024,
	"s1": 4096,
	"s2": 9216,
	"s3": 25600,
	"s4": 65536,
	"s5": 193600,
	"s6": 577600,
	"s7": 1721344,
	"s8": 5062500,
	"s9": 14745600,
}

var Tprecalcs = [...]string{"s2", "s3", "s4", "s5", "s6"}
