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

func (mw *MainWindow) createBankAccountTab() *container.TabItem {
	var accounts []models.BankAccount
	columns := []string{"ID", "Тип счета", "Номер счета", "Баланс", "Валюта", "Дата открытия", "Статус", "ID клиента"}

	table := widget.NewTable(
		func() (int, int) {
			return len(accounts) + 1, len(columns)
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
				if rowIdx < len(accounts) {
					acc := accounts[rowIdx]
					var text string
					switch id.Col {
					case 0:
						text = strconv.Itoa(acc.ID)
					case 1:
						text = acc.AccountType
					case 2:
						text = acc.AccountNumber
					case 3:
						text = fmt.Sprintf("%.2f", acc.Balance)
					case 4:
						text = acc.Currency
					case 5:
						text = acc.OpenDate
					case 6:
						text = acc.AccountStatus
					case 7:
						if acc.ClientID != nil {
							text = strconv.Itoa(*acc.ClientID)
						}
					}
					label.SetText(text)
				}
			}
		},
	)

	// Set column widths for bank accounts table
	table.SetColumnWidth(0, 60)  // ID
	table.SetColumnWidth(1, 140) // Тип счета
	table.SetColumnWidth(2, 200) // Номер счета
	table.SetColumnWidth(3, 120) // Баланс
	table.SetColumnWidth(4, 100) // Валюта
	table.SetColumnWidth(5, 180) // Дата открытия
	table.SetColumnWidth(6, 120) // Статус
	table.SetColumnWidth(7, 100) // ID клиента

	refreshTable := func() {
		var err error
		accounts, err = mw.repos.BankAccount.GetAll()
		if err != nil {
			showError(err, mw.window)
			return
		}
		table.Refresh()
	}

	createBtn := widget.NewButton("Создать", func() {
		mw.showBankAccountForm(nil, refreshTable)
	})

	var selectedRow int = -1
	table.OnSelected = func(id widget.TableCellID) {
		selectedRow = id.Row
	}

	updateBtn := widget.NewButton("Обновить", func() {
		if selectedRow <= 0 {
			showInfo("Выберите запись для обновления", mw.window)
			return
		}
		if selectedRow-1 < len(accounts) {
			mw.showBankAccountForm(&accounts[selectedRow-1], refreshTable)
		}
	})

	deleteBtn := widget.NewButton("Удалить", func() {
		if selectedRow <= 0 {
			showInfo("Выберите запись для удаления", mw.window)
			return
		}
		if selectedRow-1 < len(accounts) {
			dialog.ShowConfirm("Подтверждение", "Удалить запись?", func(ok bool) {
				if ok {
					if err := mw.repos.BankAccount.Delete(accounts[selectedRow-1].ID); err != nil {
						showError(err, mw.window)
					} else {
						refreshTable()
					}
				}
			}, mw.window)
		}
	})

	exportBtn := widget.NewButton("Экспорт в Excel", func() {
		filename := fmt.Sprintf("bankaccount_%s.xlsx", time.Now().Format("20060102_150405"))
		if err := mw.utils.ExportManager.ExportTableToExcel("bankaccount", filename); err != nil {
			showError(err, mw.window)
		} else {
			showInfo(fmt.Sprintf("Экспортировано в %s", filename), mw.window)
		}
	})

	refreshBtn := widget.NewButton("Обновить данные таблицы", refreshTable)

	buttons := container.NewHBox(createBtn, updateBtn, deleteBtn, exportBtn, refreshBtn)
	content := container.NewBorder(buttons, nil, nil, nil, table)
	refreshTable()

	return container.NewTabItem("Банковские счета", content)
}

func (mw *MainWindow) showBankAccountForm(acc *models.BankAccount, onSuccess func()) {
	form := widget.NewForm()

	idEntry := widget.NewEntry()
	accountTypeEntry := widget.NewEntry()
	accountNumberEntry := widget.NewEntry()
	balanceEntry := widget.NewEntry()
	currencyEntry := widget.NewEntry()
	openDateEntry := widget.NewEntry()
	accountStatusEntry := widget.NewSelect([]string{"ACTIVE", "BLOCKED", "CLOSED", "SUSPENDED"}, nil)
	clientIDEntry := widget.NewEntry()

	if acc != nil {
		idEntry.SetText(strconv.Itoa(acc.ID))
		idEntry.Disable()
		accountTypeEntry.SetText(acc.AccountType)
		accountNumberEntry.SetText(acc.AccountNumber)
		balanceEntry.SetText(fmt.Sprintf("%.2f", acc.Balance))
		currencyEntry.SetText(acc.Currency)
		openDateEntry.SetText(acc.OpenDate)
		accountStatusEntry.SetSelected(acc.AccountStatus)
		if acc.ClientID != nil {
			clientIDEntry.SetText(strconv.Itoa(*acc.ClientID))
		}
	} else {
		openDateEntry.SetText(time.Now().Format("2006-01-02"))
		accountStatusEntry.SetSelected("ACTIVE")
	}

	form.Append("Тип счета", formField(accountTypeEntry))
	form.Append("Номер счета", formField(accountNumberEntry))
	form.Append("Баланс", formField(balanceEntry))
	form.Append("Валюта", formField(currencyEntry))
	form.Append("Дата открытия", formField(openDateEntry))
	form.Append("Статус", formField(accountStatusEntry))
	form.Append("ID клиента", formField(clientIDEntry))

	form.OnSubmit = func() {
		balance, err := strconv.ParseFloat(balanceEntry.Text, 64)
		if err != nil {
			showError(fmt.Errorf("неверный баланс: %w", err), mw.window)
			return
		}

		var clientID *int
		if clientIDEntry.Text != "" {
			id, err := strconv.Atoi(clientIDEntry.Text)
			if err == nil {
				clientID = &id
			}
		}

		account := &models.BankAccount{
			AccountType:   accountTypeEntry.Text,
			AccountNumber: accountNumberEntry.Text,
			Balance:       balance,
			Currency:      currencyEntry.Text,
			OpenDate:      openDateEntry.Text,
			AccountStatus: accountStatusEntry.Selected,
			ClientID:      clientID,
		}

		if acc != nil {
			account.ID = acc.ID
			if err := mw.repos.BankAccount.Update(account); err != nil {
				showError(err, mw.window)
				return
			}
		} else {
			if err := mw.repos.BankAccount.Create(account); err != nil {
				showError(err, mw.window)
				return
			}
		}

		showInfo("Операция выполнена успешно", mw.window)
		onSuccess()
	}

	showFormDialog("Банковский счет", form, mw.window)
}
