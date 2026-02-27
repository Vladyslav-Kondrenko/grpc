package product

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type ProductCreateRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type ProductUpdateRequest struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}
