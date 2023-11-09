


package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	cat           ListProductsResponse
	reloadCatalog bool
)


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

func ListProducts() {
    for _, product := range parseCatalog() {
        fmt.Printf("Product ID: %s\n", product.Id)
        fmt.Printf("Name: %s\n", product.Name)
        fmt.Printf("Description: %s\n", product.Description)
        fmt.Printf("PriceUSD: %+v\n", product.PriceUsd) // Print PriceUSD
        // Print other product fields as needed
        fmt.Println("------------------------")
    }
    //return &ListProductsResponse{Products: parseCatalog()}
}



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

func main() {
	ListProducts()
}

func readCatalogFile(catalog *ListProductsResponse) error {
	catalogJSON, err := ioutil.ReadFile("./products.json")
	if err != nil {
		print("Error in reading file")
	}

	if err := json.Unmarshal([]byte(catalogJSON), catalog); err != nil {
		return err
	}
	return nil
}



type SearchProductsRequest struct {
	Query                string   
}

type SearchProductsResponse struct {
	Results              []*Product 

}
func SearchProducts(req *SearchProductsRequest) (*SearchProductsResponse) {
	print(req.Query)
	var ps []*Product
	for _, p := range parseCatalog() {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(req.Query)) ||
			strings.Contains(strings.ToLower(p.Description), strings.ToLower(req.Query)) {
			ps = append(ps, p)
		}
	}
	return &SearchProductsResponse{Results: ps}
}
