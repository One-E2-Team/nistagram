package dto

type ShowProductDTO struct {
	Name            string        `json:"name"`
	PicturePath     string        `json:"picturePath"`
	PricePerItem    float32       `json:"pricePerItem"`
}
