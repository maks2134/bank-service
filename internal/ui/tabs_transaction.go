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

func (mw *MainWindow) createTransactionTab() *container.TabItem {
	table := widget.NewTable(
		func() (int, int) { return 0, 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, obj fyne.CanvasObject) {},
	)

	refreshTable := func() {
		transactions, err := mw.repos.Transaction.GetAll()
		if err != nil {
			showError(err, mw.window)
			return
		}

		data := make([][]interface{}, len(transactions))
		for i, t := range transactions {
			purpose := ""
			if t.Purpose != nil {
				purpose = *t.Purpose
			}
			data[i] = []interface{}{
				t.ID, t.Amount, t.OperationDate.Format("2006-01-02 15:04:05"),
				t.OperationType, purpose, t.OperationStatus, t.Currency,
				t.BranchID, t.AccountID,
			}
		}

		columns := []string{"ID", "Сумма", "Дата операции", "Тип операции", "Назначение", "Статус", "Валюта", "ID отделения", "ID счета"}
		table = createTableWidget(data, columns)
	}

	refreshTable()

	// Set column widths
	table.SetColumnWidth(0, 60)  // ID
	table.SetColumnWidth(1, 120) // Сумма
	table.SetColumnWidth(2, 200) // Дата операции
	table.SetColumnWidth(3, 180) // Тип операции
	table.SetColumnWidth(4, 250) // Назначение
	table.SetColumnWidth(5, 140) // Статус
	table.SetColumnWidth(6, 100) // Валюта
	table.SetColumnWidth(7, 120) // ID отделения
	table.SetColumnWidth(8, 100) // ID счета

	var selectedRow int = -1
	table.OnSelected = func(id widget.TableCellID) {
		selectedRow = id.Row
	}

	createBtn := widget.NewButton("Создать", func() {
		mw.showTransactionForm(nil, refreshTable)
	})

	updateBtn := widget.NewButton("Обновить", func() {
		if selectedRow <= 0 {
			showInfo("Выберите запись", mw.window)
			return
		}
		transactions, _ := mw.repos.Transaction.GetAll()
		if selectedRow-1 < len(transactions) {
			mw.showTransactionForm(&transactions[selectedRow-1], refreshTable)
		}
	})

	deleteBtn := widget.NewButton("Удалить", func() {
		if selectedRow <= 0 {
			showInfo("Выберите запись", mw.window)
			return
		}
		transactions, _ := mw.repos.Transaction.GetAll()
		if selectedRow-1 < len(transactions) {
			dialog.ShowConfirm("Подтверждение", "Удалить запись?", func(ok bool) {
				if ok {
					if err := mw.repos.Transaction.Delete(transactions[selectedRow-1].ID); err != nil {
						showError(err, mw.window)
					} else {
						refreshTable()
					}
				}
			}, mw.window)
		}
	})

	exportBtn := widget.NewButton("Экспорт в Excel", func() {
		filename := fmt.Sprintf("transaction_%s.xlsx", time.Now().Format("20060102_150405"))
		if err := mw.utils.ExportManager.ExportTableToExcel("transaction", filename); err != nil {
			showError(err, mw.window)
		} else {
			showInfo(fmt.Sprintf("Экспортировано в %s", filename), mw.window)
		}
	})

	refreshBtn := widget.NewButton("Обновить данные таблицы", refreshTable)

	buttons := container.NewHBox(createBtn, updateBtn, deleteBtn, exportBtn, refreshBtn)
	content := container.NewBorder(buttons, nil, nil, nil, table)

	return container.NewTabItem("Транзакции", content)
}

func (mw *MainWindow) showTransactionForm(t *models.Transaction, onSuccess func()) {
	form := widget.NewForm()

	amountEntry := widget.NewEntry()
	operationDateEntry := widget.NewEntry()
	operationTypeEntry := widget.NewEntry()
	purposeEntry := widget.NewEntry()
	operationStatusEntry := widget.NewEntry()
	currencyEntry := widget.NewEntry()
	branchIDEntry := widget.NewEntry()
	accountIDEntry := widget.NewEntry()

	if t != nil {
		amountEntry.SetText(fmt.Sprintf("%.2f", t.Amount))
		operationDateEntry.SetText(t.OperationDate.Format("2006-01-02 15:04:05"))
		operationTypeEntry.SetText(t.OperationType)
		if t.Purpose != nil {
			purposeEntry.SetText(*t.Purpose)
		}
		operationStatusEntry.SetText(t.OperationStatus)
		currencyEntry.SetText(t.Currency)
		branchIDEntry.SetText(strconv.Itoa(t.BranchID))
		accountIDEntry.SetText(strconv.Itoa(t.AccountID))
	} else {
		operationDateEntry.SetText(time.Now().Format("2006-01-02 15:04:05"))
	}

	form.Append("Сумма", amountEntry)
	form.Append("Дата операции", operationDateEntry)
	form.Append("Тип операции", operationTypeEntry)
	form.Append("Назначение", purposeEntry)
	form.Append("Статус", operationStatusEntry)
	form.Append("Валюта", currencyEntry)
	form.Append("ID отделения", branchIDEntry)
	form.Append("ID счета", accountIDEntry)

	form.OnSubmit = func() {
		amount, err := strconv.ParseFloat(amountEntry.Text, 64)
		if err != nil {
			showError(fmt.Errorf("неверная сумма: %w", err), mw.window)
			return
		}

		operationDate, err := time.Parse("2006-01-02 15:04:05", operationDateEntry.Text)
		if err != nil {
			showError(fmt.Errorf("неверная дата: %w", err), mw.window)
			return
		}

		branchID, err := strconv.Atoi(branchIDEntry.Text)
		if err != nil {
			showError(fmt.Errorf("неверный ID отделения: %w", err), mw.window)
			return
		}

		accountID, err := strconv.Atoi(accountIDEntry.Text)
		if err != nil {
			showError(fmt.Errorf("неверный ID счета: %w", err), mw.window)
			return
		}

		var purpose *string
		if purposeEntry.Text != "" {
			purpose = &purposeEntry.Text
		}

		transaction := &models.Transaction{
			Amount:          amount,
			OperationDate:   operationDate,
			OperationType:   operationTypeEntry.Text,
			Purpose:         purpose,
			OperationStatus: operationStatusEntry.Text,
			Currency:        currencyEntry.Text,
			BranchID:        branchID,
			AccountID:       accountID,
		}

		if t != nil {
			transaction.ID = t.ID
			if err := mw.repos.Transaction.Update(transaction); err != nil {
				showError(err, mw.window)
				return
			}
		} else {
			if err := mw.repos.Transaction.Create(transaction); err != nil {
				showError(err, mw.window)
				return
			}
		}

		showInfo("Операция выполнена успешно", mw.window)
		onSuccess()
	}

	dialog.ShowCustom("Транзакция", "Отмена", form, mw.window)
}
