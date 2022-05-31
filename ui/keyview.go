package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type KeyView struct {
	widget.BaseWidget
	filenpath string
	button    widget.Button
	view      widget.Entry
}

func (c *KeyView) Layout(size fyne.Size) {

}

func (c *KeyView) MinSize() {

}

func (c *KeyView) Refresh() {

}

func (c *KeyView) CreateRenderer() fyne.WidgetRenderer {
	return nil
}
