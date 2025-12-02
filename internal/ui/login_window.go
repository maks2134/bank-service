package ui

import (
	"bank_service/internal/config"
	"bank_service/pkg/db"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type LoginWindow struct {
	app    fyne.App
	window fyne.Window
}

func NewLoginWindow() *LoginWindow {
	a := app.New()
	a.Settings().SetTheme(&CustomTheme{})

	w := a.NewWindow("Авторизация")
	w.Resize(fyne.NewSize(400, 300))

	return &LoginWindow{
		app:    a,
		window: w,
	}
}

func (lw *LoginWindow) ShowAndRun() {
	lw.buildUI()
	lw.window.ShowAndRun()
}

func (lw *LoginWindow) buildUI() {
	host := widget.NewEntry()
	host.SetText("localhost")

	port := widget.NewEntry()
	port.SetText("5432")

	user := widget.NewEntry()
	user.SetText("bankuser")

	password := widget.NewPasswordEntry()
	password.SetText("bankpassword")

	dbname := widget.NewEntry()
	dbname.SetText("bank")

	loginBtn := widget.NewButton("Войти", func() {

		cfg := &config.Config{
			DBHost:     host.Text,
			DBPort:     port.Text,
			DBUser:     user.Text,
			DBPassword: password.Text,
			DBName:     dbname.Text,
		}

		// Подключение к БД
		err := db.Init(cfg)
		if err != nil {
			dialog.ShowError(err, lw.window)
			return
		}

		mainWindow := NewMainWindow(lw.app, cfg)
		lw.window.Hide()
		mainWindow.Show()
	})

	form := container.NewVBox(
		widget.NewLabel("Параметры подключения"),
		host, port, user, password, dbname,
		loginBtn,
	)

	lw.window.SetContent(form)
}
