package ui

import (
	"bank_service/internal/models"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (mw *MainWindow) createBankBranchTab() *container.TabItem {
	var branches []models.BankBranch
	columns := []string{"ID", "Зона обслуживания", "Адрес", "Телефон", "Рабочие дни", "Часы работы", "Тип отделения"}

	table := widget.NewTable(
		func() (int, int) {
			return len(branches) + 1, len(columns)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id.Row == 0 {
				if id.Col < len(columns) {
					label.SetText(columns[id.Col])
				}
			} else {
				rowIdx := id.Row - 1
				if rowIdx < len(branches) {
					b := branches[rowIdx]
					var text string
					switch id.Col {
					case 0:
						text = fmt.Sprintf("%d", b.ID)
					case 1:
						text = b.ServiceZone
					case 2:
						text = b.FullAddress
					case 3:
						text = b.ContactPhone
					case 4:
						text = b.WorkingDays
					case 5:
						text = b.WorkingHours
					case 6:
						text = b.BranchType
					}
					label.SetText(text)
				}
			}
		},
	)

	// Set column widths
	table.SetColumnWidth(0, 60)  // ID
	table.SetColumnWidth(1, 180) // Зона обслуживания
	table.SetColumnWidth(2, 300) // Адрес
	table.SetColumnWidth(3, 160) // Телефон
	table.SetColumnWidth(4, 140) // Рабочие дни
	table.SetColumnWidth(5, 140) // Часы работы
	table.SetColumnWidth(6, 200) // Тип отделения

	refreshTable := func() {
		var err error
		branches, err = mw.repos.BankBranch.GetAll()
		if err != nil {
			showError(err, mw.window)
			return
		}
		table.Refresh()
	}

	var selectedRow int = -1
	table.OnSelected = func(id widget.TableCellID) {
		selectedRow = id.Row
	}

	createBtn := widget.NewButton("Создать", func() {
		mw.showBankBranchForm(nil, refreshTable)
	})

	updateBtn := widget.NewButton("Обновить", func() {
		if selectedRow <= 0 {
			showInfo("Выберите запись", mw.window)
			return
		}
		if selectedRow-1 < len(branches) {
			mw.showBankBranchForm(&branches[selectedRow-1], refreshTable)
		}
	})

	deleteBtn := widget.NewButton("Удалить", func() {
		if selectedRow <= 0 {
			showInfo("Выберите запись", mw.window)
			return
		}
		if selectedRow-1 < len(branches) {
			dialog.ShowConfirm("Подтверждение", "Удалить запись?", func(ok bool) {
				if ok {
					if err := mw.repos.BankBranch.Delete(branches[selectedRow-1].ID); err != nil {
						showError(err, mw.window)
					} else {
						refreshTable()
					}
				}
			}, mw.window)
		}
	})

	exportBtn := widget.NewButton("Экспорт в Excel", func() {
		filename := fmt.Sprintf("bankbranch_%s.xlsx", time.Now().Format("20060102_150405"))
		if err := mw.utils.ExportManager.ExportTableToExcel("bankbranch", filename); err != nil {
			showError(err, mw.window)
		} else {
			showInfo(fmt.Sprintf("Экспортировано в %s", filename), mw.window)
		}
	})

	refreshBtn := widget.NewButton("Обновить данные таблицы", refreshTable)

	buttons := container.NewHBox(createBtn, updateBtn, deleteBtn, exportBtn, refreshBtn)
	content := container.NewBorder(buttons, nil, nil, nil, table)

	return container.NewTabItem("Отделения", content)
}

func (mw *MainWindow) showBankBranchForm(branch *models.BankBranch, onSuccess func()) {
	form := widget.NewForm()

	serviceZoneEntry := widget.NewEntry()
	fullAddressEntry := widget.NewEntry()
	contactPhoneEntry := widget.NewEntry()
	workingDaysEntry := widget.NewEntry()
	workingHoursEntry := widget.NewEntry()
	branchTypeEntry := widget.NewEntry()

	if branch != nil {
		serviceZoneEntry.SetText(branch.ServiceZone)
		fullAddressEntry.SetText(branch.FullAddress)
		contactPhoneEntry.SetText(branch.ContactPhone)
		workingDaysEntry.SetText(branch.WorkingDays)
		workingHoursEntry.SetText(branch.WorkingHours)
		branchTypeEntry.SetText(branch.BranchType)
	}

	form.Append("Зона обслуживания", formField(serviceZoneEntry))
	form.Append("Адрес", formField(fullAddressEntry))
	form.Append("Телефон", formField(contactPhoneEntry))
	form.Append("Рабочие дни", formField(workingDaysEntry))
	form.Append("Часы работы", formField(workingHoursEntry))
	form.Append("Тип отделения", formField(branchTypeEntry))

	form.OnSubmit = func() {
		b := &models.BankBranch{
			ServiceZone:  serviceZoneEntry.Text,
			FullAddress:  fullAddressEntry.Text,
			ContactPhone: contactPhoneEntry.Text,
			WorkingDays:  workingDaysEntry.Text,
			WorkingHours: workingHoursEntry.Text,
			BranchType:   branchTypeEntry.Text,
		}

		if branch != nil {
			b.ID = branch.ID
			if err := mw.repos.BankBranch.Update(b); err != nil {
				showError(err, mw.window)
				return
			}
		} else {
			if err := mw.repos.BankBranch.Create(b); err != nil {
				showError(err, mw.window)
				return
			}
		}

		showInfo("Операция выполнена успешно", mw.window)
		onSuccess()
	}

	showFormDialog("Отделение банка", form, mw.window)
}
