package views

import (
	"bytes"
	"fmt"

	"github.com/IacopoMelani/the-blockchain-pub/database"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/rivo/tview"
)

type TxReviewViewModel struct {
	ViewModel

	key       *keystore.Key
	toAddress string
	amount    uint
}

func NewTxReviewViewModel(app *tview.Application, key *keystore.Key, toAddress string, amount uint) *TxReviewViewModel {

	vm := new(TxReviewViewModel)

	vm.ViewModel = *NewViewModel(app)

	vm.key = key
	vm.toAddress = toAddress
	vm.amount = amount

	vm.View = vm.TxReviewView()

	return vm
}

func (vm *TxReviewViewModel) TxReviewView() tview.Primitive {

	buffer := bytes.NewBuffer([]byte(""))

	fmt.Fprintf(buffer, "Review your transaction:\n")
	fmt.Fprintf(buffer, "To: %v\n", vm.toAddress)
	fmt.Fprintf(buffer, "Amount: %v\n", vm.amount)
	fmt.Fprintf(buffer, "Fee: %v\n", database.TxFee)

	modal := tview.NewModal()

	modal.AddButtons([]string{"Send", "Cancel"}).
		SetText(buffer.String()).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Send" {
				vm.done(nil)
			} else {
				vm.Cancel()
			}
		})

	modal.SetTitle(" Review your transaction ").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)

	return modal
}
