package widgets

import (
	"image"
	"log"

	"github.com/gizak/termui/v3"
)

// Sparkline is like: ▅▆▂▂▅▇▂▂▃▆▆▆▅▃. The data points should be non-negative integers.
type Sparkline struct {
	Data       []int
	Title1     string
	Title2     string
	TitleColor termui.Color
	LineColor  termui.Color
}

// SparklineGroup is a renderable widget which groups together the given sparklines.
type SparklineGroup struct {
	*termui.Block
	Lines []*Sparkline
}

// Add appends a given Sparkline to the *SparklineGroup.
func (self *SparklineGroup) Add(sl Sparkline) {
	self.Lines = append(self.Lines, &sl)
}

// NewSparkline returns an unrenderable single sparkline that intended to be added into a SparklineGroup.
func NewSparkline() *Sparkline {
	return &Sparkline{}
}

// NewSparklineGroup return a new *SparklineGroup with given Sparklines, you can always add a new Sparkline later.
func NewSparklineGroup(ss ...*Sparkline) *SparklineGroup {
	return &SparklineGroup{
		Block: termui.NewBlock(),
		Lines: ss,
	}
}

func (self *SparklineGroup) Draw(buf *termui.Buffer) {
	self.Block.Draw(buf)

	lc := len(self.Lines) // lineCount

	// renders each sparkline and its titles
	for i, line := range self.Lines {

		// prints titles
		title1Y := self.Inner.Min.Y + 1 + (self.Inner.Dy()/lc)*i
		title2Y := self.Inner.Min.Y + 2 + (self.Inner.Dy()/lc)*i
		title1 := termui.TrimString(line.Title1, self.Inner.Dx())
		title2 := termui.TrimString(line.Title2, self.Inner.Dx())
		if self.Inner.Dy() > 5 {
			buf.SetString(
				title1,
				termui.NewStyle(line.TitleColor, termui.ColorClear, termui.ModifierBold),
				image.Pt(self.Inner.Min.X, title1Y),
			)
		}
		if self.Inner.Dy() > 6 {
			buf.SetString(
				title2,
				termui.NewStyle(line.TitleColor, termui.ColorClear, termui.ModifierBold),
				image.Pt(self.Inner.Min.X, title2Y),
			)
		}

		sparkY := (self.Inner.Dy() / lc) * (i + 1)
		// finds max data in current view used for relative heights
		max := 1
		for i := len(line.Data) - 1; i >= 0 && self.Inner.Dx()-((len(line.Data)-1)-i) >= 1; i-- {
			if line.Data[i] > max {
				max = line.Data[i]
			}
		}
		// prints sparkline
		for x := self.Inner.Dx(); x >= 1; x-- {
			char := termui.BARS[1]
			if (self.Inner.Dx() - x) < len(line.Data) {
				offset := self.Inner.Dx() - x
				curItem := line.Data[(len(line.Data)-1)-offset]
				percent := float64(curItem) / float64(max)
				index := int(percent*float64(len(termui.BARS)-2)) + 1
				if index < 1 || index >= len(termui.BARS) {
					log.Printf(
						"invalid sparkline data value. len(data) %d, offset %d, curItem %d, percent %f, index %d",
						len(line.Data), offset, curItem, percent, index,
					)
				} else {
					char = termui.BARS[index]
				}
			}
			buf.SetCell(
				termui.NewCell(char, termui.NewStyle(line.LineColor)),
				image.Pt(self.Inner.Min.X+x-1, self.Inner.Min.Y+sparkY-1),
			)
		}
		dx := self.Inner.Dx()
		if len(line.Data) > 4*dx {
			line.Data = line.Data[dx-1:]
		}
	}
}