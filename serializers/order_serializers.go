package serializers

import (
	"ginserver/models"
)

// OrdersSubsetJSON ..
type OrdersSubsetJSON struct {
	ListMeta
	Values []models.Order `json:"data"`
}

// NewOrdersSubsetJSON ...
func NewOrdersSubsetJSON(orders []models.Order, count int, skip int, total int) OrdersSubsetJSON {
	json := OrdersSubsetJSON{
		Values: []models.Order{},
		ListMeta: ListMeta{
			Count: count,
			Skip:  skip,
			Total: total,
		},
	}
	for _, order := range orders {
		json.Values = append(json.Values, order)
	}
	return json
}

// SerializeOrders ...
func SerializeOrders(orders []models.Order, params ...interface{}) interface{} {
	count := params[0].(int64)
	skip := params[1].(int64)
	total := params[2].(int)
	orderSubsetJSON := NewOrdersSubsetJSON(orders, int(count), int(skip), total)
	return NewResponse(0, orderSubsetJSON, "Success")
}

// SerializeOrder ..
func SerializeOrder(order models.Order) interface{} {
	return NewResponse(0, order, "Success")
}
