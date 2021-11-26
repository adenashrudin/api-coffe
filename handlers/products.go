// Package classification of Product API
//
// Documentation for Product API
//
// Schemes: http
// Basepath: /
// Version: 1.0.0
//
// Consumes:
// - aplication/json
//
// Produces:
// - aplication/json
// swagger:meta
package handlers

import (
	"api_coffe/data"
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

// A list of products returns in the responses
// swagger:response productsResponses
type productsResponseWrapper struct{
	// All products in the system
	// in: body
	Body []data.Product
}

// swagger:parameters deleteProduct
type productIDParameterWrapper struct {
	// The id of the product to delete from the database
	// in:path
	// required:true
	ID int `json:"id"`
}

// swagger: response noContent
type productNoContent struct {

}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//	200: productsResponses

// GetProducts returns the product from the data base
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unnable marshal JSON", http.StatusInternalServerError)
	}
}

// swagger:route POST /products products product
// Return a status 
// responses :
// 	200: productsResponses

// AddProduct store data product to the data base
func (p *Products) AddProducts(rw http.ResponseWriter, r *http.Request) {

	p.l.Println("ADD/POST PRODUCT")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
	rw.Write([]byte("Success"))

}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unnable get ID", http.StatusBadRequest)
		return
	}

	p.l.Println("Get ID", id)

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found ", http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

// swagger:route DELETE /products/{id} products deleteProduct
// Returns a list of products
// responses: 
//	200: noContent

// DeleteProduct delete a product from the database 
func(p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request){
	// this will always convert because of the router
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p.l.Println("Handle DELETE Product-->", id)

	err := data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not Found ", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Something wrong ", http.StatusInternalServerError)
		return
	}
	
	rw.Write([]byte("Success"))

}

// func(p *Products) GetProductByID (rw http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["id"])

// 	if err!=nil {
// 		http.Error(rw, "Unable Get ID", http.StatusBadRequest)
// 		return
// 	}

// 	// pr,err := data.GetProductById(id)

// 	// err := pr.TO

// 	p.l.Println("GET PRODUCT BY ID", id)
// }

// func(p *Products) deleteProducts(rw http.ResponseWriter, r *http.Request){

// 	p.l.Println("DELETE PRODUCT")
// }

type KeyProduct struct{}

//MiddelwareVlaidateProduct validates the product in the Request and call the nex if ok!
func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		prod := data.Product{}
		err := prod.FromJSON(r.Body)

		if err != nil {
			p.l.Println("[ERROR] deseriliazting product", err)
			http.Error(rw, "Unnable unmarshal JSON", http.StatusBadRequest)
			return
		}

		//Add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		//Call the next handler, which can be another middleware in the chain, or the final handler
		next.ServeHTTP(rw, r)
	})
}
