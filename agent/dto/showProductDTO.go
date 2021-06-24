package dto

type ShowProductDTO struct {
	ID              uint          `json:"id"`
	Name            string        `json:"name"`
	PicturePath     string        `json:"picturePath"`
	PricePerItem    float64       `json:"pricePerItem"`
}
