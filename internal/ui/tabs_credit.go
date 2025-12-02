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

func (mw *MainWindow) createCreditTab() *container.TabItem {
	table := widget.NewTable(
		func() (int, int) { return 0, 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, obj fyne.CanvasObject) {},
	)

	// Set column widths
	table.SetColumnWidth(0, 50)  // ID
	table.SetColumnWidth(1, 100) // Сумма
	table.SetColumnWidth(2, 120) // Процентная ставка
	table.SetColumnWidth(3, 200) // Цель
	table.SetColumnWidth(4, 120) // Дата выдачи
	table.SetColumnWidth(5, 100) // Статус
	table.SetColumnWidth(6, 120) // Дата погашения
	table.SetColumnWidth(7, 80)  // Валюта
	table.SetColumnWidth(8, 100) // Уровень риска

	refreshTable := func() {
		credits, err := mw.repos.Credit.GetAll()
		if err != nil {
			showError(err, mw.window)
			return
		}

		data := make([][]interface{}, len(credits))
		for i, c := range credits {
			repaymentDate := ""
			if c.RepaymentDate != nil {
				repaymentDate = *c.RepaymentDate
			}
			riskLevel := ""
			if c.RiskLevel != nil {
				riskLevel = *c.RiskLevel
			}
			data[i] = []interface{}{
				c.ID, c.Amount, c.InterestRate, c.Purpose,
				c.IssueDate, c.Status, repaymentDate, c.Currency, riskLevel,
			}
		}

		columns := []string{"ID", "Сумма", "Процентная ставка", "Цель", "Дата выдачи", "Статус", "Дата погашения", "Валюта", "Уровень риска"}
		table = createTableWidget(data, columns)
		// Set column widths after creating table
		table.SetColumnWidth(0, 60)  // ID
		table.SetColumnWidth(1, 120) // Сумма
		table.SetColumnWidth(2, 150) // Процентная ставка
		table.SetColumnWidth(3, 250) // Цель
		table.SetColumnWidth(4, 180) // Дата выдачи
		table.SetColumnWidth(5, 120) // Статус
		table.SetColumnWidth(6, 180) // Дата погашения
		table.SetColumnWidth(7, 100) // Валюта
		table.SetColumnWidth(8, 130) // Уровень риска
	}

	refreshTable()

	var selectedRow int = -1
	table.OnSelected = func(id widget.TableCellID) {
		selectedRow = id.Row
	}

	createBtn := widget.NewButton("Создать", func() {
		mw.showCreditForm(nil, refreshTable)
	})

	updateBtn := widget.NewButton("Обновить", func() {
		if selectedRow <= 0 {
			showInfo("Выберите запись", mw.window)
			return
		}
		credits, _ := mw.repos.Credit.GetAll()
		if selectedRow-1 < len(credits) {
			mw.showCreditForm(&credits[selectedRow-1], refreshTable)
		}
	})

	deleteBtn := widget.NewButton("Удалить", func() {
		if selectedRow <= 0 {
			showInfo("Выберите запись", mw.window)
			return
		}
		credits, _ := mw.repos.Credit.GetAll()
		if selectedRow-1 < len(credits) {
			dialog.ShowConfirm("Подтверждение", "Удалить запись?", func(ok bool) {
				if ok {
					if err := mw.repos.Credit.Delete(credits[selectedRow-1].ID); err != nil {
						showError(err, mw.window)
					} else {
						refreshTable()
					}
				}
			}, mw.window)
		}
	})

	exportBtn := widget.NewButton("Экспорт в Excel", func() {
		filename := fmt.Sprintf("credit_%s.xlsx", time.Now().Format("20060102_150405"))
		if err := mw.utils.ExportManager.ExportTableToExcel("credit", filename); err != nil {
			showError(err, mw.window)
		} else {
			showInfo(fmt.Sprintf("Экспортировано в %s", filename), mw.window)
		}
	})

	refreshBtn := widget.NewButton("Обновить данные таблицы", refreshTable)

	buttons := container.NewHBox(createBtn, updateBtn, deleteBtn, exportBtn, refreshBtn)
	content := container.NewBorder(buttons, nil, nil, nil, table)

	return container.NewTabItem("Кредиты", content)
}

func (mw *MainWindow) showCreditForm(c *models.Credit, onSuccess func()) {
	form := widget.NewForm()

	amountEntry := widget.NewEntry()
	interestRateEntry := widget.NewEntry()
	purposeEntry := widget.NewEntry()
	issueDateEntry := widget.NewEntry()
	statusEntry := widget.NewEntry()
	repaymentDateEntry := widget.NewEntry()
	currencyEntry := widget.NewEntry()
	riskLevelEntry := widget.NewSelect([]string{"LOW", "MEDIUM", "HIGH", "CRITICAL"}, nil)

	if c != nil {
		amountEntry.SetText(fmt.Sprintf("%.2f", c.Amount))
		interestRateEntry.SetText(fmt.Sprintf("%.2f", c.InterestRate))
		purposeEntry.SetText(c.Purpose)
		issueDateEntry.SetText(c.IssueDate)
		statusEntry.SetText(c.Status)
		if c.RepaymentDate != nil {
			repaymentDateEntry.SetText(*c.RepaymentDate)
		}
		currencyEntry.SetText(c.Currency)
		if c.RiskLevel != nil {
			riskLevelEntry.SetSelected(*c.RiskLevel)
		}
	} else {
		issueDateEntry.SetText(time.Now().Format("2006-01-02"))
	}

	form.Append("Сумма", formField(amountEntry))
	form.Append("Процентная ставка", formField(interestRateEntry))
	form.Append("Цель", formField(purposeEntry))
	form.Append("Дата выдачи", formField(issueDateEntry))
	form.Append("Статус", formField(statusEntry))
	form.Append("Дата погашения", formField(repaymentDateEntry))
	form.Append("Валюта", formField(currencyEntry))
	form.Append("Уровень риска", formField(riskLevelEntry))

	form.OnSubmit = func() {
		amount, err := strconv.ParseFloat(amountEntry.Text, 64)
		if err != nil {
			showError(fmt.Errorf("неверная сумма: %w", err), mw.window)
			return
		}

		interestRate, err := strconv.ParseFloat(interestRateEntry.Text, 64)
		if err != nil {
			showError(fmt.Errorf("неверная процентная ставка: %w", err), mw.window)
			return
		}

		var repaymentDate *string
		if repaymentDateEntry.Text != "" {
			repaymentDate = &repaymentDateEntry.Text
		}

		var riskLevel *string
		if riskLevelEntry.Selected != "" {
			riskLevel = &riskLevelEntry.Selected
		}

		credit := &models.Credit{
			Amount:        amount,
			InterestRate:  interestRate,
			Purpose:       purposeEntry.Text,
			IssueDate:     issueDateEntry.Text,
			Status:        statusEntry.Text,
			RepaymentDate: repaymentDate,
			Currency:      currencyEntry.Text,
			RiskLevel:     riskLevel,
		}

		if c != nil {
			credit.ID = c.ID
			if err := mw.repos.Credit.Update(credit); err != nil {
				showError(err, mw.window)
				return
			}
		} else {
			if err := mw.repos.Credit.Create(credit); err != nil {
				showError(err, mw.window)
				return
			}
		}

		showInfo("Операция выполнена успешно", mw.window)
		onSuccess()
	}

	showFormDialog("Кредит", form, mw.window)
}
