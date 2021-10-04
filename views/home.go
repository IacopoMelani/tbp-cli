package views

import (
	"fmt"
	"strconv"
	"time"

	"github.com/IacopoMelani/tbp-cli/api"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/rivo/tview"
)

type HomeViewModel struct {
	ViewModel

	key        *keystore.Key
	balance    uint
	inTxs      []api.TX
	outTxs     []api.TX
	pendingTxs []api.TX

	pages *tview.Pages

	textViewBalance *tview.TextView
	formSend        *tview.Form
	listTxsIn       *tview.TextView
	listTxsOut      *tview.TextView
	listTxsPending  *tview.TextView

	tempToAddress string
	tempAmount    uint
}

func NewHomeViewModel(app *tview.Application, key *keystore.Key) *HomeViewModel {

	vm := new(HomeViewModel)

	vm.app = app

	vm.key = key
	vm.inTxs = make([]api.TX, 0)
	vm.outTxs = make([]api.TX, 0)
	vm.pendingTxs = make([]api.TX, 0)

	go vm.fetch()

	vm.View = vm.HomeView()

	return vm
}

func (vm *HomeViewModel) fetch() {

	for {
		vm.fetchBalances()
		vm.fetchTxs()
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

func (vm *HomeViewModel) fetchTxs() {

	txsIn, err := api.GetTxs(vm.key, "in")
	if err == nil {
		vm.inTxs = txsIn
		setListTxsTextView(vm.listTxsIn, vm.inTxs)

	}
	txsOut, err := api.GetTxs(vm.key, "out")
	if err == nil {
		vm.outTxs = txsOut
		setListTxsTextView(vm.listTxsOut, vm.outTxs)

	}
	txsPending, err := api.GetTxs(vm.key, "pending")
	if err == nil {
		vm.pendingTxs = txsPending
		setListTxsTextView(vm.listTxsPending, vm.pendingTxs)

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
			AddItem(vm.HomeTxsView(" TX in ", vm.inTxs, func(textView *tview.TextView) {
				vm.listTxsIn = textView
			}), 0, 1, false).
			AddItem(vm.HomeTxsView(" TX out ", vm.outTxs, func(textView *tview.TextView) {
				vm.listTxsOut = textView
			}), 0, 1, false).
			AddItem(vm.HomeTxsView(" TX pending ", vm.pendingTxs, func(textView *tview.TextView) {
				vm.listTxsPending = textView
			}), 0, 1, false),
			0, 7, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(vm.HomeBalanceView(), 0, 1, false).
			AddItem(vm.HomeSendView(), 0, 8, true).
			AddItem(vm.HomeReceiveView(), 0, 1, false),
			0, 4, true)

	pages.AddPage("home", homeView, true, true)

	vm.pages = pages

	return pages
}

func (vm *HomeViewModel) HomeTxsView(title string, txs []api.TX, registerFunc func(textView *tview.TextView)) tview.Primitive {

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	setListTxsTextView(textView, txs)

	textView.SetTextAlign(tview.AlignLeft)

	textView.SetBorder(true).SetTitle(title)

	if registerFunc != nil {
		registerFunc(textView)
	}

	return textView

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

	fmt.Fprintf(textView, "You can receive [yellow]TBP[white] at this address: \n[red]%s", vm.key.Address.Hex())

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

func setListTxsTextView(textView *tview.TextView, txs []api.TX) {
	textView.Clear()
	for _, tx := range txs {
		fmt.Fprintf(textView, "⚪Block: [yellow]%v[white]\n⚫Hash: [yellow]%v[white]\n🔴From: [yellow]%v[white]\n🟠To: [yellow]%v[white]\n🟡Nonce: [yellow]%v[white]\n🟢Amount: [red]%v[yellow] TBP\n\n\n", tx.BlockHash, tx.TxHash, tx.Tx.From, tx.Tx.To, tx.Nonce, tx.Tx.Value)
	}
}
