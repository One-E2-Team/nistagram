package dto

type OrderDTO struct {
	Items        []ItemDTO        `json:"items"`
}

type ItemDTO struct{
	ProductId       uint       `json:"productId"`
	Quantity        uint       `json:"quantity"`
}
