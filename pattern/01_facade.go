package pattern

import "fmt"

// Преймуществом паттерна фасад, является сокрытие сложной реализацией за простым интерфейсом
// Недостаток данного паттерна - Фасад может превраться в "божественный объект"
// в качестве примера можно привети оплату заказа, конвертацию видео, регистрацию на сайте
// всё это может реализовываться с помощью паттерна фасад, который будет скрывать за собой множество
// объектов, учавствующих в процессе

// Payment - структура отвечающая за оплату
type Payment struct{}

// PaymentHistory - структура отвечающая историю оплат
type PaymentHistory struct{}

// Mailer - структура отвечающая рассылку писем на почту
type Mailer struct{}

// Product - структура продукта
type Product struct{}

// OrderFacade - структура реализующая паттерн фасад
type OrderFacade struct {
	pmt        Payment
	pmtHistory PaymentHistory
	mailer     Mailer
	prod       Product
}

func (pmt *Payment) MakePayment() {
	fmt.Println("Make a new payment")
}

func (pmtH *PaymentHistory) SavePayment() {
	fmt.Println("save new payment")
}

func (m *Mailer) SendNotification() {
	fmt.Println("send notification about payment")
}

func (prd *Product) ChangeCount() {
	fmt.Println("Change count products in data base")
}

// MakeOrder - метод вызывающий необходимые для заказа методы други структур
func (ord *OrderFacade) MakeOrder() {
	ord.pmt.MakePayment()
	ord.pmtHistory.SavePayment()
	ord.mailer.SendNotification()
	ord.prod.ChangeCount()
}
