// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/productcatalogservice/genproto"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)
var (
	cat          pb.ListProductsResponse
	reloadCatalog bool
)


type MyEvent struct {
	Type  string `json:"type" default:"invalid"`
	Query string `json:"query" default:"" `
	ID    string `json:"id" default:""`
}


func HandleLambdaEvent(request MyEvent)(events.APIGatewayProxyResponse, error) {
	inputData := request.Type
	//fmt.Println(inputData)
	var response interface{}
	var err error
	switch inputData {
	case "list":
		emptyInstance := &pb.Empty{}
		products := ListProducts(emptyInstance)
		response = products
	case "search":
		query := request.Query
		searchRequest := &pb.SearchProductsRequest{Query: query}
		searchResponse := SearchProducts(searchRequest)
		response = searchResponse
	case "get":
		id:=request.ID
		getRequest := &pb.GetProductRequest{Id: id}
		getResponse := GetProduct(getRequest)
		response = getResponse
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid type parameter in JSON request",
		}, nil
	}


	responseJSON, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseJSON),
	}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}

func readCatalogFile(catalog *pb.ListProductsResponse) error {
	catalogJSON, err := ioutil.ReadFile("products.json")
	if err!=nil{
		return err
	}

    if err := json.Unmarshal([]byte(catalogJSON), catalog); err != nil {
        //fmt.Println("Failed to parse the catalog JSON:", err)
        return err
    }
    return nil
}
type productCatalog struct{}


func parseCatalog() []*pb.Product {
	reloadCatalog = true
	cat = pb.ListProductsResponse{}
	if reloadCatalog || len(cat.Products) == 0 {
		err := readCatalogFile(&cat)
		if err != nil {
			return []*pb.Product{}
		}
	}
	return cat.Products
}



func ListProducts(*pb.Empty) *pb.ListProductsResponse {
    // for _, product := range parseCatalog() {
    //     fmt.Printf("Product ID: %s\n", product.Id)
    //     fmt.Printf("Name: %s\n", product.Name)
    //     fmt.Printf("Description: %s\n", product.Description)
    //     fmt.Printf("PriceUSD: %+v\n", product.PriceUsd) // Print PriceUSD
    //     // Print other product fields as needed
    //     fmt.Println("------------------------")
    // }
    return &pb.ListProductsResponse{Products: parseCatalog()}
}

func  GetProduct(req *pb.GetProductRequest) (*pb.Product) {
	var found *pb.Product
	found_value := 0
	//print(len(parseCatalog()))
	for i := 0; i < len(parseCatalog()); i++ {
		//(parseCatalog()[i].Id)
		if req.Id == parseCatalog()[i].Id {
			found_value = 1
			found = parseCatalog()[i]
			//fmt.Printf("PriceUSD: %+v\n", found.PriceUsd)
		}
	}
	if found_value ==1 {
		return found
	}else {
		return nil
	}
	
	
}

func SearchProducts(req *pb.SearchProductsRequest) (*pb.SearchProductsResponse) {
	//print(req.Query)
	var ps []*pb.Product
	for _, p := range parseCatalog() {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(req.Query)) ||
			strings.Contains(strings.ToLower(p.Description), strings.ToLower(req.Query)) {
			ps = append(ps, p)
		}
	}
	return &pb.SearchProductsResponse{Results: ps}
}
