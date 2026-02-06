package dto

type Checkout struct {
	Items []CheckoutItem `json:"items"`
}

type CheckoutItem struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}
