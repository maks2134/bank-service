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

func (mw *MainWindow) createClientTab() *container.TabItem {
	table := widget.NewTable(
		func() (int, int) { return 0, 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, obj fyne.CanvasObject) {},
	)

	// Set column widths
	table.SetColumnWidth(0, 50)  // ID
	table.SetColumnWidth(1, 150) // Телефон
	table.SetColumnWidth(2, 150) // Логин
	table.SetColumnWidth(3, 200) // ФИО
	table.SetColumnWidth(4, 120) // Тип клиента
	table.SetColumnWidth(5, 200) // Email
	table.SetColumnWidth(6, 120) // Паспорт
	table.SetColumnWidth(7, 80)  // ID кредита

	refreshTable := func() {
		clients, err := mw.repos.Client.GetAll()
		if err != nil {
			showError(err, mw.window)
			return
		}

		data := make([][]interface{}, len(clients))
		for i, c := range clients {
			creditID := ""
			if c.CreditID != nil {
				creditID = strconv.Itoa(*c.CreditID)
			}
			data[i] = []interface{}{
				c.ID, c.Phone, c.Login, c.FullName,
				c.ClientType, c.Email, c.Passport, creditID,
			}
		}

		columns := []string{"ID", "Телефон", "Логин", "ФИО", "Тип клиента", "Email", "Паспорт", "ID кредита"}
		table = createTableWidget(data, columns)
		// Set column widths after creating table
		table.SetColumnWidth(0, 60)  // ID
		table.SetColumnWidth(1, 160) // Телефон
		table.SetColumnWidth(2, 160) // Логин
		table.SetColumnWidth(3, 250) // ФИО
		table.SetColumnWidth(4, 150) // Тип клиента
		table.SetColumnWidth(5, 220) // Email
		table.SetColumnWidth(6, 140) // Паспорт
		table.SetColumnWidth(7, 100) // ID кредита
	}

	refreshTable()

	var selectedRow int = -1
	table.OnSelected = func(id widget.TableCellID) {
		selectedRow = id.Row
	}

	createBtn := widget.NewButton("Создать", func() {
		mw.showClientForm(nil, refreshTable)
	})

	updateBtn := widget.NewButton("Обновить", func() {
		if selectedRow <= 0 {
			showInfo("Выберите запись", mw.window)
			return
		}
		clients, _ := mw.repos.Client.GetAll()
		if selectedRow-1 < len(clients) {
			mw.showClientForm(&clients[selectedRow-1], refreshTable)
		}
	})

	deleteBtn := widget.NewButton("Удалить", func() {
		if selectedRow <= 0 {
			showInfo("Выберите запись", mw.window)
			return
		}
		clients, _ := mw.repos.Client.GetAll()
		if selectedRow-1 < len(clients) {
			dialog.ShowConfirm("Подтверждение", "Удалить запись?", func(ok bool) {
				if ok {
					if err := mw.repos.Client.Delete(clients[selectedRow-1].ID); err != nil {
						showError(err, mw.window)
					} else {
						refreshTable()
					}
				}
			}, mw.window)
		}
	})

	exportBtn := widget.NewButton("Экспорт в Excel", func() {
		filename := fmt.Sprintf("client_%s.xlsx", time.Now().Format("20060102_150405"))
		if err := mw.utils.ExportManager.ExportTableToExcel("client", filename); err != nil {
			showError(err, mw.window)
		} else {
			showInfo(fmt.Sprintf("Экспортировано в %s", filename), mw.window)
		}
	})

	refreshBtn := widget.NewButton("Обновить данные таблицы", refreshTable)

	buttons := container.NewHBox(createBtn, updateBtn, deleteBtn, exportBtn, refreshBtn)
	content := container.NewBorder(buttons, nil, nil, nil, table)

	return container.NewTabItem("Клиенты", content)
}

func (mw *MainWindow) showClientForm(c *models.Client, onSuccess func()) {
	form := widget.NewForm()

	phoneEntry := widget.NewEntry()
	loginEntry := widget.NewEntry()
	fullNameEntry := widget.NewEntry()
	clientTypeEntry := widget.NewSelect([]string{"Физическое лицо", "Юридическое лицо"}, nil)
	passwordEntry := widget.NewPasswordEntry()
	emailEntry := widget.NewEntry()
	passportEntry := widget.NewEntry()
	creditIDEntry := widget.NewEntry()

	if c != nil {
		phoneEntry.SetText(c.Phone)
		loginEntry.SetText(c.Login)
		fullNameEntry.SetText(c.FullName)
		clientTypeEntry.SetSelected(c.ClientType)
		passwordEntry.SetText(c.Password)
		emailEntry.SetText(c.Email)
		passportEntry.SetText(c.Passport)
		if c.CreditID != nil {
			creditIDEntry.SetText(strconv.Itoa(*c.CreditID))
		}
	} else {
		clientTypeEntry.SetSelected("Физическое лицо")
	}

	form.Append("Телефон", formField(phoneEntry))
	form.Append("Логин", formField(loginEntry))
	form.Append("ФИО", formField(fullNameEntry))
	form.Append("Тип клиента", formField(clientTypeEntry))
	form.Append("Пароль", formField(passwordEntry))
	form.Append("Email", formField(emailEntry))
	form.Append("Паспорт", formField(passportEntry))
	form.Append("ID кредита", formField(creditIDEntry))

	form.OnSubmit = func() {
		var creditID *int
		if creditIDEntry.Text != "" {
			id, err := strconv.Atoi(creditIDEntry.Text)
			if err == nil {
				creditID = &id
			}
		}

		client := &models.Client{
			Phone:      phoneEntry.Text,
			Login:      loginEntry.Text,
			FullName:   fullNameEntry.Text,
			ClientType: clientTypeEntry.Selected,
			Password:   passwordEntry.Text,
			Email:      emailEntry.Text,
			Passport:   passportEntry.Text,
			CreditID:   creditID,
		}

		if c != nil {
			client.ID = c.ID
			if err := mw.repos.Client.Update(client); err != nil {
				showError(err, mw.window)
				return
			}
		} else {
			if err := mw.repos.Client.Create(client); err != nil {
				showError(err, mw.window)
				return
			}
		}

		showInfo("Операция выполнена успешно", mw.window)
		onSuccess()
	}

	showFormDialog("Клиент", form, mw.window)
}
