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

func (mw *MainWindow) createBankStaffTab() *container.TabItem {
	table := widget.NewTable(
		func() (int, int) { return 0, 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, obj fyne.CanvasObject) {},
	)

	// Set column widths
	table.SetColumnWidth(0, 50)  // ID
	table.SetColumnWidth(1, 200) // ФИО
	table.SetColumnWidth(2, 120) // Паспорт
	table.SetColumnWidth(3, 180) // Должность
	table.SetColumnWidth(4, 120) // Дата найма
	table.SetColumnWidth(5, 120) // Уровень доступа
	table.SetColumnWidth(6, 200) // Квалификация
	table.SetColumnWidth(7, 80)  // ID отделения

	refreshTable := func() {
		staff, err := mw.repos.BankStaff.GetAll()
		if err != nil {
			showError(err, mw.window)
			return
		}

		data := make([][]interface{}, len(staff))
		for i, s := range staff {
			qual := ""
			if s.Qualification != nil {
				qual = *s.Qualification
			}
			data[i] = []interface{}{
				s.ID, s.FullName, s.Passport, s.Position,
				s.HireDate, s.AccessLevel, qual, s.BranchID,
			}
		}

		columns := []string{"ID", "ФИО", "Паспорт", "Должность", "Дата найма", "Уровень доступа", "Квалификация", "ID отделения"}
		table = createTableWidget(data, columns)
		// Set column widths after creating table
		table.SetColumnWidth(0, 60)  // ID
		table.SetColumnWidth(1, 250) // ФИО
		table.SetColumnWidth(2, 140) // Паспорт
		table.SetColumnWidth(3, 220) // Должность
		table.SetColumnWidth(4, 150) // Дата найма
		table.SetColumnWidth(5, 150) // Уровень доступа
		table.SetColumnWidth(6, 250) // Квалификация
		table.SetColumnWidth(7, 120) // ID отделения
	}

	refreshTable()

	var selectedRow int = -1
	table.OnSelected = func(id widget.TableCellID) {
		selectedRow = id.Row
	}

	createBtn := widget.NewButton("Создать", func() {
		mw.showBankStaffForm(nil, refreshTable)
	})

	updateBtn := widget.NewButton("Обновить", func() {
		if selectedRow <= 0 {
			showInfo("Выберите запись", mw.window)
			return
		}
		staff, _ := mw.repos.BankStaff.GetAll()
		if selectedRow-1 < len(staff) {
			mw.showBankStaffForm(&staff[selectedRow-1], refreshTable)
		}
	})

	deleteBtn := widget.NewButton("Удалить", func() {
		if selectedRow <= 0 {
			showInfo("Выберите запись", mw.window)
			return
		}
		staff, _ := mw.repos.BankStaff.GetAll()
		if selectedRow-1 < len(staff) {
			dialog.ShowConfirm("Подтверждение", "Удалить запись?", func(ok bool) {
				if ok {
					if err := mw.repos.BankStaff.Delete(staff[selectedRow-1].ID); err != nil {
						showError(err, mw.window)
					} else {
						refreshTable()
					}
				}
			}, mw.window)
		}
	})

	exportBtn := widget.NewButton("Экспорт в Excel", func() {
		filename := fmt.Sprintf("bankstaff_%s.xlsx", time.Now().Format("20060102_150405"))
		if err := mw.utils.ExportManager.ExportTableToExcel("bankstaff", filename); err != nil {
			showError(err, mw.window)
		} else {
			showInfo(fmt.Sprintf("Экспортировано в %s", filename), mw.window)
		}
	})

	refreshBtn := widget.NewButton("Обновить данные таблицы", refreshTable)

	buttons := container.NewHBox(createBtn, updateBtn, deleteBtn, exportBtn, refreshBtn)
	content := container.NewBorder(buttons, nil, nil, nil, table)

	return container.NewTabItem("Сотрудники", content)
}

func (mw *MainWindow) showBankStaffForm(s *models.BankStaff, onSuccess func()) {
	form := widget.NewForm()

	fullNameEntry := widget.NewEntry()
	passportEntry := widget.NewEntry()
	positionEntry := widget.NewEntry()
	hireDateEntry := widget.NewEntry()
	accessLevelEntry := widget.NewEntry()
	qualificationEntry := widget.NewEntry()
	branchIDEntry := widget.NewEntry()

	if s != nil {
		fullNameEntry.SetText(s.FullName)
		passportEntry.SetText(s.Passport)
		positionEntry.SetText(s.Position)
		hireDateEntry.SetText(s.HireDate)
		accessLevelEntry.SetText(s.AccessLevel)
		if s.Qualification != nil {
			qualificationEntry.SetText(*s.Qualification)
		}
		branchIDEntry.SetText(strconv.Itoa(s.BranchID))
	} else {
		hireDateEntry.SetText(time.Now().Format("2006-01-02"))
	}

	form.Append("ФИО", formField(fullNameEntry))
	form.Append("Паспорт", formField(passportEntry))
	form.Append("Должность", formField(positionEntry))
	form.Append("Дата найма", formField(hireDateEntry))
	form.Append("Уровень доступа", formField(accessLevelEntry))
	form.Append("Квалификация", formField(qualificationEntry))
	form.Append("ID отделения", formField(branchIDEntry))

	form.OnSubmit = func() {
		branchID, err := strconv.Atoi(branchIDEntry.Text)
		if err != nil {
			showError(fmt.Errorf("неверный ID отделения: %w", err), mw.window)
			return
		}

		var qual *string
		if qualificationEntry.Text != "" {
			qual = &qualificationEntry.Text
		}

		staff := &models.BankStaff{
			FullName:      fullNameEntry.Text,
			Passport:      passportEntry.Text,
			Position:      positionEntry.Text,
			HireDate:      hireDateEntry.Text,
			AccessLevel:   accessLevelEntry.Text,
			Qualification: qual,
			BranchID:      branchID,
		}

		if s != nil {
			staff.ID = s.ID
			if err := mw.repos.BankStaff.Update(staff); err != nil {
				showError(err, mw.window)
				return
			}
		} else {
			if err := mw.repos.BankStaff.Create(staff); err != nil {
				showError(err, mw.window)
				return
			}
		}

		showInfo("Операция выполнена успешно", mw.window)
		onSuccess()
	}

	showFormDialog("Сотрудник банка", form, mw.window)
}
