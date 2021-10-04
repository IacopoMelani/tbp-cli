package views

import (
	"github.com/IacopoMelani/tbp-cli/wallet"
	"github.com/rivo/tview"
)

type GenerateViewModel struct {
	ViewModel
}

func NewGenerateViewModel(app *tview.Application) *GenerateViewModel {

	vm := new(GenerateViewModel)
	vm.app = app

	vm.View = vm.GenerateView()

	return vm
}

func (vm *GenerateViewModel) GenerateView() tview.Primitive {

	mnemonic, err := wallet.GenerateNewMnemonic()
	if err != nil {
		panic(err)
	}

	form := tview.NewForm()
	form.AddInputField("Mnemonic", mnemonic, 0, nil, func(text string) {
		form.GetFormItem(0).(*tview.InputField).SetText(mnemonic)
	}).
		AddButton("Done", func() {
			vm.Cancel()
		})

	form.SetTitle(" This is your 24 words, use this to generate a new wallet with 'recover' module, do not share this words with anyone ").SetBorder(true)

	flex := tview.NewFlex().
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox(), 0, 3, false).
			AddItem(form, 0, 1, true).
			AddItem(tview.NewBox(), 0, 3, false), 0, 5, true).
		AddItem(tview.NewBox(), 0, 1, false)

	return flex
}
