package data

import "testing"


func TestCheckValidation(t *testing.T){
	p := &Product{
		Name: "Nic",
		Price: 1.00,
		SKU: "cf-003",
	}

	err := p.Valdate()

	if err!=nil {
		t.Fatal(err)
	}
}