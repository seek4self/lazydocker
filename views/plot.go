package views

import (
	"fmt"
	"image"
	"math"
	"strconv"

	ui "github.com/gizak/termui/v3"
)

const (
	xAxisLabelsHeight = 1
	yAxisLabelsWidth  = 4
	xAxisLabelsGap    = 2
	yAxisLabelsGap    = 1
)

type PlotType uint

const (
	LineChart PlotType = iota
	ScatterPlot
)

type PlotMarker uint

const (
	MarkerBraille PlotMarker = iota
	MarkerDot
)

type DrawDirection uint

const (
	DrawLeft DrawDirection = iota
	DrawRight
)

type Plot struct {
	*ui.Block

	Data       [][]float64
	DataLabels []string
	MaxVal     float64

	LineColors []ui.Color
	AxesStyle  ui.Style // TODO
	ShowAxes   bool

	Marker          PlotMarker
	DotMarkerRune   rune
	PlotType        PlotType
	HorizontalScale int
	DrawDirection   DrawDirection // TODO

	XLabelGap, YLabelGap int

	ox, oy int // Coordinate origin
}

func NewPlot() *Plot {
	return &Plot{
		Block:           ui.NewBlock(),
		LineColors:      ui.Theme.Plot.Lines,
		AxesStyle:       ui.NewStyle(ui.Theme.Plot.Axes),
		Marker:          MarkerBraille,
		DotMarkerRune:   ui.DOT,
		Data:            make([][]float64, 0),
		HorizontalScale: 1,
		DrawDirection:   DrawRight,
		ShowAxes:        true,
		PlotType:        LineChart,

		XLabelGap: 5,
		YLabelGap: 3,
	}
}

func (p *Plot) renderBraille(buf *ui.Buffer, drawArea image.Rectangle) {
	canvas := ui.NewCanvas()
	canvas.Rectangle = drawArea

	switch p.PlotType {
	case ScatterPlot:
		for i, line := range p.Data {
			for j, val := range line {
				height := int((val / p.MaxVal) * float64(drawArea.Dy()-1))
				canvas.SetPoint(
					image.Pt(
						(drawArea.Min.X+(j*p.HorizontalScale))*2,
						(drawArea.Max.Y-height-1)*4,
					),
					ui.SelectColor(p.LineColors, i),
				)
			}
		}
	case LineChart:
		for i, line := range p.Data {
			previousHeight := int((line[1] / p.MaxVal) * float64(drawArea.Dy()-1))
			for j, val := range line[1:] {
				height := int((val / p.MaxVal) * float64(drawArea.Dy()-1))
				canvas.SetLine(
					image.Pt(
						(drawArea.Min.X+(j*p.HorizontalScale))*2,
						(drawArea.Max.Y-previousHeight-1)*4,
					),
					image.Pt(
						(drawArea.Min.X+((j+1)*p.HorizontalScale))*2,
						(drawArea.Max.Y-height-1)*4,
					),
					ui.SelectColor(p.LineColors, i),
				)
				previousHeight = height
			}
		}
	}

	canvas.Draw(buf)
}

func (p *Plot) renderDot(buf *ui.Buffer, drawArea image.Rectangle) {
	switch p.PlotType {
	case ScatterPlot:
		for i, line := range p.Data {
			for j, val := range line {
				height := int((val / p.MaxVal) * float64(drawArea.Dy()-1))
				point := image.Pt(drawArea.Min.X+(j*p.HorizontalScale), drawArea.Max.Y-1-height)
				if point.In(drawArea) {
					buf.SetCell(
						ui.NewCell(p.DotMarkerRune, ui.NewStyle(ui.SelectColor(p.LineColors, i))),
						point,
					)
				}
			}
		}
	case LineChart:
		for i, line := range p.Data {
			for j := 0; j < len(line) && j*p.HorizontalScale < drawArea.Dx(); j++ {
				val := line[j]
				height := int((val / p.MaxVal) * float64(drawArea.Dy()-1))
				buf.SetCell(
					ui.NewCell(p.DotMarkerRune, ui.NewStyle(ui.SelectColor(p.LineColors, i))),
					image.Pt(drawArea.Min.X+(j*p.HorizontalScale), drawArea.Max.Y-1-height),
				)
			}
		}
	}
}

// draw the axes
func (p *Plot) plotAxes(buf *ui.Buffer) {
	// draw origin cell
	p.drawAxisOrigin(buf)
	// draw x axis line
	p.drawXAxis(buf)
	// draw y axis line
	p.drawYAxis(buf)
	// draw x axis labels
	// for x := p.ox + (xAxisLabelsGap)*p.HorizontalScale + 1; x < p.Inner.Max.X-1; {
	// 	label := fmt.Sprintf(
	// 		"%d",
	// 		(x-(p.ox)-1)/(p.HorizontalScale)+1,
	// 	)
	// 	buf.SetString(
	// 		label,
	// 		ui.NewStyle(ui.ColorWhite),
	// 		image.Pt(x, p.Inner.Max.Y-1),
	// 	)
	// 	x += (len(label) + xAxisLabelsGap) * p.HorizontalScale
	// }
	// draw y axis labels
	// verticalScale := p.MaxVal / float64(p.Inner.Dy()-xAxisLabelsHeight-1)
	// for i := 0; i*(yAxisLabelsGap+1) < p.Inner.Dy()-1; i++ {
	// 	buf.SetString(
	// 		fmt.Sprintf("%.2f", float64(i)*verticalScale*(yAxisLabelsGap+1)),
	// 		ui.NewStyle(ui.ColorWhite),
	// 		image.Pt(p.Inner.Min.X, p.Inner.Max.Y-(i*(yAxisLabelsGap+1))-2),
	// 	)
	// }
}

func (p *Plot) drawAxisOrigin(buf *ui.Buffer) {
	p.ox = p.Inner.Min.X + yAxisLabelsWidth
	p.oy = p.Inner.Max.Y - xAxisLabelsHeight - 1
	buf.SetCell(
		ui.NewCell(ui.BOTTOM_LEFT, p.AxesStyle),
		image.Pt(p.ox, p.oy),
	)
	buf.SetCell(
		ui.NewCell('0', p.AxesStyle),
		image.Pt(p.ox-1, p.Inner.Max.Y-1),
	)
}

func (p *Plot) xLen() int {
	return len(p.Data[0])
}

func (p *Plot) xOffset(n int) int {
	xLen := p.xLen()
	if p.xLen() > p.Inner.Dx()-yAxisLabelsWidth {
		xLen = p.Inner.Dx() - yAxisLabelsWidth - 3
	}
	x := int(math.Round(float64(xLen) / float64(p.XLabelGap) * float64(n)))
	if x > p.xLen() {
		return p.xLen()
	}
	return x
}

func (p *Plot) xLabels() int {
	return int(math.Round(float64(p.xLen()) / float64(p.XLabelGap)))
}

func (p *Plot) drawXAxis(buf *ui.Buffer) {
	xAxis := ui.NewCell(ui.HORIZONTAL_DASH, p.AxesStyle)
	for i := yAxisLabelsWidth + 1; i < p.Inner.Dx(); i++ {
		buf.SetCell(xAxis, image.Pt(i+p.Inner.Min.X, p.oy))
	}
	// num := p.xLabels()
	for i := 1; i*p.XLabelGap+p.ox < p.Inner.Max.X-1; i++ {
		offset := i * p.XLabelGap
		label := strconv.Itoa(offset * p.HorizontalScale)
		buf.SetString(
			label,
			p.AxesStyle,
			image.Pt(p.ox+offset, p.Inner.Max.Y-1),
		)
	}
}

func (p *Plot) drawYAxis(buf *ui.Buffer) {
	yAxis := ui.NewCell(ui.VERTICAL_DASH, p.AxesStyle)
	for i := 0; i < p.Inner.Dy()-xAxisLabelsHeight-1; i++ {
		buf.SetCell(yAxis, image.Pt(p.ox, i+p.Inner.Min.Y))
	}
	verticalScale := p.MaxVal / float64(p.Inner.Dy()-xAxisLabelsHeight-1)
	for i := 1; i*p.YLabelGap < p.Inner.Dy()-1; i++ {
		label := fmt.Sprintf("%.1f", float64(i*p.YLabelGap)*verticalScale)
		buf.SetString(
			label,
			p.AxesStyle,
			image.Pt(p.Inner.Min.X+(yAxisLabelsWidth-len(label)), p.Inner.Max.Y-i*p.YLabelGap-2),
		)
	}
}

func (p *Plot) Draw(buf *ui.Buffer) {
	p.Block.Draw(buf)

	if p.MaxVal == 0 {
		p.MaxVal, _ = ui.GetMaxFloat64From2dSlice(p.Data)
	}

	if p.ShowAxes {
		p.plotAxes(buf)
	}

	drawArea := p.Inner
	if p.ShowAxes {
		drawArea = image.Rect(
			p.ox+1, p.Inner.Min.Y,
			p.Inner.Max.X, p.oy,
		)
	}

	switch p.Marker {
	case MarkerBraille:
		p.renderBraille(buf, drawArea)
	case MarkerDot:
		p.renderDot(buf, drawArea)
	}
}
