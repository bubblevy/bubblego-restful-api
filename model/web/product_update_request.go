package web

type ProductUpdateRequest struct {
	Id          int    `validate:"required"`
	ProductName string `validate:"required,max=255,min=1" json:"product_name"`
	Price       int    `validate:"required" json:"price"`
}
