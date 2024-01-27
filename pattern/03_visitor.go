package pattern

/*
Преймущества:
	Упрощает добавление новой логики, для работы со сложными структурами
	Объединяет родственные операции водной структуре
	Может накапливать состояние при обходе структуры элементов
Недостатки:
	Нерелевантен, если иерархия элементов часто меняется
	Может привести к нарушению инкапсуляцию элементов
*/

type Order struct {
	TrackNumber string
	Entry       string
	DateCreated string
}

type Delivery struct {
	Name  string
	Phone string
	Email string
}

// Convertor - интерфейс посетителя
type Convertor interface {
	ConvertorOrder(ordRep *OrderRep) *Order
	ConvertorDelivery(delRep *DeliveryRep) *Delivery
}

// Getter - интерфейс объекта
type Getter interface {
	Get() map[string]string
	Accept(v Convertor)
}

type OrderRep struct{}

type DeliveryRep struct{}

type MapToStructConvertor struct{}

func (ordRep *OrderRep) Get() map[string]string {
	return map[string]string{
		"track_number": "WBILMTESTTRACK",
		"entry":        "WBIL",
		"date_created": "2021-11-26T06:22:19Z",
	}
}

func (ordRep *OrderRep) Accept(v Convertor) {
	v.ConvertorOrder(ordRep)
}

func (delRep *DeliveryRep) Get() map[string]string {
	return map[string]string{
		"name":  "Test Testov",
		"phone": "+9720000000",
		"email": "test@gmail.com",
	}
}

func (delRep *DeliveryRep) Accept(v Convertor) {
	v.ConvertorDelivery(delRep)
}

func (mtsc *MapToStructConvertor) ConvertorOrder(ordRep *OrderRep) *Order {
	ordMap := ordRep.Get()
	ord := new(Order)
	ord.TrackNumber = ordMap["track_number"]
	ord.Entry = ordMap["entry"]
	ord.DateCreated = ordMap["date_created"]

	return ord
}

func (mtsc *MapToStructConvertor) ConvertorDelivery(delRep *DeliveryRep) *Delivery {
	delMap := delRep.Get()
	del := new(Delivery)
	del.Name = delMap["name"]
	del.Phone = delMap["phone"]
	del.Email = delMap["email"]

	return del
}
