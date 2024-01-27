package pattern

import "fmt"

/*
Преймущества:
	- Горячая замена алгоритмов на лету.
	- Изолирует код и данные алгоритмов от остальных классов.
	- Уход от наследования к делегированию.
	- Реализует принцип открытости/закрытости.
Недостатки:
	- Усложняет программу за счёт дополнительных классов.
	- Клиент должен знать, в чём состоит разница между стратегиями, чтобы выбрать подходящую.
*/

// FastDelivery - стратегия быстрой доставки
type FastDelivery struct{}

// StandardDelivery - стратегия стандартной доставки
type StandardDelivery struct{}

// OrderDelivery - контекст, использующий контекст доставки
type OrderDelivery struct {
	DeliveryStrategy DeliveryStrategy
	OrderID          string
}

// Delivery - метод быстрой доставки
func (fd *FastDelivery) Delivery(orderID string) string {
	return fmt.Sprintf("fast delivery id: %s", orderID)
}

// Delivery - метод стандартной доставки
func (sd *StandardDelivery) Delivery(orderID string) string {
	return fmt.Sprintf("standard delivery id: %s", orderID)
}

// SetDeliveryStrategy - метод для установки стратегии
func (od *OrderDelivery) SetDeliveryStrategy(strategy DeliveryStrategy) {
	od.DeliveryStrategy = strategy
}

// ProcessDelivery - метод выполнения доставки с использованием текущей стратегией
func (od *OrderDelivery) ProcessDelivery() (string, error) {
	if od.DeliveryStrategy == nil {
		return "", fmt.Errorf("unset delivery strategy")
	}

	return od.DeliveryStrategy.Delivery(od.OrderID), nil
}

// DeliveryStrategy - интерфейс стратегии доставки
type DeliveryStrategy interface {
	Delivery(orderID string) string
}
