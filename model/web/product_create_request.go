package web

type ProductCreateRequest struct {
	ProductName string `validate:"required,max=255,min=1" json:"product_name"`
	Price       int    `validate:"required" json:"price"`
}
