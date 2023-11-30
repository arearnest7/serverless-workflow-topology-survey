package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	pb "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/genproto"
	"github.com/aws/aws-lambda-go/events"
	money "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/money"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
)


const (
	usdCurrency = "USD"
)


type MyEvent struct {
    UserId string `json:"userId"`
    UserCurrency string `json:"userCurrency"`
    Address Address `json:"address"`
    Email string `json:"email"`
    CreditCard CreditCardInfo `json:"creditCard"`
}

type Address struct {
    StreetAddress string `json:"streetAddress"`
    City string `json:"city"`
    State string `json:"state"`
    Country string `json:"country"`
    ZipCode int32 `json:"zipCode"`
}

type CreditCardInfo struct {
    CreditCardNumber string `json:"creditCardNumber"`
    CreditCardCvv int32 `json:"creditCardCvv"`
    CreditCardExpirationYear int32 `json:"creditCardExpirationYear"`
    CreditCardExpirationMonth int32 `json:"creditCardExpirationMonth"`
}


func HandleLambdaEvent(event events.APIGatewayProxyRequest) (string, error) {
	var myEvent pb.PlaceOrderRequest
	if err := json.Unmarshal([]byte(event.Body), &myEvent); err != nil {
		return "",err
	}
	result,_:= PlaceOrder(&myEvent)
    return result,nil
} 

func main() {
    lambda.Start(HandleLambdaEvent)
}

func PlaceOrder(req *pb.PlaceOrderRequest) (string, error) {
	orderID, err := uuid.NewUUID()
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to generate order uuid")
	}
	prep, err := prepareOrderItemsAndShippingQuoteFromCart(req.UserId, req.UserCurrency, req.Address)
	if err != nil {
		return "", status.Errorf(codes.Internal, err.Error())
	}
	total := pb.Money{CurrencyCode: req.UserCurrency,
		Units: 0,
		Nanos: 0}
	total = money.Must(money.Sum(total, *prep.shippingCostLocalized))
	for _, it := range prep.orderItems {
		multPrice := money.MultiplySlow(*it.Cost, uint32(it.GetItem().GetQuantity()))
		total = money.Must(money.Sum(total, multPrice))
	}
	txID, err := chargeCard(&total, req.CreditCard)
	print(txID)
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to charge card: %+v", err)
	}
	shippingTrackingID, err := shipOrder(req.Address, prep.cartItems)
	if err != nil {
		return "", status.Errorf(codes.Unavailable, "shipping error: %+v", err)
	}


	orderResult := &pb.OrderResult{
		OrderId:            orderID.String(),
		ShippingTrackingId: shippingTrackingID,
		ShippingCost:       prep.shippingCostLocalized,
		ShippingAddress:    req.Address,
		Items:              prep.orderItems,
	}

	resp := &pb.PlaceOrderResponse{Order: orderResult}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	fmt.Println(string(jsonResp))
	fmt.Println(sendOrderConfirmation(req.Email, resp))
	return txID+string(string(jsonResp))+sendOrderConfirmation(req.Email, resp), nil
}

type orderPrep struct {
	orderItems            []*pb.OrderItem
	cartItems             []*pb.CartItem
	shippingCostLocalized *pb.Money
}

func prepareOrderItemsAndShippingQuoteFromCart(userID, userCurrency string, address *pb.Address) (orderPrep, error) {
	var out orderPrep
	cartItems, err := getUserCart(userID)
	if err != nil {
		return out, fmt.Errorf("cart failure: %+v", err)
	}
	orderItems, err := prepOrderItems(cartItems, userCurrency)
	if err != nil {
		return out, fmt.Errorf("failed to prepare order: %+v", err)
	}
	//fmt.Println("Qutote Shipping Function...")
	shippingUSD, err := quoteShipping(address, cartItems)
	if err != nil {
		return out, fmt.Errorf("shipping quote failure: %+v", err)
	}
	//fmt.Print("ShippingUSD:")
	//fmt.Println(shippingUSD)

	shippingPrice, err := convertCurrency(shippingUSD, userCurrency)
	if err != nil {
		return out, fmt.Errorf("failed to convert shipping cost to currency: %+v", err)
	}

	out.shippingCostLocalized = shippingPrice
	out.cartItems = cartItems
	out.orderItems = orderItems
	return out, nil
}

type QuoteResponse struct {
    ShipOrderResponse struct {
        TrackingID string `json:"tracking_id"`
    } `json:"ShipOrderResponse"`
    GetQuoteResponse struct {
        CostUSD struct {
            CurrencyCode string `json:"currency_code"`
            Units        int64    `json:"units"`
            Nanos        int32    `json:"nanos"`
        } `json:"cost_usd"`
    } `json:"GetQuoteResponse"`
}

func  quoteShipping(address *pb.Address, items []*pb.CartItem) (*pb.Money, error) {
	apiURL := "https://tekw7om46d.execute-api.us-east-2.amazonaws.com/default/shipping"
	addressMap := map[string]interface{}{
		"city": address.City,
		"country": address.Country,
		"state": address.State,
		"zipCode": address.ZipCode,
		"streetAddress": address.StreetAddress,     
	}

	itemsMap := make([]map[string]interface{}, len(items))
	for i, item := range items {
		itemMap := map[string]interface{}{
			"product_id": item.ProductId,
			"quantity": item.Quantity,
		}
		itemsMap[i] = itemMap
	}

	data := map[string]interface{}{
		"ShipOrder": map[string]interface{}{
			"Address": addressMap,
			"Items":   itemsMap,
		},
		"GetQuote": map[string]interface{}{
			"Address": addressMap,
			"Items":   itemsMap,
		},
	}
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", apiURL,bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var jsonStr string
	err = json.NewDecoder(resp.Body).Decode(&jsonStr)
	if err != nil {
		return nil, err
	}
	var response QuoteResponse
	
    err = json.Unmarshal([]byte(jsonStr), &response)
    if err != nil {
        //fmt.Println("Error:", err)
        return nil,err
    }
    getQuoteResponse := response.GetQuoteResponse
    //fmt.Printf("Currency Code: %s, Units: %d, Nanos: %d\n", getQuoteResponse.CostUSD.CurrencyCode, getQuoteResponse. CostUSD.Units, getQuoteResponse.CostUSD.Nanos)
	resultMoney := &pb.Money{
		CurrencyCode: getQuoteResponse.CostUSD.CurrencyCode,
		Units:        getQuoteResponse.CostUSD.Units,
		Nanos:        getQuoteResponse.CostUSD.Nanos,
	}
	
	return resultMoney, nil
}

func getUserCart(userID string) ([]*pb.CartItem, error) {
	apiURL := "https://c60ekgpu4l.execute-api.us-east-2.amazonaws.com/default/cartservice"
	data := map[string]interface{}{
		"userID": userID,
	}
	jsonData, _ := json.Marshal(data)

	// Create a POST request with JSON data
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println("Error making the request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// Unmarshal the JSON response into a slice of strings (productIDs)
	var productIDs []string
	err = json.Unmarshal(body, &productIDs)
	if err != nil {
		//fmt.Println("Error unmarshaling JSON:", err)
		return nil, err
	}
	//t(productIDs)
	var cartItems []*pb.CartItem
	for _, productID := range productIDs {
		cartItem := &pb.CartItem{
			ProductId: productID,
			Quantity:  int32(rand.Intn(10)),
		}
		cartItems = append(cartItems, cartItem)
	}

	return cartItems, nil
	
}


func GetPriceUsd(m *pb.Product) (*pb.Money) {
	if m != nil {
		return m.PriceUsd
	}
	return nil
}

func  GetProductId(m *pb.CartItem) string {
	if m != nil {
		return m.ProductId
	}
	return ""
}
func  prepOrderItems(items []*pb.CartItem, userCurrency string) ([]*pb.OrderItem, error) {
	out := make([]*pb.OrderItem, len(items))

	apiURL := "https://o346ng6ah7.execute-api.us-east-2.amazonaws.com/default/productcatalog"

	// JSON data to send in the request body

	for i, item := range items {
		data := map[string]interface{}{
			"type": "get",
			"id" : GetProductId(item),
		}
		jsonData, _ := json.Marshal(data)
		// Create a POST request with JSON data
		req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
		// Perform the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			//tln("Error making the request:", err)
			return nil, err
		}
		defer resp.Body.Close()
		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//fmt.Println("Error reading response body:", err)
			return nil, err
		}
		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
			
		}
		bodyJSON, _ := response["body"].(string)
		var bodyData map[string]interface{}
		if err := json.Unmarshal([]byte(bodyJSON), &bodyData); err != nil {
			//fmt.Println("Error parsing 'body' as JSON:", err)
			return nil, err
		}

		product := &pb.Product{}

		product.Id = bodyData["id"].(string)
		product.Name = bodyData["name"].(string)
		product.Description = bodyData["description"].(string)
		product.Picture = bodyData["picture"].(string)

		// Assuming "PriceUsd" is a complex structure, you need to further parse it from bodyData.
		priceUsd := &pb.Money{
			CurrencyCode: bodyData["price_usd"].(map[string]interface{})["currency_code"].(string),
			Units:       int64(bodyData["price_usd"].(map[string]interface{})["units"].(float64)),
			Nanos:       int32(bodyData["price_usd"].(map[string]interface{})["nanos"].(float64)),
		}

		product.PriceUsd = priceUsd

		// Categories should be converted to a slice of strings
		categoriesInterface := bodyData["categories"].([]interface{})
		categories := make([]string, len(categoriesInterface))
		for i, v := range categoriesInterface {
			categories[i] = v.(string)
		}
		product.Categories = categories

		price, err := convertCurrency(GetPriceUsd(product), userCurrency)
		if err != nil {
			return nil, fmt.Errorf("failed to convert price of %q to %s", item.GetProductId(), userCurrency)
		}
		out[i] = &pb.OrderItem{
			Item: item,
			Cost: price}
	}
	return out, nil
}

func convertCurrency(from *pb.Money, toCurrency string) (*pb.Money, error) {
	apiURL := "https://iks1v50mvk.execute-api.us-east-2.amazonaws.com/default/conversion"
	data := map[string]interface{}{
		"currency_code": from.CurrencyCode,
		"units": from.Units,
		"nanos": from.Nanos,
		"to_code": toCurrency,
	}
	jsonData, _ := json.Marshal(data)
	result, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, fmt.Errorf("failed to get product")
		}
	responseBody, err := ioutil.ReadAll(result.Body)
	var money struct {
		Units        int    `json:"units"`
		Nanos        int    `json:"nanos"`
		CurrencyCode string `json:"currency_code"`
	}
	
	if err := json.Unmarshal([]byte(responseBody), &money); err != nil {
		return nil, err
	}
	
	pbMoney := &pb.Money{
		CurrencyCode: money.CurrencyCode,
		Units:        int64(money.Units),
		Nanos:        int32(money.Nanos),
	}
	return pbMoney,nil

}

func chargeCard(amount *pb.Money, paymentInfo *pb.CreditCardInfo) (string, error) {
	apiURL := "https://kwrzd713al.execute-api.us-east-2.amazonaws.com/default/payment"
	data := map[string]interface{}{
		"amount": map[string]interface{}{
			"currency_code": amount.CurrencyCode,
			"units":        amount.Units,
			"nanos":        amount.Nanos,
		},
		"credit_card": map[string]interface{}{
			"credit_card_number":          paymentInfo.CreditCardNumber,
			"credit_card_expiration_month": paymentInfo.CreditCardExpirationMonth,
			"credit_card_expiration_year":  paymentInfo.CreditCardExpirationYear,
		},
	}	

	jsonData, _ := json.Marshal(data)
	fmt.Print("Value is being sent")
	fmt.Println(string(jsonData))
	result, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	result.Header.Set("Content-Type", "application/json")
	httpClient := &http.Client{}
	resp, _ := httpClient.Do(result)
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)

	htmlResponse := string(responseBody)
	return htmlResponse,nil
}

func sendOrderConfirmation(email string, order *pb.PlaceOrderResponse) string {
	data := map[string]interface{}{
		"email": email,
		"order": map[string]interface{}{
			"order_id":             order.Order.OrderId,
			"shipping_tracking_id": order.Order.ShippingTrackingId,
			"shipping_cost": map[string]interface{}{
				"units":         order.Order.ShippingCost.Units,
				"nanos":         order.Order.ShippingCost.Nanos,
				"currency_code": order.Order.ShippingCost.CurrencyCode,
			},
			"shipping_address": map[string]interface{}{
				"street_address_1": order.Order.ShippingAddress.StreetAddress,
				"street_address_2": "",
				"city":             order.Order.ShippingAddress.City,
				"country":          order.Order.ShippingAddress.Country,
				"zip_code":         order.Order.ShippingAddress.ZipCode,
			},
			"items": order.Order.Items,
		},
	}
	requestData, _ := json.Marshal(data)
	httpClient := &http.Client{}
	url := "https://tsvhk254jrn6hfqz7jo62wcufi0rmxzr.lambda-url.us-east-2.on.aws/"
	buf := bytes.NewBuffer(requestData)
	result, _ := http.NewRequest("POST", url, buf)

	result.Header.Set("Content-Type", "application/json")
	resp, _ := httpClient.Do(result)
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)

	htmlResponse := string(responseBody)
	//fmt.Println(htmlResponse)
	return htmlResponse
}


func  shipOrder(address *pb.Address, items []*pb.CartItem) (string, error) {
	apiURL := "https://tekw7om46d.execute-api.us-east-2.amazonaws.com/default/shipping"
	data := map[string]interface{}{
		"ShipOrder": map[string]interface{}{
			"Address": address,
			"Items":   items,
		},
		"GetQuote": map[string]interface{}{
			"Address": address,
			"Items":   items,
		},
	}	
	jsonData, err := json.Marshal(data)
	result, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return "failed",err
		}
	responseBody, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return "failed to read response body",err
	}
	
	var response struct {
		ShipOrderResponse struct {
			TrackingID string `json:"tracking_id"`
		} `json:"ShipOrderResponse"`
	}
	
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "failed to unmarshal response", err
	}
	
	trackingID := response.ShipOrderResponse.TrackingID
	return trackingID, nil
}
