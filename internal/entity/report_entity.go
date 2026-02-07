package entity

type Report struct {
	TotalRevenue     int           `gorm:"column:total_revenue"`
	TotalTransaction int           `gorm:"column:total_transaction"`
	BestProduct      []BestProduct `gorm:"-"`
}

type BestProduct struct {
	Name     string `gorm:"column:name"`
	Quantity int    `gorm:"column:quantity"`
	Subtotal int    `gorm:"column:subtotal"`
}
