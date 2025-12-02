package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func showError(err error, parent fyne.Window) {
	dialog.ShowError(err, parent)
}

func showInfo(message string, parent fyne.Window) {
	dialog.ShowInformation("Информация", message, parent)
}

func createTableWidget(data [][]interface{}, columns []string) *widget.Table {
	if len(columns) == 0 {
		return widget.NewTable(
			func() (int, int) { return 0, 0 },
			func() fyne.CanvasObject { return widget.NewLabel("") },
			func(id widget.TableCellID, obj fyne.CanvasObject) {},
		)
	}

	rowCount := len(data) + 1
	colCount := len(columns)

	headerStyle := fyne.TextStyle{Bold: true}
	cellStyle := fyne.TextStyle{}

	table := widget.NewTable(
		func() (int, int) {
			return rowCount, colCount
		},
		func() fyne.CanvasObject {
			lbl := widget.NewLabelWithStyle("", fyne.TextAlignLeading, cellStyle)
			lbl.Truncation = fyne.TextTruncateEllipsis
			lbl.Wrapping = fyne.TextWrapOff
			return lbl
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id.Row == 0 {
				label.TextStyle = headerStyle
				if id.Col < len(columns) {
					label.SetText(columns[id.Col])
				} else {
					label.SetText("")
				}
				return
			}

			label.TextStyle = cellStyle
			rowIdx := id.Row - 1
			if rowIdx < len(data) && id.Col < len(data[rowIdx]) {
				val := data[rowIdx][id.Col]
				if val != nil {
					label.SetText(fmt.Sprintf("%v", val))
				} else {
					label.SetText("")
				}
			} else {
				label.SetText("")
			}
		},
	)

	// Set row heights: header taller than data rows
	table.SetRowHeight(0, 32)
	for r := 1; r < rowCount; r++ {
		table.SetRowHeight(r, 28)
	}

	const (
		minColumnWidth float32 = 90
		maxColumnWidth float32 = 420
		charPixelWidth float32 = 8
		extraPadding   float32 = 30
	)

	for colIdx := 0; colIdx < colCount; colIdx++ {
		maxLen := len([]rune(columns[colIdx]))
		for _, row := range data {
			if colIdx < len(row) {
				strVal := fmt.Sprintf("%v", row[colIdx])
				if l := len([]rune(strVal)); l > maxLen {
					maxLen = l
				}
			}
		}

		width := float32(maxLen)*charPixelWidth + extraPadding
		if width < minColumnWidth {
			width = minColumnWidth
		}
		if width > maxColumnWidth {
			width = maxColumnWidth
		}
		table.SetColumnWidth(colIdx, width)
	}

	return table
}

const defaultFormFieldWidth float32 = 280

func formField(obj fyne.CanvasObject) fyne.CanvasObject {
	min := obj.MinSize()
	width := defaultFormFieldWidth
	if min.Width > width {
		width = min.Width
	}
	height := min.Height
	if height < 36 {
		height = 36
	}
	return container.New(layout.NewGridWrapLayout(fyne.NewSize(width, height)), obj)
}

func showFormDialog(title string, form *widget.Form, parent fyne.Window) {
	form.SubmitText = "Сохранить"
	content := container.NewPadded(form)
	dialog := dialog.NewCustom(title, "Отмена", content, parent)
	min := content.MinSize()
	if min.Width < 520 {
		min.Width = 520
	}
	if min.Height < 360 {
		min.Height = 360
	}
	dialog.Resize(min)
	dialog.Show()
}

func createCRUDButtons(onCreate, onUpdate, onDelete func(), parent fyne.Window) fyne.CanvasObject {
	createBtn := widget.NewButton("Создать", func() {
		if onCreate != nil {
			onCreate()
		}
	})
	updateBtn := widget.NewButton("Обновить", func() {
		if onUpdate != nil {
			onUpdate()
		}
	})
	deleteBtn := widget.NewButton("Удалить", func() {
		if onDelete != nil {
			onDelete()
		}
	})
	exportBtn := widget.NewButton("Экспорт в Excel", func() {
		// Will be implemented in each tab
	})

	return container.NewHBox(createBtn, updateBtn, deleteBtn, exportBtn)
}
