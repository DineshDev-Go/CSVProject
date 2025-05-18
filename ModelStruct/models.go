package modelstruct

type ResponseStruct struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

type RevenueResponse struct {
	TotalRevenue float64 `json:"totRevenue"`
	ErrMsg       string  `json:"errMsg"`
	Status       string  `json:"status"`
}

type CustomerData struct {
	CustomerId      string `gorm:"column:customer_id;primaryKey"`
	CustomerName    string `gorm:"column:customer_name"`
	CustomerAddress string `gorm:"column:customer_address"`
	CustomerEmail   string `gorm:"column:customer_email"`
	CreatedDate     string `gorm:"column:created_date"`
	CreatedBy       string `gorm:"column:created_by"`
}

type ProductData struct {
	ProductId   string  `gorm:"column:product_id;primaryKey"`
	ProductName string  `gorm:"column:product_name"`
	UnitPrice   float64 `gorm:"column:unit_price"`
	Category    string  `gorm:"column:category"`
	CreatedDate string  `gorm:"column:created_date"`
	CreatedBy   string  `gorm:"column:created_by"`
}

type AllDataStruct struct {
	OrderId         string  `gorm:"column:order_id"`
	ProductId       string  `gorm:"column:product_id"`
	CustomerId      string  `gorm:"column:customer_id"`
	ProductName     string  `gorm:"column:product_name"`
	Discount        float64 `gorm:"column:discount"`
	ShippingCost    float64 `gorm:"column:shipping_cost"`
	PaymentMethod   string  `gorm:"column:payment_method"`
	CustomerName    string  `gorm:"column:customer_name"`
	CustomerEmail   string  `gorm:"column:customer_email"`
	Category        string  `gorm:"column:category"`
	Region          string  `gorm:"column:region"`
	DateOfSale      string  `gorm:"column:date_of_sale"`
	QuantitySold    int     `gorm:"column:quantity_sold"`
	UnitPrice       float64 `gorm:"column:unit_price"`
	CustomerAddress string  `gorm:"column:customer_address"`
	CreatedDate     string  `gorm:"column:created_date"`
	CreatedBy       string  `gorm:"column:created_by"`
	UniqueId        string  `gorm:"column:unique_id"`
}

type OrderData struct {
	OrderId       string  `gorm:"column:order_id;primaryKey"`
	QuantitySold  int     `gorm:"quantity_sold"`
	UnitPrice     float64 `gorm:"unit_price"`
	Discount      float64 `gorm:"discount"`
	Region        string  `json:"region"`
	ShippingCost  float64 `gorm:"shipping_cost"`
	PaymentMethod string  `gorm:"payment_method"`
	DateOfSale    string  `gorm:"date_of_sale"`
	ProductId     string  `gorm:"column:product_id"`
	CustomerId    string  `gorm:"column:customer_id"`
	CreatedDate   string  `gorm:"column:created_date"`
	CreatedBy     string  `gorm:"column:created_by"`
}
