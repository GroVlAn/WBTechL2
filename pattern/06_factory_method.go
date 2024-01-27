package pattern

import "fmt"

/*
Преймущества:
	- Избавляет класс от привязки к конкретным классам продуктов.
	- Выделяет код производства продуктов в одно место, упрощая поддержку кода.
	- Упрощает добавление новых продуктов в программу.
	- Реализует принцип открытости/закрытости.
Недостатки:
	- Может привести к созданию больших параллельных иерархий классов,
	так как для каждого класса продукта надо создать свой подкласс создателя.
*/

// SberPayment - структура оплаты через сбербанк
type SberPayment struct{}

// TinkoffPayment - структура оплаты через сбербанк
type TinkoffPayment struct{}

// VTBPayment - структура оплаты через ВТБ
type VTBPayment struct{}

// Pay - оплата через сбербанк
func (sb *SberPayment) Pay(rubles, kopecks int) string {
	return fmt.Sprintf("Pay with seberbank card: %d.%d", rubles, kopecks)
}

// Pay - оплата через тинькофф
func (tp *TinkoffPayment) Pay(rubles, kopecks int) string {
	return fmt.Sprintf("Pay with tinkoff card: %d.%d", rubles, kopecks)
}

// Pay - оплата через втб
func (vtbp *VTBPayment) Pay(rubles, kopecks int) string {
	return fmt.Sprintf("Pay with vtb card: %d.%d", rubles, kopecks)
}

// CreatePayment - функция для создания конкретного способа оплаты
func CreatePayment(name string) (Paymenter, error) {
	switch name {
	case "sber":
		return &SberPayment{}, nil
	case "tinkoff":
		return &TinkoffPayment{}, nil
	case "vtb":
		return &VTBPayment{}, nil
	default:
		return nil, fmt.Errorf("wrong payment type passed")
	}
}

// Paymenter - интерфейс для структур оплаты
type Paymenter interface {
	Pay(rubles, kopecks int) string
}
