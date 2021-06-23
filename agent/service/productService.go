package service

import (
	"nistagram/agent/dto"
	"nistagram/agent/model"
	"nistagram/agent/repository"
	"time"
)

type ProductService struct {
	ProductRepository *repository.ProductRepository
}

func (service *ProductService) CreateProduct(dto dto.ProductDTO, loggedUserId uint) error{

	product := model.Product{Name: dto.Name, PicturePath: dto.PicturePath}

	err := service.ProductRepository.CreateProduct(&product)
	if err != nil{
		return err
	}

	agentProduct := model.AgentProduct{UserID: loggedUserId, ProductID: product.ID, Quantity: dto.Quantity,
		PricePerItem: dto.PricePerItem, ValidFrom: time.Now(), IsValid: true}

	err = service.ProductRepository.CreateAgentProduct(&agentProduct)
	if err != nil{
		return err
	}

	return nil
}
