package handlers

import (
	"context"
	"hello/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)



type Products struct{
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}



func(p*Products) GetProducts(rw http.ResponseWriter, r *http.Request){
	lp := data.GetProducts()
	err := lp.ToJSON(rw)

	if err !=nil {
		http.Error(rw, "Unnable marshal JSON", http.StatusInternalServerError)
	}
}

func(p *Products) AddProducts (rw http.ResponseWriter,  r *http.Request){
	
	p.l.Println("ADD/POST PRODUCT")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
	rw.Write([]byte("Success"))

}

func(p *Products) UpdateProduct (rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		 http.Error(rw, "Unnable get ID", http.StatusBadRequest)
		 return
	 }

	p.l.Println("Get ID",id)

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

type KeyProduct struct {}

//MiddelwareVlaidateProduct validates the product in the Request and call the nex if ok!
func(p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		
		prod := data.Product{}
		err := prod.FromJSON(r.Body) 

		if err!=nil {
			p.l.Println("[ERROR] deseriliazting product",err)
			http.Error(rw, "Unnable unmarshal JSON", http.StatusBadRequest)
			return
		}

		 //Add the product to the context
		 ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		 r  = r.WithContext(ctx)

		 //Call the next handler, which can be another middleware in the chain, or the final handler
		 next.ServeHTTP(rw,r)
	})
}