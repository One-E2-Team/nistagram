package service

import (
	"errors"
	"nistagram/agent/dto"
	"nistagram/agent/model"
	"nistagram/agent/repository"
	"nistagram/agent/util"
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
		agentProduct, err := service.ProductRepository.GetValidAgentProductByProductId(p.ID)
		if err != nil{
			continue
		}
		retItem := dto.ShowProductDTO{ID: p.ID, Name: p.Name, PicturePath: p.PicturePath,
			PricePerItem: agentProduct.PricePerItem, Quantity: agentProduct.Quantity}
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

func (service *ProductService) CreateOrder(dto dto.OrderDTO, loggedUserId uint) error{
	err := service.checkQuantities(dto)
	if err != nil {
		return err
	}

	order, err := service.makeOrder(dto, loggedUserId)
	if err != nil {
		return err
	}

	err = service.ProductRepository.CreateOrder(order)
	return err
}

func (service *ProductService) GetProductById(productId uint) (*model.Product, error){
	return service.ProductRepository.GetProductById(productId)
}

func (service *ProductService) makeOrder(dto dto.OrderDTO, loggedUserId uint) (*model.Order, error) {
	var items []model.Item
	var agentId uint
	fullPrice := 0.0
	for _, i := range dto.Items {
		agentProduct, err := service.ProductRepository.GetValidAgentProductByProductId(i.ProductId)
		if err != nil {
			return nil, err
		}
		agentProduct.Quantity -= i.Quantity
		err = service.ProductRepository.UpdateAgentProduct(agentProduct)
		agentId = agentProduct.UserID
		item := &model.Item{ProductId: i.ProductId, Quantity: i.Quantity}
		items = append(items, *item)
		fullPrice += agentProduct.PricePerItem * float64(i.Quantity)
	}

	order := &model.Order{UserId: loggedUserId, Timestamp: time.Now(),
		Items: items, FullPrice: fullPrice, AgentId: agentId}
	return order, nil
}

func (service *ProductService) checkQuantities(dto dto.OrderDTO) error {
	for _, i := range dto.Items{
		agentProduct, err := service.ProductRepository.GetValidAgentProductByProductId(i.ProductId)
		if err != nil{
			return err
		}
		if agentProduct.Quantity < i.Quantity{
			return errors.New("Product id " + util.Uint2String(i.ProductId) + " does not have enough quantity!")
		}
	}
	return nil
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

