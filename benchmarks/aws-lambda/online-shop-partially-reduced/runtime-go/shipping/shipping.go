package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)
type CartItem struct {
	ProductId string
	Quantity  int32
}

type Address struct {
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       int32
}


type ShipOrderRequest struct {
	Address              *Address    `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Items                []*CartItem `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
}

type GetQuoteRequest struct {
	Address              *Address    `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Items                []*CartItem `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
}
type MyEvent struct {
    ShipOrder *ShipOrderRequest
    GetQuote  *GetQuoteRequest
}
type ShipOrderResponse struct {
	TrackingId           string  
}

type Money struct {
	CurrencyCode string `json:"currency_code"`
	Units        int64  `json:"units"`
	Nanos        int32  `json:"nanos"`
}

type GetQuoteResponse struct {
	CostUsd              *Money   `protobuf:"bytes,1,opt,name=cost_usd,json=costUsd,proto3" json:"cost_usd,omitempty"`
}

func main() {
	arg1 := os.Args[1] 
	var shipRequest ShipOrderRequest
	if err := json.Unmarshal([]byte(arg1), &shipRequest); err != nil {
		fmt.Println("Error parsing the argument:", err)
		return
	}
	var quoteRequest GetQuoteRequest
	if err := json.Unmarshal([]byte(arg1), &quoteRequest); err != nil {
		fmt.Println("Error parsing the argument:", err)
		return
	}
    shipResponse, _ := ShipOrder(&shipRequest)
    quoteResponse, _ := GetQuote(&quoteRequest)
    response := struct {
        ShipOrderResponse *ShipOrderResponse
        GetQuoteResponse  *GetQuoteResponse
    }{
        ShipOrderResponse: shipResponse,
        GetQuoteResponse:  quoteResponse,
    }
    jsonResponse, _:= json.Marshal(response)
    fmt.Print(string(jsonResponse))
}

type Quote struct {
	Dollars uint32
	Cents   uint32
}
func CreateQuoteFromFloat(value float64) Quote {
	units, fraction := math.Modf(value)
	return Quote{
		uint32(units),
		uint32(math.Trunc(fraction * 100)),
	}
}
func CreateQuoteFromCount(count int) Quote {
	return CreateQuoteFromFloat(8.99)
}
func GetQuote(in *GetQuoteRequest) (*GetQuoteResponse, error) {
	quote := CreateQuoteFromCount(0)
	return &GetQuoteResponse{
		CostUsd: &Money{
			CurrencyCode: "USD",
			Units:        int64(quote.Dollars),
			Nanos:        int32(quote.Cents * 10000000)},
	}, nil

}
var seeded bool = false
// getRandomLetterCode generates a code point value for a capital letter.
func getRandomLetterCode() uint32 {
	return 65 + uint32(rand.Intn(25))
}

// getRandomNumber generates a string representation of a number with the requested number of digits.
func getRandomNumber(digits int) string {
	str := ""
	for i := 0; i < digits; i++ {
		str = fmt.Sprintf("%s%d", str, rand.Intn(10))
	}

	return str
}
func CreateTrackingId(salt string) string {
	if !seeded {
		rand.Seed(time.Now().UnixNano())
		seeded = true
	}

	return fmt.Sprintf("%c%c-%d%s-%d%s",
		getRandomLetterCode(),
		getRandomLetterCode(),
		len(salt),
		getRandomNumber(3),
		len(salt)/2,
		getRandomNumber(7),
	)
}

// ShipOrder mocks that the requested items will be shipped.
// It supplies a tracking ID for notional lookup of shipment delivery status.
func  ShipOrder(in *ShipOrderRequest) (*ShipOrderResponse, error) {
	// 1. Create a Tracking ID
	baseAddress := fmt.Sprintf("%s, %s, %s", in.Address.StreetAddress, in.Address.City, in.Address.State)
	id := CreateTrackingId(baseAddress)

	// 2. Generate a response.
	return &ShipOrderResponse{
		TrackingId: id,
	}, nil
}

