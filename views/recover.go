package views

import "github.com/rivo/tview"

type RecoverDataDoneFunc struct {
	mnemonic string
	pwd      string
}

type RecoverViewModel struct {
	ViewModel

	mnemonic string
	pwd      string
}

func NewRecoverViewModel(app *tview.Application) *RecoverViewModel {

	vm := new(RecoverViewModel)
	vm.ViewModel = *NewViewModel(app)

	vm.View = vm.RecoverView()

	return vm
}

func (vm *RecoverViewModel) RecoverView() tview.Primitive {

	form := tview.NewForm()

	form.AddInputField("Mnemonic", "", 0, nil, func(text string) {
		vm.mnemonic = text
	})
	form.AddInputField("Provide a password", "", 0, nil, func(text string) {
		vm.pwd = text
	})
	form.AddButton("Recover", func() {
		if vm.mnemonic != "" && vm.pwd != "" {
			vm.Done(RecoverDataDoneFunc{vm.mnemonic, vm.pwd})
		}
	})
	form.AddButton("Cancel", func() {
		vm.Cancel()
	})

	form.SetBorder(true).SetTitle(" Insert your words ")

	flex := tview.NewFlex().
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox(), 0, 3, false).
			AddItem(form, 0, 1, true).
			AddItem(tview.NewBox(), 0, 3, false), 0, 3, true).
		AddItem(tview.NewBox(), 0, 1, false)

	return flex
}
