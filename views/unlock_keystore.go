package views

import (
	"errors"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/rivo/tview"
)

type UnlockKeystoreViewModel struct {
	ViewModel
	keystorePath string
	password     string
	form         *tview.Form
}

func NewUnlockKeystoreViewModel(app *tview.Application) *UnlockKeystoreViewModel {
	vm := new(UnlockKeystoreViewModel)
	vm.ViewModel = *NewViewModel(app)
	vm.View = vm.UnlockKeystoreView()
	return vm
}

func (vm *UnlockKeystoreViewModel) Reset() {
	vm.keystorePath = ""
	vm.password = ""
	vm.form.GetFormItem(0).(*tview.InputField).SetText("")
	vm.form.GetFormItem(1).(*tview.InputField).SetText("")
}

func (vm *UnlockKeystoreViewModel) UnlockKeystoreView() tview.Primitive {

	form := tview.NewForm().
		AddInputField("Path to your keystore", "", 0, nil, func(text string) {
			vm.keystorePath = text
		}).
		AddPasswordField("Password", "", 0, '*', func(text string) {
			vm.password = text
		}).
		AddButton("Unlock", func() {

			if vm.keystorePath == "" || vm.password == "" {
				vm.Error(errors.New("compile required fields"))
				return
			}

			keyJson, err := ioutil.ReadFile(vm.keystorePath)
			if err != nil {
				vm.Error(err)
				return
			}

			key, err := keystore.DecryptKey(keyJson, vm.password)
			if err != nil {
				vm.Error(err)
				return
			}

			vm.done(key)
		}).
		AddButton("Clear", func() {
			vm.Reset()
		}).
		AddButton("Quit", func() {
			vm.app.Stop()
		})

	form.SetBorder(true).SetTitle(" First unlock your keystore ").SetTitleAlign(tview.AlignCenter)

	vm.form = form

	flexForm := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, 0, 1, true).
			AddItem(nil, 0, 1, false), 100, 1, true).
		AddItem(nil, 0, 1, false)

	return flexForm
}
