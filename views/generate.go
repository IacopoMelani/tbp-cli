package views

import (
	"fmt"
	"strings"

	"github.com/IacopoMelani/tbp-cli/wallet"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
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

	mnemonicStr := fmt.Sprintf("%s\n\n", " This is your 24 words, use this to generate a new wallet with 'recover' module, do not share this words with anyone ")
	for _, str := range strings.Split(mnemonic, " ") {
		mnemonicStr += fmt.Sprintf("%s\n", str)
	}

	modal := tview.NewModal()
	modal.AddButtons([]string{"Copy to clipboard", "Close"})
	modal.SetTitleAlign(tview.AlignCenter)
	modal.SetBorder(true)
	modal.SetText(mnemonicStr)
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonIndex == 1 {
			vm.Cancel()
		} else if buttonIndex == 0 {
			clipboard.Write(clipboard.FmtText, []byte(mnemonic))
		}
	})

	flex := tview.NewFlex().
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox(), 0, 2, false).
			AddItem(modal, 0, 1, true).
			AddItem(tview.NewBox(), 0, 2, false), 0, 5, true).
		AddItem(tview.NewBox(), 0, 1, false)

	return flex
}
