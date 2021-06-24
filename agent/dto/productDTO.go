package dto

type ProductDTO struct {
	Name            string        `json:"name"`
	Quantity        uint          `json:"quantity"`
	PricePerItem    float32       `json:"pricePerItem"`
}
