package logic

type orderPayload struct {
	OrderId   string `json:"order_id"`
	ProductId string `json:"product_id"`
}
