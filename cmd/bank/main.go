package main

import (
	"bank_service/internal/ui"
)

func main() {
	login := ui.NewLoginWindow()
	login.ShowAndRun()
}
