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
	"fmt"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	pb "github.com/GoogleCloudPlatform/microservices-demo/src/shippingservice/genproto"
)

type MyEvent struct {
    ShipOrder *pb.ShipOrderRequest
    GetQuote  *pb.GetQuoteRequest
}

func HandleLambdaEvent(event MyEvent) (string, error){
    shipResponse, shipErr := ShipOrder(event.ShipOrder)

    // Process the GetQuote request
    quoteResponse, quoteErr := GetQuote(event.GetQuote)

    if shipErr != nil || quoteErr != nil {
        return "", shipErr
    }

    // Combine the ShipOrder and GetQuote responses into a JSON structure
    response := struct {
        ShipOrderResponse *pb.ShipOrderResponse
        GetQuoteResponse  *pb.GetQuoteResponse
    }{
        ShipOrderResponse: shipResponse,
        GetQuoteResponse:  quoteResponse,
    }

    jsonResponse, err := json.Marshal(response)
    if err != nil {
        return "", err
    }

    return string(jsonResponse), nil
}
func main() {
	lambda.Start(HandleLambdaEvent)
}

func GetQuote(in *pb.GetQuoteRequest) (*pb.GetQuoteResponse, error) {
	print("[GetQuote] received request")
	print("[GetQuote] completed request")

	// 1. Generate a quote based on the total number of items to be shipped.
	quote := CreateQuoteFromCount(0)

	// 2. Generate a response.
	return &pb.GetQuoteResponse{
		CostUsd: &pb.Money{
			CurrencyCode: "USD",
			Units:        int64(quote.Dollars),
			Nanos:        int32(quote.Cents * 10000000)},
	}, nil

}

// ShipOrder mocks that the requested items will be shipped.
// It supplies a tracking ID for notional lookup of shipment delivery status.
func  ShipOrder(in *pb.ShipOrderRequest) (*pb.ShipOrderResponse, error) {
	print("[ShipOrder] received request")
	// 1. Create a Tracking ID
	baseAddress := fmt.Sprintf("%s, %s, %s", in.Address.StreetAddress, in.Address.City, in.Address.State)
	id := CreateTrackingId(baseAddress)

	// 2. Generate a response.
	return &pb.ShipOrderResponse{
		TrackingId: id,
	}, nil
}

