package main

import (
	"github.com/IacopoMelani/tbp-cli/views"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/rivo/tview"
)

const (
	pageUnlockKeystore = "unlockKeystore"
	pageError          = "error"
	pageHome           = "home"
)

func main() {

	app := tview.NewApplication()

	pages := tview.NewPages()

	unlockKeystoreViewModel := views.NewUnlockKeystoreViewModel(app)

	unlockKeystoreViewModel.SetDoneFunc(func(data interface{}) {

		key := data.(*keystore.Key)
		homeViewModel := views.NewHomeViewModel(app, key)

		pages.AddAndSwitchToPage(pageHome, homeViewModel.View, true)
	})
	unlockKeystoreViewModel.SetErrorFunc(func(err error) {

		errorViewModel := views.NewErrorViewModel(app, err)
		pages.AddAndSwitchToPage(pageError, errorViewModel.View, true)

		errorViewModel.SetDoneFunc(func(data interface{}) {
			unlockKeystoreViewModel.Reset()
			pages.SwitchToPage(pageUnlockKeystore)
		})
	})

	pages.AddPage(pageUnlockKeystore, unlockKeystoreViewModel.View, true, true)

	showView(app, pages)
}

func showView(app *tview.Application, view tview.Primitive) {
	if err := app.SetRoot(view, true).Run(); err != nil {
		panic(err)
	}
}
