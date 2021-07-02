package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"mime/multipart"
	"net/http"
	"nistagram/agent/dto"
	"nistagram/agent/service"
	"nistagram/agent/util"
	"os"
	"strings"
)

type ProductHandler struct {
	ProductService *service.ProductService
}

func (handler *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request){
	loggedUserID := util.GetLoggedUserIDFromToken(r)
	var productDto dto.ProductDTO

	if err := r.ParseMultipartForm(0); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := r.MultipartForm.Value["data"]

	if err := json.Unmarshal([]byte(data[0]), &productDto); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	picture, picHeader, err := r.FormFile("file")
	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uid := uuid.NewString()
	fileSplitted := strings.Split(picHeader.Filename, ".")
	fileName := uid + "." + fileSplitted[1]

	err = handler.ProductService.CreateProduct(productDto, loggedUserID, fileName)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = handler.savePicture(fileName, picture)
	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *ProductHandler) savePicture(fileName string, picture multipart.File) error {
	f, err := os.OpenFile("../../agentstaticdata/data/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	defer func(picture multipart.File) {
		_ = picture.Close()
	}(picture)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, picture)
	if err != nil {
		return err
	}
	return nil
}

func (handler *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request){
	products := handler.ProductService.GetAllProducts()
	js, err := json.Marshal(products)
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	productId := util.String2Uint(vars["id"])

	err := handler.ProductService.DeleteProduct(productId)
	if err != nil{
		fmt.Println(err)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request){
	var updateProductDto dto.UpdateProductDTO
	err := json.NewDecoder(r.Body).Decode(&updateProductDto)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.ProductService.UpdateProduct(updateProductDto)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *ProductHandler) CreateOrder(w http.ResponseWriter, r *http.Request){
	loggedUserId := util.GetLoggedUserIDFromToken(r)
	var orderDto dto.OrderDTO
	err := json.NewDecoder(r.Body).Decode(&orderDto)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.ProductService.CreateOrder(orderDto, loggedUserId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	productId := util.String2Uint(vars["id"])

	result, err := handler.ProductService.GetProductById(productId)
	if err != nil{
		fmt.Println(err)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	w.Header().Set("Content-Type", "application/json")
}
