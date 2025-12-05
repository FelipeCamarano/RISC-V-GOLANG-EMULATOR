package components

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type SpeedControl struct {
	Container fyne.CanvasObject
	Value     binding.Float
}

func NewSpeedControl() *SpeedControl {
	val := binding.NewFloat()
	val.Set(50.0)

	label := widget.NewLabel("Clock Speed")
	
	slider := widget.NewSliderWithData(1, 100, val)
	slider.Step = 1

	lblValue := widget.NewLabelWithData(binding.FloatToStringWithFormat(val, "%.0f%%"))
	lblValue.TextStyle = fyne.TextStyle{Monospace: true}

	content := container.NewBorder(nil, nil, label, lblValue, slider)
	
	return &SpeedControl{
		Value:     val,
		Container: content,
	}
}

func (sc *SpeedControl) GetDelay() time.Duration {
	v, _ := sc.Value.Get()
	if v >= 100 {
		return 0
	}
	return time.Duration(100-v) * time.Millisecond
}