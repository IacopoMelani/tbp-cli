package views

import (
	"github.com/rivo/tview"
)

type ErrorViewModel struct {
	ViewModel
	err error
}

func NewErrorViewModel(app *tview.Application, err error) *ErrorViewModel {
	vm := new(ErrorViewModel)
	vm.ViewModel = *NewViewModel(app)
	vm.err = err
	vm.View = vm.ErrorView()
	return vm
}

func (vm *ErrorViewModel) ErrorView() tview.Primitive {

	modal := tview.NewModal()
	modal.AddButtons([]string{"Close"})
	modal.SetTitle(" Error! ")
	modal.SetTitleAlign(tview.AlignCenter)
	modal.SetBorder(true)
	modal.SetText(vm.err.Error())
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		vm.done(nil)
	})

	return modal
}
