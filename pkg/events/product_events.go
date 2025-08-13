package events

// ProductCreated evento de produto criado
type ProductCreated struct {
	BaseEvent
	Product ProductData `json:"product"`
}

// ProductData dados do produto
type ProductData struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

// ProductUpdated evento de produto atualizado
type ProductUpdated struct {
	BaseEvent
	Product ProductData `json:"product"`
}
