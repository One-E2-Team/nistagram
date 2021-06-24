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

func (service *ProductService) CreateProduct(dto dto.ProductDTO, loggedUserId uint, fileName string) error{

	product := model.Product{Name: dto.Name, PicturePath: fileName}

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

func (service *ProductService) GetAllProducts() []dto.ShowProductDTO{
	var ret []dto.ShowProductDTO
	products := service.ProductRepository.GetAllProducts()

	for _, p := range products{
		retItem := dto.ShowProductDTO{Name: p.Name, PicturePath: p.PicturePath,
			PricePerItem: service.ProductRepository.GetProductsValidPrice(p.ID)}
		ret = append(ret, retItem)
	}

	return ret
}

func (service *ProductService) DeleteProduct(productId uint) error{
	return service.ProductRepository.DeleteProduct(productId)
}

func (service *ProductService) UpdateProduct(dto dto.UpdateProductDTO) error{
	err := service.updateProduct(dto)
	if err != nil {
		return err
	}

	agentProduct, err := service.ProductRepository.GetValidAgentProductByProductId(dto.ProductId)
	if err != nil {
		return err
	}

	if agentProduct.PricePerItem != dto.PricePerItem {
		err = service.changePrice(dto, agentProduct)
	} else {
		agentProduct.Quantity = dto.Quantity
		err = service.ProductRepository.UpdateAgentProduct(agentProduct)
	}

	return err
}

func (service *ProductService) changePrice(dto dto.UpdateProductDTO, agentProduct *model.AgentProduct) error {
	agentProduct.IsValid = false
	err := service.ProductRepository.UpdateAgentProduct(agentProduct)
	if err != nil {
		return err
	}
	newAgentProduct := &model.AgentProduct{UserID: agentProduct.UserID, ProductID: agentProduct.ProductID,
		Quantity: dto.Quantity, PricePerItem: dto.PricePerItem, ValidFrom: time.Now(), IsValid: true}
	err = service.ProductRepository.CreateAgentProduct(newAgentProduct)
	return err
}

func (service *ProductService) updateProduct(dto dto.UpdateProductDTO) error {
	product, err := service.ProductRepository.GetProductById(dto.ProductId)
	if err != nil {
		return err
	}
	product.Name = dto.Name
	err = service.ProductRepository.UpdateProduct(product)
	return err
}

