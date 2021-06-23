package dto

type ProductDTO struct {
	Name            string        `json:"name"`
	PicturePath     string        `json:"picturePath"`
	Quantity        uint          `json:"quantity"`
	PricePerItem    float32       `json:"pricePerItem"`
}
