package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nistagram/agent/dto"
	"nistagram/agent/service"
	"nistagram/util"
)

type ProductHandler struct {
	ProductService *service.ProductService
}

func (handler *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request){
	loggedUserID := util.GetLoggedUserIDFromToken(r)
	var productDto dto.ProductDTO
	err := json.NewDecoder(r.Body).Decode(&productDto)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.ProductService.CreateProduct(productDto, loggedUserID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request){
	products := handler.ProductService.GetAllProducts()
	js, err := json.Marshal(products)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}
