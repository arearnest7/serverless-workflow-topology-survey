package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	cat           ListProductsResponse
	reloadCatalog bool
)

type GetProductRequest struct {
	Id string
}

func main() {
	myEvent := GetProductRequest{
		Id: os.Args[1],
	}
	result := GetProduct(&myEvent)
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Print(string(jsonBytes))
}

func readCatalogFile(catalog *ListProductsResponse) error {
	catalogJSON, err := ioutil.ReadFile("products.json")
	if err != nil {
		print("Error in reading file")
	}

	if err := json.Unmarshal([]byte(catalogJSON), catalog); err != nil {
		return err
	}
	return nil
}

type productCatalog struct{}

func parseCatalog() []*Product {
	reloadCatalog = true
	cat = ListProductsResponse{}
	if reloadCatalog || len(cat.Products) == 0 {
		err := readCatalogFile(&cat)
		if err != nil {
			return []*Product{}
		}
	}
	return cat.Products
}

type Empty struct {
}

type ListProductsResponse struct {
	Products []*Product
}

type Money struct {
	CurrencyCode string `protobuf:"bytes,1,opt,name=currency_code,json=currencyCode,proto3" json:"currency_code,omitempty"`
	Units        int64  `protobuf:"varint,2,opt,name=units,proto3" json:"units,omitempty"`
	Nanos        int32  `protobuf:"varint,3,opt,name=nanos,proto3" json:"nanos,omitempty"`
}
type Product struct {
	Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Picture     string `protobuf:"bytes,4,opt,name=picture,proto3" json:"picture,omitempty"`
	PriceUsd    *Money `protobuf:"bytes,5,opt,name=price_usd,json=priceUsd,proto3" json:"price_usd,omitempty"`
}

func GetProduct(req *GetProductRequest) *Product {
	var found *Product
	found_value := 0
	for i := 0; i < len(parseCatalog()); i++ {
		if req.Id == parseCatalog()[i].Id {
			found_value = 1
			found = parseCatalog()[i]
		}
	}
	if found_value == 1 {
		return found
	} else {
		return nil
	}
}
