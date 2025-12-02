package ui

import (
	"bank_service/internal/config"
	"bank_service/internal/repository"
	"bank_service/internal/utils"
	"bank_service/pkg/db"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type MainWindow struct {
	app            fyne.App
	window         fyne.Window
	repos          *Repositories
	utils          *Utils
	currentContent *container.AppTabs
}

type Repositories struct {
	BankAccount     *repository.BankAccountRepository
	BankBranch      *repository.BankBranchRepository
	BankStaff       *repository.BankStaffRepository
	Client          *repository.ClientRepository
	Credit          *repository.CreditRepository
	CreditBankStaff *repository.CreditBankStaffRepository
	Transaction     *repository.TransactionRepository
}

type Utils struct {
	TableManager  *utils.TableManager
	BackupManager *utils.BackupManager
	QueryManager  *utils.QueryManager
	ExportManager *utils.ExportManager
}

func NewMainWindow(a fyne.App, cfg *config.Config) *MainWindow {
	window := a.NewWindow("Банковская система")
	window.Resize(fyne.NewSize(1200, 800))

	repos := &Repositories{
		BankAccount:     repository.NewBankAccountRepository(),
		BankBranch:      repository.NewBankBranchRepository(),
		BankStaff:       repository.NewBankStaffRepository(),
		Client:          repository.NewClientRepository(),
		Credit:          repository.NewCreditRepository(),
		CreditBankStaff: repository.NewCreditBankStaffRepository(),
		Transaction:     repository.NewTransactionRepository(),
	}

	utils := &Utils{
		TableManager:  utils.NewTableManager(),
		BackupManager: utils.NewBackupManager(),
		QueryManager:  utils.NewQueryManager(),
		ExportManager: utils.NewExportManager(),
	}

	mw := &MainWindow{
		app:    a,
		window: window,
		repos:  repos,
		utils:  utils,
	}

	mw.setupUI()
	return mw
}

func (mw *MainWindow) Show() {
	mw.window.Show()
}

func (mw *MainWindow) setupUI() {
	tabs := container.NewAppTabs(
		mw.createBankAccountTab(),
		mw.createBankBranchTab(),
		mw.createBankStaffTab(),
		mw.createClientTab(),
		mw.createCreditTab(),
		mw.createCreditBankStaffTab(),
		mw.createTransactionTab(),
		mw.createTableManagementTab(),
		mw.createBackupTab(),
		mw.createQueriesTab(),
		mw.createLabQueriesTab(),
	)

	mw.currentContent = tabs
	mw.window.SetContent(tabs)
}

func (mw *MainWindow) ShowAndRun() {
	mw.window.ShowAndRun()
	db.Close()
}
