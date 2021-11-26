package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type Product struct {
	ID			int		`json:"id"`
	Name		string	`json:"name" validate:"required"`
	Description string	`json:"description"`
	Price		float32	`json:"price" validate:"gt=0"`
	SKU			string	`json:"sku" validate:"required,sku"`
	CreatedOn	string	`json:"-"`
	UpdatedOn	string	`json:"-"`
	DeletedOn	string	`json:"-"`
}

func( p *Product) Valdate()  error {
	validate := validator.New()

	validate.RegisterValidation("sku",validateSKU)

	return validate.Struct(p)

}

func validateSKU(fl validator.FieldLevel) bool {
	rg := regexp.MustCompile(`cf+-[0-9]+`)
	
	matches := rg.FindAllString(fl.Field().String(),-1)
	if len(matches) != 1 {
		return false
	}
	
	return true 
}

type Products []*Product


func(p *Product) FromJSON(r io.Reader) error {
	e:= json.NewDecoder(r)
	return e.Decode(p)
}


func(p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetProducts() Products{
	return productList
}

// func GetProductById(id int) (Products, error) {

// 	_, _,err :=FindProduct(id)
// 	if err!= nil {
// 		return nil,err
// 	}
// 	return productList[0],nil
// }


func AddProduct(p *Product) {
	p.ID = GetNextID()
	productList =append(productList, p)
}

func UpdateProduct(id int , p *Product) error {
	_ , i, err := FindProduct(id)
	if err !=nil {
		return err
	}

	p.ID = id
	productList[i] = p

	return nil
	
}

var ErrProductNotFound = fmt.Errorf("Product Not Found!")

func FindProduct(id int) (*Product, int, error) {
	for ix, p:= range productList {
		
		if p.ID== id {
		return p, ix, nil
		}
	}

	return nil, -1, ErrProductNotFound
}

func GetNextID () int {
	lp := productList[len(productList)-1]
	return lp.ID +1
}



var productList = []*Product{
	&Product{ 
		ID:				1,
		Name:			"Latte",
		Description: 	"Frothy milky coffe",
		Price: 			2.45,
		SKU: 			"cf-001",
		CreatedOn: 		time.Now().UTC().String(),	
		UpdatedOn: 		time.Now().UTC().String(),	
	},

	&Product{ 
		ID:				2,
		Name:			"Espresso",
		Description: 	"Short and strong coffee without milk",
		Price: 			1.99,
		SKU: 			"cf-002",
		CreatedOn: 		time.Now().UTC().String(),	
		UpdatedOn: 		time.Now().UTC().String(),	
	},
}