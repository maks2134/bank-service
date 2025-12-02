package ui

import (
	"bank_service/internal/utils"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (mw *MainWindow) createQueriesTab() *container.TabItem {
	queryText := widget.NewMultiLineEntry()
	queryText.SetPlaceHolder("Введите SQL запрос здесь...")

	resultTable := widget.NewTable(
		func() (int, int) { return 0, 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, obj fyne.CanvasObject) {},
	)

	executeBtn := widget.NewButton("Выполнить запрос", func() {
		query := queryText.Text
		if query == "" {
			showInfo("Введите запрос", mw.window)
			return
		}

		data, columns, err := mw.utils.QueryManager.ExecuteQuery(query)
		if err != nil {
			showError(err, mw.window)
			return
		}

		resultTable = createTableWidget(data, columns)
		showInfo(fmt.Sprintf("Запрос выполнен. Найдено записей: %d", len(data)), mw.window)
	})

	saveQueryBtn := widget.NewButton("Сохранить запрос", func() {
		query := queryText.Text
		if query == "" {
			showInfo("Введите запрос для сохранения", mw.window)
			return
		}

		form := widget.NewForm()
		nameEntry := widget.NewEntry()
		descEntry := widget.NewMultiLineEntry()

		form.Append("Название запроса", nameEntry)
		form.Append("Описание", descEntry)

		form.OnSubmit = func() {
			if nameEntry.Text == "" {
				showInfo("Введите название запроса", mw.window)
				return
			}

			savedQuery := utils.SavedQuery{
				Name:        nameEntry.Text,
				Query:       query,
				Description: descEntry.Text,
			}

			if err := mw.utils.QueryManager.SaveQuery(savedQuery); err != nil {
				showError(err, mw.window)
				return
			}

			showInfo("Запрос сохранен", mw.window)
		}

		dialog.ShowCustom("Сохранить запрос", "Отмена", form, mw.window)
	})

	exportResultsBtn := widget.NewButton("Экспорт результатов в Excel", func() {
		query := queryText.Text
		if query == "" {
			showInfo("Выполните запрос сначала", mw.window)
			return
		}

		filename := fmt.Sprintf("query_results_%s.xlsx", time.Now().Format("20060102_150405"))
		if err := mw.utils.ExportManager.ExportQueryResultsToExcel(query, filename); err != nil {
			showError(err, mw.window)
		} else {
			showInfo(fmt.Sprintf("Результаты экспортированы в %s", filename), mw.window)
		}
	})

	// Saved queries list
	var savedQueriesList *widget.List
	var allQueries []utils.SavedQuery
	var refreshSavedQueries func()

	loadSavedQuery := func(query utils.SavedQuery) {
		queryText.SetText(query.Query)
	}

	refreshSavedQueries = func() {
		queries, err := mw.utils.QueryManager.GetSavedQueries()
		if err != nil {
			showError(err, mw.window)
			return
		}

		// Add predefined queries
		predefined := mw.utils.QueryManager.GetPredefinedQueries()
		allQueries = append(predefined, queries...)

		savedQueriesList = widget.NewList(
			func() int { return len(allQueries) },
			func() fyne.CanvasObject {
				return container.NewHBox(
					widget.NewLabel(""),
					widget.NewButton("Загрузить", nil),
					widget.NewButton("Удалить", nil),
				)
			},
			func(id widget.ListItemID, obj fyne.CanvasObject) {
				if id >= len(allQueries) {
					return
				}
				query := allQueries[id]

				// Update existing widgets in the container
				if hbox, ok := obj.(interface{ Objects() []fyne.CanvasObject }); ok {
					objs := hbox.Objects()
					if len(objs) >= 3 {
						// Update label
						if lbl, ok := objs[0].(*widget.Label); ok {
							lbl.SetText(fmt.Sprintf("%s - %s", query.Name, query.Description))
						}
						// Update load button
						if btn, ok := objs[1].(*widget.Button); ok {
							btn.OnTapped = func() {
								loadSavedQuery(query)
							}
						}
						// Update delete button
						if btn, ok := objs[2].(*widget.Button); ok {
							if id < len(predefined) {
								btn.Hide()
							} else {
								btn.Show()
								btn.OnTapped = func() {
									dialog.ShowConfirm("Подтверждение", fmt.Sprintf("Удалить запрос '%s'?", query.Name), func(ok bool) {
										if ok {
											if err := mw.utils.QueryManager.DeleteQuery(query.Name); err != nil {
												showError(err, mw.window)
											} else {
												refreshSavedQueries()
											}
										}
									}, mw.window)
								}
							}
						}
					}
				}
			},
		)
	}

	refreshSavedQueries()

	refreshSavedBtn := widget.NewButton("Обновить список", refreshSavedQueries)
	leftPanel := container.NewBorder(
		widget.NewLabel("Сохраненные запросы"),
		refreshSavedBtn,
		nil,
		nil,
		savedQueriesList,
	)

	resultContainer := container.NewBorder(nil, resultTable, nil, nil, queryText)
	rightPanel := container.NewBorder(
		container.NewHBox(executeBtn, saveQueryBtn, exportResultsBtn),
		nil,
		nil,
		nil,
		resultContainer,
	)

	split := container.NewHSplit(leftPanel, rightPanel)
	split.SetOffset(0.3)

	return container.NewTabItem("Запросы", split)
}
