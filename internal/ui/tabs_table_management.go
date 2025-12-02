package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (mw *MainWindow) createTableManagementTab() *container.TabItem {
	tablesList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			label.SetText("")
		},
	)

	refreshTablesList := func() {
		tables, err := mw.utils.TableManager.GetTables()
		if err != nil {
			showError(err, mw.window)
			return
		}

		tablesList = widget.NewList(
			func() int { return len(tables) },
			func() fyne.CanvasObject {
				return widget.NewLabel("")
			},
			func(id widget.ListItemID, obj fyne.CanvasObject) {
				label := obj.(*widget.Label)
				if id < len(tables) {
					label.SetText(tables[id])
				}
			},
		)
	}

	refreshTablesList()

	createTableBtn := widget.NewButton("Создать таблицу", func() {
		mw.showCreateTableForm(refreshTablesList)
	})

	var selectedTableIdx int = -1
	tablesList.OnSelected = func(id widget.ListItemID) {
		selectedTableIdx = int(id)
	}

	deleteTableBtn := widget.NewButton("Удалить таблицу", func() {
		if selectedTableIdx < 0 {
			showInfo("Выберите таблицу", mw.window)
			return
		}
		tables, _ := mw.utils.TableManager.GetTables()
		if selectedTableIdx < len(tables) {
			tableName := tables[selectedTableIdx]
			dialog.ShowConfirm("Подтверждение", fmt.Sprintf("Удалить таблицу %s?", tableName), func(ok bool) {
				if ok {
					if err := mw.utils.TableManager.DropTable(tableName); err != nil {
						showError(err, mw.window)
					} else {
						showInfo("Таблица удалена", mw.window)
						refreshTablesList()
					}
				}
			}, mw.window)
		}
	})

	addColumnBtn := widget.NewButton("Добавить столбец", func() {
		if selectedTableIdx < 0 {
			showInfo("Выберите таблицу", mw.window)
			return
		}
		tables, _ := mw.utils.TableManager.GetTables()
		if selectedTableIdx < len(tables) {
			mw.showAddColumnForm(tables[selectedTableIdx], refreshTablesList)
		}
	})

	dropColumnBtn := widget.NewButton("Удалить столбец", func() {
		if selectedTableIdx < 0 {
			showInfo("Выберите таблицу", mw.window)
			return
		}
		tables, _ := mw.utils.TableManager.GetTables()
		if selectedTableIdx < len(tables) {
			mw.showDropColumnForm(tables[selectedTableIdx], refreshTablesList)
		}
	})

	viewStructureBtn := widget.NewButton("Просмотр структуры", func() {
		if selectedTableIdx < 0 {
			showInfo("Выберите таблицу", mw.window)
			return
		}
		tables, _ := mw.utils.TableManager.GetTables()
		if selectedTableIdx < len(tables) {
			mw.showTableStructure(tables[selectedTableIdx])
		}
	})

	refreshBtn := widget.NewButton("Обновить", refreshTablesList)

	buttons := container.NewHBox(createTableBtn, deleteTableBtn, addColumnBtn, dropColumnBtn, viewStructureBtn, refreshBtn)
	content := container.NewBorder(buttons, nil, nil, nil, tablesList)

	return container.NewTabItem("Управление таблицами", content)
}

func (mw *MainWindow) showCreateTableForm(onSuccess func()) {
	tableNameEntry := widget.NewEntry()
	tableNameEntry.SetPlaceHolder("Введите имя таблицы")

	columnsText := widget.NewMultiLineEntry()
	columnsText.SetPlaceHolder("Формат: имя_столбца тип_данных [NOT NULL] [DEFAULT значение]\nПример:\nid integer NOT NULL\nname varchar(255)\ncreated_at timestamp DEFAULT now()")
	columnsText.Wrapping = fyne.TextWrapOff

	// Create form with proper field sizing
	form := widget.NewForm()
	form.Append("Имя таблицы", formField(tableNameEntry))
	form.Append("Столбцы (по одному на строку)", formField(columnsText))

	form.OnSubmit = func() {
		// Simple implementation - in production, parse columnsText properly
		showInfo("Функция создания таблицы требует более сложной реализации парсинга", mw.window)
		onSuccess()
	}

	// Wrap form in container with proper sizing
	content := container.NewPadded(form)
	d := dialog.NewCustom("Создать таблицу", "Отмена", content, mw.window)

	// Set minimum dialog size to prevent overlapping
	minSize := content.MinSize()
	if minSize.Width < 520 {
		minSize.Width = 520
	}
	if minSize.Height < 400 {
		minSize.Height = 400
	}
	d.Resize(minSize)
	d.Show()
}

func (mw *MainWindow) showAddColumnForm(tableName string, onSuccess func()) {
	form := widget.NewForm()

	columnNameEntry := widget.NewEntry()
	columnNameEntry.SetPlaceHolder("Введите имя столбца")

	dataTypeEntry := widget.NewSelect([]string{"integer", "varchar(255)", "text", "numeric(15,2)", "date", "timestamp"}, nil)
	nullableCheck := widget.NewCheck("NULL разрешен", nil)
	nullableCheck.SetChecked(true)

	form.Append("Имя столбца", formField(columnNameEntry))
	form.Append("Тип данных", formField(dataTypeEntry))
	form.Append("", nullableCheck)

	form.OnSubmit = func() {
		if err := mw.utils.TableManager.AddColumn(tableName, columnNameEntry.Text, dataTypeEntry.Selected, nullableCheck.Checked); err != nil {
			showError(err, mw.window)
			return
		}
		showInfo("Столбец добавлен", mw.window)
		onSuccess()
	}

	content := container.NewPadded(form)
	d := dialog.NewCustom("Добавить столбец", "Отмена", content, mw.window)
	minSize := content.MinSize()
	if minSize.Width < 400 {
		minSize.Width = 400
	}
	if minSize.Height < 300 {
		minSize.Height = 300
	}
	d.Resize(minSize)
	d.Show()
}

func (mw *MainWindow) showDropColumnForm(tableName string, onSuccess func()) {
	columns, err := mw.utils.TableManager.GetTableColumns(tableName)
	if err != nil {
		showError(err, mw.window)
		return
	}

	columnNames := make([]string, len(columns))
	for i, col := range columns {
		columnNames[i] = col["name"]
	}

	form := widget.NewForm()
	columnSelect := widget.NewSelect(columnNames, nil)
	form.Append("Столбец для удаления", formField(columnSelect))

	form.OnSubmit = func() {
		if columnSelect.Selected == "" {
			showInfo("Выберите столбец", mw.window)
			return
		}
		dialog.ShowConfirm("Подтверждение", fmt.Sprintf("Удалить столбец %s?", columnSelect.Selected), func(ok bool) {
			if ok {
				if err := mw.utils.TableManager.DropColumn(tableName, columnSelect.Selected); err != nil {
					showError(err, mw.window)
				} else {
					showInfo("Столбец удален", mw.window)
					onSuccess()
				}
			}
		}, mw.window)
	}

	content := container.NewPadded(form)
	d := dialog.NewCustom("Удалить столбец", "Отмена", content, mw.window)
	minSize := content.MinSize()
	if minSize.Width < 400 {
		minSize.Width = 400
	}
	if minSize.Height < 200 {
		minSize.Height = 200
	}
	d.Resize(minSize)
	d.Show()
}

func (mw *MainWindow) showTableStructure(tableName string) {
	columns, err := mw.utils.TableManager.GetTableColumns(tableName)
	if err != nil {
		showError(err, mw.window)
		return
	}

	data := make([][]interface{}, len(columns))
	for i, col := range columns {
		data[i] = []interface{}{col["name"], col["type"], col["nullable"], col["default"]}
	}

	table := createTableWidget(data, []string{"Имя", "Тип", "NULL", "По умолчанию"})
	dialog.ShowCustom("Структура таблицы "+tableName, "Закрыть", table, mw.window)
}
