package views

import "github.com/rivo/tview"

type ViewModel struct {
	View  tview.Primitive
	app   *tview.Application
	done  func(interface{})
	error func(error)
}

func NewViewModel(app *tview.Application) *ViewModel {
	vm := new(ViewModel)
	vm.app = app
	return vm
}

func (vm *ViewModel) SetDoneFunc(done func(data interface{})) {
	vm.done = done
}

func (vm *ViewModel) SetErrorFunc(error func(error)) {
	vm.error = error
}

func (vm *ViewModel) SetApp(app *tview.Application) {
	vm.app = app
}

func (vm *ViewModel) SetView(view tview.Primitive) {
	vm.View = view
}

func (vm *ViewModel) Done(data interface{}) {
	if vm.done != nil {
		vm.done(data)
	}
}

func (vm *ViewModel) Error(err error) {
	if vm.error != nil {
		vm.error(err)
	}
}
