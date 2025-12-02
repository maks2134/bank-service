package ui

import (
	"fmt"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (mw *MainWindow) createBackupTab() *container.TabItem {
	statusLabel := widget.NewLabel("Готово к работе")

	backupTableBtn := widget.NewButton("Резервная копия таблицы", func() {
		tables, err := mw.utils.TableManager.GetTables()
		if err != nil {
			showError(err, mw.window)
			return
		}

		tableSelect := widget.NewSelect(tables, nil)
		form := widget.NewForm()
		form.Append("Выберите таблицу", tableSelect)

		form.OnSubmit = func() {
			if tableSelect.Selected == "" {
				showInfo("Выберите таблицу", mw.window)
				return
			}
			filename, err := mw.utils.BackupManager.BackupTable(tableSelect.Selected)
			if err != nil {
				showError(err, mw.window)
				return
			}
			statusLabel.SetText(fmt.Sprintf("Создана резервная копия: %s", filename))
			showInfo(fmt.Sprintf("Резервная копия создана: %s", filename), mw.window)
		}

		dialog.ShowCustom("Резервная копия таблицы", "Отмена", form, mw.window)
	})

	backupDBBtn := widget.NewButton("Резервная копия БД", func() {
		dialog.ShowConfirm("Подтверждение", "Создать резервную копию всей базы данных?", func(ok bool) {
			if ok {
				filename, err := mw.utils.BackupManager.BackupDatabase()
				if err != nil {
					showError(err, mw.window)
					return
				}
				statusLabel.SetText(fmt.Sprintf("Создана резервная копия БД: %s", filename))
				showInfo(fmt.Sprintf("Резервная копия БД создана: %s", filename), mw.window)
			}
		}, mw.window)
	})

	restoreTableBtn := widget.NewButton("Восстановить таблицу", func() {
		fileEntry := widget.NewEntry()
		fileEntry.SetPlaceHolder("Путь к файлу резервной копии")

		tables, err := mw.utils.TableManager.GetTables()
		if err != nil {
			showError(err, mw.window)
			return
		}

		tableSelect := widget.NewSelect(tables, nil)
		form := widget.NewForm()
		form.Append("Файл резервной копии", fileEntry)
		form.Append("Выберите таблицу для восстановления", tableSelect)

		form.OnSubmit = func() {
			if fileEntry.Text == "" {
				showInfo("Введите путь к файлу", mw.window)
				return
			}
			if tableSelect.Selected == "" {
				showInfo("Выберите таблицу", mw.window)
				return
			}
			if err := mw.utils.BackupManager.RestoreTable(fileEntry.Text, tableSelect.Selected); err != nil {
				showError(err, mw.window)
				return
			}
			statusLabel.SetText(fmt.Sprintf("Таблица %s восстановлена", tableSelect.Selected))
			showInfo("Таблица восстановлена", mw.window)
		}

		dialog.ShowCustom("Восстановить таблицу", "Отмена", form, mw.window)
	})

	content := container.NewVBox(
		widget.NewLabel("Резервное копирование и восстановление"),
		statusLabel,
		backupTableBtn,
		backupDBBtn,
		restoreTableBtn,
	)

	return container.NewTabItem("Резервное копирование", content)
}
