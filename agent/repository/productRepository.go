package repository

import (
	"fmt"
	"gorm.io/gorm"
	"nistagram/agent/model"
)

type ProductRepository struct {
	Database *gorm.DB
}

func (repo *ProductRepository) CreateProduct(product *model.Product) error{
	result := repo.Database.Create(product)
	if result.RowsAffected == 0 {
		return fmt.Errorf("Product not created")
	}
	fmt.Println("Product created")
	return nil
}

func (repo *ProductRepository) CreateAgentProduct(agentProduct *model.AgentProduct) error{
	result := repo.Database.Create(agentProduct)
	if result.RowsAffected == 0 {
		return fmt.Errorf("Agent product not created")
	}
	fmt.Println("Agent product created")
	return nil
}

func (repo *ProductRepository) GetAllProducts() []model.Product{
	var result []model.Product
	repo.Database.Find(&result)
	return result
}

func (repo *ProductRepository) GetProductsValidPrice(productId uint) float32{
	var result model.AgentProduct
	repo.Database.Table("agent_products").Find(&result, "is_valid = 1 and product_id = ?", productId)
	return result.PricePerItem
}
