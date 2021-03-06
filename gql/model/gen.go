// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type NewProduct struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

type Product struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	Seller string `json:"seller"`
}

type UpdateProduct struct {
	ID    string  `json:"id"`
	Name  *string `json:"name"`
	Price *int64  `json:"price"`
}
