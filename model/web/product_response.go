package web

type ProductResponse struct {
	Id          int    `json:"id"`
	ProductName string `json:"product_name"`
	Price       int    `json:"price"`
}
