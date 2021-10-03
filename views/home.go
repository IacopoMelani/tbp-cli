package views

import (
	"fmt"
	"strconv"
	"time"

	"github.com/IacopoMelani/tbp-cli/api"
	"github.com/IacopoMelani/the-blockchain-pub/database"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/rivo/tview"
)

type HomeViewModel struct {
	ViewModel
	balance uint
	inTx    []database.SignedTx
	outTx   []database.SignedTx
	key     *keystore.Key

	pages *tview.Pages

	textViewBalance *tview.TextView
	formSend        *tview.Form

	tempToAddress string
	tempAmount    uint
}

func NewHomeViewModel(app *tview.Application, key *keystore.Key) *HomeViewModel {

	vm := new(HomeViewModel)

	vm.app = app
	vm.View = vm.HomeView()

	vm.key = key
	vm.inTx = make([]database.SignedTx, 0)
	vm.outTx = make([]database.SignedTx, 0)

	go vm.fetch()

	return vm
}

func (vm *HomeViewModel) fetch() {

	for {
		vm.fetchBalances()
		vm.app.Draw()
		time.Sleep(time.Second * 5)
	}
}

func (vm *HomeViewModel) fetchBalances() {
	balance, err := api.GetBalance(vm.key.Address)
	if err == nil {
		vm.balance = balance
		if vm.textViewBalance != nil {
			setBalaceTextViewBalance(vm.textViewBalance, vm.balance)
		}
	}
}

func (vm *HomeViewModel) resetTempTx() {
	vm.tempToAddress = ""
	vm.tempAmount = 0
	vm.formSend.GetFormItem(0).(*tview.InputField).SetText("")
	vm.formSend.GetFormItem(1).(*tview.InputField).SetText("")
}

func (vm *HomeViewModel) HomeView() tview.Primitive {

	pages := tview.NewPages()

	homeView := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(tview.NewBox().SetBorder(true).SetTitle(" TX in "), 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle(" TX out "), 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle(" TX pending "), 0, 1, false),
			0, 5, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(vm.HomeBalanceView(), 0, 1, false).
			AddItem(vm.HomeSendView(), 0, 8, true).
			AddItem(vm.HomeReceiveView(), 0, 1, false),
			0, 4, true)

	pages.AddPage("home", homeView, true, true)

	vm.pages = pages

	return pages
}

func (vm *HomeViewModel) HomeBalanceView() tview.Primitive {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	setBalaceTextViewBalance(textView, vm.balance)

	textView.SetTextAlign(tview.AlignCenter)

	textView.SetBorder(true).SetTitle(" Balance ")

	vm.textViewBalance = textView

	return textView
}

func (vm *HomeViewModel) HomeReceiveView() tview.Primitive {

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	fmt.Fprintf(textView, "You can receive [yellow]TBP[white] at this address: \n[red]%s", "0x50543e830590fd03a0301faa0164d731f0e2ff7d")

	textView.SetTextAlign(tview.AlignCenter)

	textView.SetBorder(true).SetTitle(" Receive TBP ")

	return textView
}

func (vm *HomeViewModel) HomeSendView() tview.Primitive {

	form := tview.NewForm().
		AddInputField("Address", "", 0, nil, func(text string) {
			vm.tempToAddress = text
		}).
		AddInputField("Amount", "", 0, tview.InputFieldInteger, func(text string) {
			amount, err := strconv.Atoi(text)
			if err != nil {
				return
			}
			vm.tempAmount = uint(amount)
		}).
		AddButton("Send", func() {
			if vm.tempToAddress == "" || vm.tempAmount == 0 {
				return
			}
			reviewTxView := NewTxReviewViewModel(vm.app, vm.key, vm.tempToAddress, vm.tempAmount)
			reviewTxView.SetDoneFunc(func(data interface{}) {
				if err := api.SendTx(vm.key, vm.tempToAddress, vm.tempAmount); err != nil {
					errorViewModel := NewErrorViewModel(vm.app, err)
					vm.pages.AddAndSwitchToPage("error", errorViewModel.View, true)

					errorViewModel.SetDoneFunc(func(data interface{}) {
						vm.pages.SwitchToPage("home")
						vm.pages.RemovePage("error")
					})
				}
				vm.resetTempTx()
				vm.pages.RemovePage("review")
			})
			reviewTxView.SetCancelFunc(func() {
				vm.resetTempTx()
				vm.pages.RemovePage("review")

			})
			vm.pages.AddAndSwitchToPage("review", reviewTxView.View, true)
		}).
		AddButton("Cancel", func() {
			vm.resetTempTx()
		}).SetButtonsAlign(tview.AlignCenter)

	vm.formSend = form

	form.SetBorder(true).SetTitle(" Send TBP ")

	return form
}

func setBalaceTextViewBalance(textView *tview.TextView, balance uint) {
	textView.Clear()
	fmt.Fprintf(textView, "Your current balance is: %d [yellow]TBP", balance)
}
