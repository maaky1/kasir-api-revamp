package dto

type Report struct {
	ReportRange      string        `json:"repor_range"`
	TotalRevenue     int           `json:"total_revenue"`
	TotalTransaction int           `json:"total_transaction"`
	BestProduct      []BestProduct `json:"best_product"`
}

type BestProduct struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Subtotal int    `json:"subtotal"`
}
