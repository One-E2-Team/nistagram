package dto

type UpdateProductDTO struct {
	ProductId		uint		  `json:"productId"`
	Name            string        `json:"name"`
	Quantity        uint          `json:"quantity"`
	PricePerItem    float32       `json:"pricePerItem"`
}
