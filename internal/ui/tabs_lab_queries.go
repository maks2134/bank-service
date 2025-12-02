package ui

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/xuri/excelize/v2"
)

// Lab5Queries содержит все запросы из лабораторной работы 5
var Lab5Queries = map[string]string{
	// bankbranch
	"bankbranch_select": `SELECT full_address, contactphone FROM public.bankbranch;`,
	"bankbranch_from":   `SELECT servicezone, workingdays, workinghours FROM public.bankbranch;`,
	"bankbranch_where":  `SELECT * FROM public.bankbranch WHERE workingdays LIKE '%Сб%' OR workingdays LIKE '%Вс%';`,
	"bankbranch_order1": `SELECT * FROM public.bankbranch ORDER BY branchtype ASC;`,
	"bankbranch_order2": `SELECT * FROM public.bankbranch WHERE full_address LIKE 'г. Минск%';`,
	"bankbranch_join":   `SELECT bb.full_address, bb.branchtype, bs.fullname AS employee_name FROM public.bankbranch bb LEFT JOIN public.bankstaff bs ON bb.id = bs.branchid;`,

	// bankstaff
	"bankstaff_select": `SELECT fullname, "position", hiredate FROM public.bankstaff;`,
	"bankstaff_from":   `SELECT fullname, passport, accesslevel FROM public.bankstaff;`,
	"bankstaff_where":  `SELECT fullname, "position", accesslevel FROM public.bankstaff WHERE accesslevel = 'Высший';`,
	"bankstaff_order":  `SELECT fullname, "position" FROM public.bankstaff ORDER BY "position";`,
	"bankstaff_like":   `SELECT fullname, qualification FROM public.bankstaff WHERE qualification LIKE '%аналитик%';`,
	"bankstaff_join":   `SELECT bs.fullname, bs."position", bb.full_address FROM public.bankstaff bs JOIN public.bankbranch bb ON bs.branchid = bb.id;`,

	// client
	"client_select": `SELECT fullname, phone, email FROM public.client;`,
	"client_from":   `SELECT fullname, login, passport FROM public.client;`,
	"client_where":  `SELECT fullname, phone FROM public.client WHERE creditid IS NULL;`,
	"client_order":  `SELECT fullname, phone, email FROM public.client ORDER BY fullname ASC;`,
	"client_like":   `SELECT fullname, email FROM public.client WHERE email LIKE '%@mail.com';`,
	"client_join":   `SELECT c.fullname, cr.amount, cr.interestrate, cr.purpose FROM public.client c JOIN public.credit cr ON c.creditid = cr.id;`,

	// credit
	"credit_select": `SELECT amount, interestrate, currency, purpose FROM public.credit;`,
	"credit_from":   `SELECT purpose, status, issuedate, repaymentdate FROM public.credit;`,
	"credit_where":  `SELECT * FROM public.credit WHERE status = 'Закрыт';`,
	"credit_order":  `SELECT amount, interestrate, purpose FROM public.credit ORDER BY interestrate DESC;`,
	"credit_like":   `SELECT * FROM public.credit WHERE purpose LIKE 'Потребительский%';`,
	"credit_join":   `SELECT cr.amount, cr.purpose, bs.fullname AS staff_member FROM public.credit cr JOIN public.creditbankstaff cbs ON cr.id = cbs.creditid JOIN public.bankstaff bs ON cbs.staffid = bs.id;`,

	// bankaccount
	"bankaccount_select": `SELECT accountnumber, balance, currency FROM public.bankaccount;`,
	"bankaccount_from":   `SELECT accountnumber, accounttype, accountstatus, opendate FROM public.bankaccount;`,
	"bankaccount_where":  `SELECT * FROM public.bankaccount WHERE accountstatus = 'BLOCKED' OR accountstatus = 'SUSPENDED';`,
	"bankaccount_order":  `SELECT accountnumber, balance, currency FROM public.bankaccount ORDER BY balance DESC;`,
	"bankaccount_like":   `SELECT * FROM public.bankaccount WHERE accounttype LIKE 'Сберегательный%';`,
	"bankaccount_join":   `SELECT ba.accountnumber, ba.balance, c.fullname AS client_name FROM public.bankaccount ba JOIN public.client c ON ba.clientid = c.id;`,

	// transaction
	"transaction_select": `SELECT amount, operationtype, currency FROM public.transaction;`,
	"transaction_from":   `SELECT operationdate, operationtype, operationstatus, purpose FROM public.transaction;`,
	"transaction_where":  `SELECT * FROM public.transaction WHERE operationstatus = 'Отклонено';`,
	"transaction_order":  `SELECT amount, operationdate, operationtype FROM public.transaction ORDER BY operationdate DESC;`,
	"transaction_like":   `SELECT * FROM public.transaction WHERE purpose LIKE '%Перевод%';`,
	"transaction_join":   `SELECT t.amount, t.operationdate, t.purpose, c.fullname AS client_name FROM public.transaction t JOIN public.bankaccount ba ON t.accountid = ba.id JOIN public.client c ON ba.clientid = c.id;`,
}

// Lab6Queries содержит все запросы из лабораторной работы 6
var Lab6Queries = map[string]string{
	"bankbranch_group": `SELECT bb.branchtype, AVG(staff_count) as avg_staff
FROM (
    SELECT branchid, COUNT(*) as staff_count
    FROM public.bankstaff
    GROUP BY branchid
) staff_counts
JOIN public.bankbranch bb ON staff_counts.branchid = bb.id
GROUP BY bb.branchtype
HAVING AVG(staff_count) > (
    SELECT AVG(cnt)
    FROM (
        SELECT COUNT(*) as cnt
        FROM public.bankstaff
        GROUP BY branchid
    ) sub
);`,

	"bankstaff_group": `SELECT "position", AVG(hiredate::date) as avg_hire_date
FROM public.bankstaff
GROUP BY "position"
HAVING AVG(hiredate::date) > (
    SELECT AVG(hiredate::date)
    FROM public.bankstaff
);`,

	"client_group": `SELECT c.clienttype, SUM(cr.amount) as total_credit_amount
FROM public.client c
JOIN public.credit cr ON c.creditid = cr.id
GROUP BY c.clienttype
HAVING SUM(cr.amount) > (
    SELECT AVG(amount)
    FROM public.credit
);`,

	"credit_group": `SELECT purpose, 
    COUNT(*) as credit_count,
    MAX(amount) - MIN(amount) as amount_diff
FROM public.credit
WHERE purpose IN (
    SELECT purpose
    FROM public.credit
    GROUP BY purpose
    HAVING COUNT(DISTINCT currency) > 1
)
GROUP BY purpose
HAVING AVG(interestrate) > (
    SELECT AVG(interestrate)
    FROM public.credit
);`,

	"bankaccount_group": `SELECT ba.accounttype, AVG(t.amount) as avg_transaction_amount
FROM public.bankaccount ba
JOIN public.transaction t ON ba.id = t.accountid
GROUP BY ba.accounttype
HAVING AVG(t.amount) > (
    SELECT AVG(amount)
    FROM public.transaction
);`,

	"transaction_group": `SELECT 
    TO_CHAR(operationdate, 'YYYY-MM') as month,
    COUNT(*) as total_transactions,
    SUM(amount) as total_amount,
    ROUND(100.0 * COUNT(CASE WHEN operationstatus = 'Отклонено' THEN 1 END)::numeric / COUNT(*)::numeric, 2) as rejected_percent,
    AVG(amount) as avg_amount
FROM public.transaction
GROUP BY TO_CHAR(operationdate, 'YYYY-MM')
ORDER BY month;`,
}

type labQueryDefinition struct {
	Key         string
	Label       string
	Description string
}

type labQueryMenuOption struct {
	Label       string
	Query       string
	Description string
}

func (mw *MainWindow) createLabQueriesTab() *container.TabItem {
	var currentData [][]interface{}
	var currentColumns []string
	var filteredData [][]interface{}

	// Result table
	resultTable := widget.NewTable(
		func() (int, int) { return 0, 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, obj fyne.CanvasObject) {},
	)

	// Search and filter
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Поиск по результатам...")

	filterColumnSelect := widget.NewSelect([]string{"Все колонки"}, nil)
	filterColumnSelect.SetSelected("Все колонки")

	var resultContainerBorder fyne.CanvasObject
	var resultTableContainer *fyne.Container
	var searchPanel fyne.CanvasObject

	updateResultTable := func() {
		data := filteredData
		if len(data) == 0 {
			data = currentData
		}
		if len(data) == 0 || len(currentColumns) == 0 {
			return
		}
		newTable := createTableWidget(data, currentColumns)
		// Set column widths
		for i := 0; i < len(currentColumns); i++ {
			newTable.SetColumnWidth(i, 150)
		}
		resultTable = newTable
		// Update table container
		if resultTableContainer != nil {
			resultTableContainer.Objects = []fyne.CanvasObject{resultTable}
			resultTableContainer.Refresh()
		}
	}

	searchEntry.OnChanged = func(text string) {
		if len(currentData) == 0 {
			return
		}
		if text == "" {
			filteredData = nil
			updateResultTable()
			return
		}

		colIdx := -1
		if filterColumnSelect.Selected != "Все колонки" {
			for i, col := range currentColumns {
				if col == filterColumnSelect.Selected {
					colIdx = i
					break
				}
			}
		}

		filteredData = [][]interface{}{}
		textLower := strings.ToLower(text)
		for _, row := range currentData {
			match := false
			if colIdx >= 0 {
				if colIdx < len(row) && row[colIdx] != nil {
					if strings.Contains(strings.ToLower(fmt.Sprintf("%v", row[colIdx])), textLower) {
						match = true
					}
				}
			} else {
				for _, val := range row {
					if val != nil && strings.Contains(strings.ToLower(fmt.Sprintf("%v", val)), textLower) {
						match = true
						break
					}
				}
			}
			if match {
				filteredData = append(filteredData, row)
			}
		}
		updateResultTable()
	}

	executeQuery := func(query string, queryName string) {
		data, columns, err := mw.utils.QueryManager.ExecuteQuery(query)
		if err != nil {
			showError(err, mw.window)
			return
		}

		currentData = data
		currentColumns = columns
		filteredData = nil
		searchEntry.SetText("")

		// Update filter column select
		colOptions := []string{"Все колонки"}
		colOptions = append(colOptions, columns...)
		filterColumnSelect.Options = colOptions
		filterColumnSelect.SetSelected("Все колонки")

		// Create new table with data
		resultTable = createTableWidget(data, columns)
		// Set column widths
		for i := 0; i < len(columns); i++ {
			resultTable.SetColumnWidth(i, 150)
		}

		// Update container
		if resultTableContainer != nil {
			resultTableContainer.Objects = []fyne.CanvasObject{resultTable}
			resultTableContainer.Refresh()
		}

		showInfo(fmt.Sprintf("Запрос '%s' выполнен. Найдено записей: %d", queryName, len(data)), mw.window)
	}

	menuOptions := []labQueryMenuOption{}
	addButtons := func(defs []labQueryDefinition, queryMap map[string]string, groupLabel string) []fyne.CanvasObject {
		buttons := make([]fyne.CanvasObject, 0, len(defs))
		for _, definition := range defs {
			definition := definition
			queryStr, ok := queryMap[definition.Key]
			if !ok {
				continue
			}
			label := definition.Label
			description := definition.Description
			queryValue := queryStr

			buttons = append(buttons, widget.NewButton(label, func() {
				executeQuery(queryValue, description)
			}))

			menuOptions = append(menuOptions, labQueryMenuOption{
				Label:       fmt.Sprintf("%s · %s", groupLabel, label),
				Query:       queryValue,
				Description: description,
			})
		}
		return buttons
	}

	lab5Title := widget.NewRichTextFromMarkdown("## Оперативные отчеты\n### Стандартные выборки")

	lab5BankBranchDefs := []labQueryDefinition{
		{"bankbranch_select", "Адреса и телефоны отделений", "Адреса и телефоны отделений"},
		{"bankbranch_from", "Зоны обслуживания и режим работы", "Реквизиты отделений"},
		{"bankbranch_where", "Отделения, работающие в выходные", "Отделения в выходные"},
		{"bankbranch_order1", "Перечень отделений по типу", "Отделения по типу"},
		{"bankbranch_order2", "Отделения в Минске", "Отделения в Минске"},
		{"bankbranch_join", "Отделения и закрепленные сотрудники", "Отделения и сотрудники"},
	}
	lab5BankBranchGroup := widget.NewCard("Отчеты по отделениям", "", container.NewVBox(addButtons(lab5BankBranchDefs, Lab5Queries, "Отделения")...))

	lab5BankStaffDefs := []labQueryDefinition{
		{"bankstaff_select", "Список сотрудников с должностями", "Сотрудники - основные данные"},
		{"bankstaff_from", "Паспорта и уровни доступа", "Идентификационные данные"},
		{"bankstaff_where", "Сотрудники с максимальным доступом", "Высший уровень доступа"},
		{"bankstaff_order", "Персонал по должностям", "Сортировка по должности"},
		{"bankstaff_like", "Специалисты-аналитики", "Сотрудники-аналитики"},
		{"bankstaff_join", "Сотрудники и их отделения", "Сотрудники и отделения"},
	}
	lab5BankStaffGroup := widget.NewCard("Отчеты по персоналу", "", container.NewVBox(addButtons(lab5BankStaffDefs, Lab5Queries, "Персонал")...))

	lab5ClientDefs := []labQueryDefinition{
		{"client_select", "Контакты клиентов", "Клиенты - контакты"},
		{"client_from", "Документы и логины клиентов", "Идентификационные данные клиентов"},
		{"client_where", "Клиенты без действующих кредитов", "Клиенты без кредита"},
		{"client_order", "Клиенты по алфавиту", "Клиенты по ФИО"},
		{"client_like", "Клиенты c почтой на mail.com", "Клиенты с mail.com"},
		{"client_join", "Клиенты и параметры кредитов", "Клиенты и кредиты"},
	}
	lab5ClientGroup := widget.NewCard("Отчеты по клиентам", "", container.NewVBox(addButtons(lab5ClientDefs, Lab5Queries, "Клиенты")...))

	lab5CreditDefs := []labQueryDefinition{
		{"credit_select", "Размеры кредитов и валюта", "Кредиты - финансовые данные"},
		{"credit_from", "Цель, статус и даты кредитов", "Параметры кредитов"},
		{"credit_where", "Закрытые кредиты", "Закрытые кредиты"},
		{"credit_order", "Кредиты по процентной ставке", "Кредиты по ставке"},
		{"credit_like", "Потребительские кредиты", "Потребительские кредиты"},
		{"credit_join", "Кредиты и ответственные сотрудники", "Кредиты и сотрудники"},
	}
	lab5CreditGroup := widget.NewCard("Отчеты по кредитам", "", container.NewVBox(addButtons(lab5CreditDefs, Lab5Queries, "Кредиты")...))

	lab5BankAccountDefs := []labQueryDefinition{
		{"bankaccount_select", "Счета: номер, баланс, валюта", "Счета - основные данные"},
		{"bankaccount_from", "Статусы счетов", "Статус счетов"},
		{"bankaccount_where", "Заблокированные и приостановленные счета", "Заблокированные счета"},
		{"bankaccount_order", "Клиенты с крупными остатками", "Счета по балансу"},
		{"bankaccount_like", "Сберегательные счета", "Сберегательные счета"},
		{"bankaccount_join", "Счета и владельцы", "Счета и клиенты"},
	}
	lab5BankAccountGroup := widget.NewCard("Отчеты по счетам", "", container.NewVBox(addButtons(lab5BankAccountDefs, Lab5Queries, "Счета")...))

	lab5TransactionDefs := []labQueryDefinition{
		{"transaction_select", "Основные параметры транзакций", "Транзакции - основные данные"},
		{"transaction_from", "Типы, статусы и назначения", "Статус транзакций"},
		{"transaction_where", "Отклоненные транзакции", "Отклоненные транзакции"},
		{"transaction_order", "Хронология операций", "Транзакции по дате"},
		{"transaction_like", "Переводы средств", "Транзакции-переводы"},
		{"transaction_join", "Транзакции с привязкой к клиентам", "Транзакции и клиенты"},
	}
	lab5TransactionGroup := widget.NewCard("Отчеты по транзакциям", "", container.NewVBox(addButtons(lab5TransactionDefs, Lab5Queries, "Транзакции")...))

	lab5Scroll := container.NewScroll(container.NewVBox(
		lab5Title,
		lab5BankBranchGroup,
		lab5BankStaffGroup,
		lab5ClientGroup,
		lab5CreditGroup,
		lab5BankAccountGroup,
		lab5TransactionGroup,
	))

	// Lab 6 buttons
	lab6Title := widget.NewRichTextFromMarkdown("## Аналитика и агрегаты\n### Сводные отчеты")

	lab6Defs := []labQueryDefinition{
		{"bankbranch_group", "Среднее количество сотрудников по типу отделений", "Отделения - среднее количество сотрудников"},
		{"bankstaff_group", "Должности с более поздним наймом", "Сотрудники - средняя дата приема"},
		{"client_group", "Типы клиентов с суммой кредитов выше среднего", "Клиенты - сумма кредитов"},
		{"credit_group", "Цели кредитов с повышенными ставками", "Кредиты - анализ по целям"},
		{"bankaccount_group", "Типы счетов с высокой активностью транзакций", "Счета - средняя сумма транзакций"},
		{"transaction_group", "Ежемесячная статистика транзакций", "Транзакции - статистика по месяцам"},
	}
	lab6Group := widget.NewCard("Стратегическая аналитика", "", container.NewVBox(addButtons(lab6Defs, Lab6Queries, "Аналитика")...))

	var menuSelect *widget.Select
	menuLabels := make([]string, len(menuOptions))
	for i, opt := range menuOptions {
		menuLabels[i] = opt.Label
	}
	menuSelect = widget.NewSelect(menuLabels, func(value string) {
		for _, opt := range menuOptions {
			if opt.Label == value {
				executeQuery(opt.Query, opt.Description)
				return
			}
		}
	})
	menuSelect.PlaceHolder = "Быстрый выбор отчета"

	lab6Scroll := container.NewScroll(container.NewVBox(
		lab6Title,
		lab6Group,
	))

	// Left panel with tabs for Lab5 and Lab6
	labTabs := container.NewAppTabs(
		container.NewTabItem("Оперативные отчеты", lab5Scroll),
		container.NewTabItem("Бизнес-аналитика", lab6Scroll),
	)

	// Search and filter panel with menu
	searchPanel = container.NewVBox(
		widget.NewLabel("Быстрое меню отчетов:"),
		menuSelect,
		widget.NewSeparator(),
		widget.NewLabel("Поиск и фильтрация:"),
		searchEntry,
		filterColumnSelect,
		widget.NewButton("Очистить фильтр", func() {
			searchEntry.SetText("")
			filteredData = nil
			updateResultTable()
		}),
		widget.NewButton("Экспорт в Excel", func() {
			if len(currentData) == 0 {
				showInfo("Нет данных для экспорта", mw.window)
				return
			}
			data := filteredData
			if len(data) == 0 {
				data = currentData
			}
			// Export filtered data
			filename := fmt.Sprintf("lab_query_%s.xlsx", time.Now().Format("20060102_150405"))
			if err := mw.exportDataToExcel(data, currentColumns, filename); err != nil {
				showError(err, mw.window)
			} else {
				showInfo(fmt.Sprintf("Данные экспортированы в %s", filename), mw.window)
			}
		}),
	)

	// Right panel with results
	resultTableContainer = container.NewMax(resultTable)
	resultContainerBorder = container.NewBorder(
		searchPanel,
		nil,
		nil,
		nil,
		resultTableContainer,
	)

	// Main split
	split := container.NewHSplit(labTabs, resultContainerBorder)
	split.SetOffset(0.35)

	// Автоматически выполнить первый запрос при открытии вкладки
	if len(menuOptions) > 0 {
		first := menuOptions[0]
		executeQuery(first.Query, first.Description)
		if menuSelect != nil {
			menuSelect.SetSelected(first.Label)
		}
	}

	return container.NewTabItem("Отчеты и аналитика", split)
}

func (mw *MainWindow) exportDataToExcel(data [][]interface{}, columns []string, filename string) error {
	if len(data) == 0 {
		return fmt.Errorf("нет данных для экспорта")
	}

	// Use excelize directly to export data
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sheet1"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Write headers
	for i, col := range columns {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, col)
	}

	// Write data
	for rowNum, row := range data {
		for colNum, val := range row {
			if colNum < len(columns) {
				cell := fmt.Sprintf("%c%d", 'A'+colNum, rowNum+2)
				if val != nil {
					f.SetCellValue(sheetName, cell, val)
				}
			}
		}
	}

	f.SetActiveSheet(index)
	return f.SaveAs(filename)
}
