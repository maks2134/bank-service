package ui

import (
	"bank_service/internal/models"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (mw *MainWindow) createCreditBankStaffTab() *container.TabItem {
	table := widget.NewTable(
		func() (int, int) { return 0, 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, obj fyne.CanvasObject) {},
	)

	// Set column widths
	table.SetColumnWidth(0, 100) // ID кредита
	table.SetColumnWidth(1, 100) // ID сотрудника

	refreshTable := func() {
		relations, err := mw.repos.CreditBankStaff.GetAll()
		if err != nil {
			showError(err, mw.window)
			return
		}

		data := make([][]interface{}, len(relations))
		for i, r := range relations {
			data[i] = []interface{}{r.CreditID, r.StaffID}
		}

		columns := []string{"ID кредита", "ID сотрудника"}
		table = createTableWidget(data, columns)
		// Set column widths after creating table
		table.SetColumnWidth(0, 100) // ID кредита
		table.SetColumnWidth(1, 100) // ID сотрудника
	}

	refreshTable()

	var selectedRow int = -1
	table.OnSelected = func(id widget.TableCellID) {
		selectedRow = id.Row
	}

	createBtn := widget.NewButton("Создать", func() {
		mw.showCreditBankStaffForm(nil, refreshTable)
	})

	deleteBtn := widget.NewButton("Удалить", func() {
		if selectedRow <= 0 {
			showInfo("Выберите запись", mw.window)
			return
		}
		relations, _ := mw.repos.CreditBankStaff.GetAll()
		if selectedRow-1 < len(relations) {
			dialog.ShowConfirm("Подтверждение", "Удалить запись?", func(ok bool) {
				if ok {
					if err := mw.repos.CreditBankStaff.Delete(relations[selectedRow-1].CreditID, relations[selectedRow-1].StaffID); err != nil {
						showError(err, mw.window)
					} else {
						refreshTable()
					}
				}
			}, mw.window)
		}
	})

	exportBtn := widget.NewButton("Экспорт в Excel", func() {
		filename := fmt.Sprintf("creditbankstaff_%s.xlsx", time.Now().Format("20060102_150405"))
		if err := mw.utils.ExportManager.ExportTableToExcel("creditbankstaff", filename); err != nil {
			showError(err, mw.window)
		} else {
			showInfo(fmt.Sprintf("Экспортировано в %s", filename), mw.window)
		}
	})

	refreshBtn := widget.NewButton("Обновить", refreshTable)

	buttons := container.NewHBox(createBtn, deleteBtn, exportBtn, refreshBtn)
	content := container.NewBorder(buttons, nil, nil, nil, table)

	return container.NewTabItem("Связь кредит-сотрудник", content)
}

func (mw *MainWindow) showCreditBankStaffForm(rel *models.CreditBankStaff, onSuccess func()) {
	form := widget.NewForm()

	creditIDEntry := widget.NewEntry()
	staffIDEntry := widget.NewEntry()

	if rel != nil {
		creditIDEntry.SetText(strconv.Itoa(rel.CreditID))
		staffIDEntry.SetText(strconv.Itoa(rel.StaffID))
	}

	form.Append("ID кредита", creditIDEntry)
	form.Append("ID сотрудника", staffIDEntry)

	form.OnSubmit = func() {
		creditID, err := strconv.Atoi(creditIDEntry.Text)
		if err != nil {
			showError(fmt.Errorf("неверный ID кредита: %w", err), mw.window)
			return
		}

		staffID, err := strconv.Atoi(staffIDEntry.Text)
		if err != nil {
			showError(fmt.Errorf("неверный ID сотрудника: %w", err), mw.window)
			return
		}

		relation := &models.CreditBankStaff{
			CreditID: creditID,
			StaffID:  staffID,
		}

		if err := mw.repos.CreditBankStaff.Create(relation); err != nil {
			showError(err, mw.window)
			return
		}

		showInfo("Операция выполнена успешно", mw.window)
		onSuccess()
	}

	dialog.ShowCustom("Связь кредит-сотрудник", "Отмена", form, mw.window)
}
