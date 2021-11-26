package handlers

import (
	"hello/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)



type Products struct{
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}


func(p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		p.getProducts(rw,r)
		return
	}else if r.Method == http.MethodPut {
		
		reg := regexp.MustCompile(`/([0-9]+)`)
		g :=reg.FindAllStringSubmatch(r.URL.Path,-1)
		if len(g) != 1 {
			http.Error(rw, "Invalid URL", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			http.Error(rw, "Invalid URI-2", http.StatusBadRequest)
			return
		}

		idString := g[0][1]

		id,err := strconv.Atoi(idString)

		if err !=nil {
			http.Error(rw, "Invalid URI-3", http.StatusBadRequest)
			return
		}

		p.l.Println("Got ID PATH-->",id)

		p.putProducts(id,rw,r)
		return
	}else if r.Method == http.MethodPost {
		p.addProducts(rw,r)
		return
	}else if r.Method == http.MethodDelete {
		p.deleteProducts(rw,r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func(p*Products) getProducts(rw http.ResponseWriter, r *http.Request){
	lp := data.GetProducts()
	err := lp.ToJSON(rw)

	if err !=nil {
		http.Error(rw, "Unnable marshal JSON", http.StatusInternalServerError)
	}
}

func(p *Products) addProducts (rw http.ResponseWriter,  r *http.Request){
	p.l.Println("ADD/POST PRODUCT")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body) 

	if err!=nil {
		http.Error(rw, "Unnable unmarshal JSON", http.StatusInternalServerError)
	}

	p.l.Printf("Product: %#v", prod)

	data.AddProduct(prod)

	rw.Write([]byte("Success"))

}

func(p *Products) putProducts (id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("PUT PRODUCT")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body) 

	if err!=nil {
		http.Error(rw, "Unnable unmarshal JSON", http.StatusBadRequest)
	}

	p.l.Printf("Product: %#v", prod)

	err = data.UpdateProduct(id, prod)
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

func(p *Products) deleteProducts(rw http.ResponseWriter, r *http.Request){
	p.l.Println("DELETE PRODUCT")
}