package views

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/IacopoMelani/tbp-cli/wallet"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/rivo/tview"
)

type UnlockKeystoreViewModel struct {
	ViewModel
	keystorePath string
	password     string

	form  *tview.Form
	pages *tview.Pages
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
		AddButton("Recover", func() {
			recover := NewRecoverViewModel(vm.app)
			vm.pages.AddAndSwitchToPage("recover", recover.View, true)
			recover.SetCancelFunc(func() {
				vm.pages.SwitchToPage("unlock")
				vm.pages.RemovePage("recover")
			})
			recover.SetDoneFunc(func(d interface{}) {
				data := d.(RecoverDataDoneFunc)
				filename := fmt.Sprintf("key-%d", time.Now().Unix())
				if err := wallet.GenerateKeyFromMnemonicAndStore(filename, data.pwd, data.mnemonic); err != nil {
					vm.Error(err)
					return
				}
				vm.pages.SwitchToPage("unlock")
				vm.pages.RemovePage("recover")
			})
		}).
		AddButton("New", func() {
			generate := NewGenerateViewModel(vm.app)
			vm.pages.AddAndSwitchToPage("generate", generate.View, true)
			generate.SetCancelFunc(func() {
				vm.pages.SwitchToPage("unlock")
				vm.pages.RemovePage("generate")
			})
		})

	form.SetBorder(true).SetTitle(" First unlock your keystore ").SetTitleAlign(tview.AlignCenter)

	vm.form = form

	pages := tview.NewPages()

	flexForm := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, 0, 1, true).
			AddItem(nil, 0, 1, false), 100, 1, true).
		AddItem(nil, 0, 1, false)

	pages.AddPage("unlock", flexForm, true, true)

	vm.pages = pages

	return pages
}
