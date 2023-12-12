package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"math/rand"
	"net/http"
)

const (
	usdCurrency = "USD"
)

type OrderResult struct {
	OrderId            string
	ShippingTrackingId string
	ShippingCost       *Money
	ShippingAddress    *Address
	Items              []*OrderItem
}

type PlaceOrderResponse struct {
	Order *OrderResult
}

const (
	nanosMin = -999999999
	nanosMax = +999999999
	nanosMod = 1000000000
)

func (m *Money) GetCurrencyCode() string {
	if m != nil {
		return m.CurrencyCode
	}
	return ""
}

func (m *Money) GetUnits() int64 {
	if m != nil {
		return m.Units
	}
	return 0
}

func (m *Money) GetNanos() int32 {
	if m != nil {
		return m.Nanos
	}
	return 0
}

var (
	ErrInvalidValue        = errors.New("one of the specified money values is invalid")
	ErrMismatchingCurrency = errors.New("mismatching currency codes")
)

// IsValid checks if specified value has a valid units/nanos signs and ranges.
func IsValid(m Money) bool {
	return signMatches(m) && validNanos(m.GetNanos())
}

func signMatches(m Money) bool {
	return m.GetNanos() == 0 || m.GetUnits() == 0 || (m.GetNanos() < 0) == (m.GetUnits() < 0)
}

func validNanos(nanos int32) bool { return nanosMin <= nanos && nanos <= nanosMax }

// IsZero returns true if the specified money value is equal to zero.
func IsZero(m Money) bool { return m.GetUnits() == 0 && m.GetNanos() == 0 }

// IsPositive returns true if the specified money value is valid and is
// positive.
func IsPositive(m Money) bool {
	return IsValid(m) && m.GetUnits() > 0 || (m.GetUnits() == 0 && m.GetNanos() > 0)
}

// IsNegative returns true if the specified money value is valid and is
// negative.
func IsNegative(m Money) bool {
	return IsValid(m) && m.GetUnits() < 0 || (m.GetUnits() == 0 && m.GetNanos() < 0)
}

// AreSameCurrency returns true if values l and r have a currency code and
// they are the same values.
func AreSameCurrency(l, r Money) bool {
	return l.GetCurrencyCode() == r.GetCurrencyCode() && l.GetCurrencyCode() != ""
}

// AreEquals returns true if values l and r are the equal, including the
// currency. This does not check validity of the provided values.
func AreEquals(l, r Money) bool {
	return l.GetCurrencyCode() == r.GetCurrencyCode() &&
		l.GetUnits() == r.GetUnits() && l.GetNanos() == r.GetNanos()
}

// Negate returns the same amount with the sign negated.
func Negate(m Money) Money {
	return Money{
		Units:        -m.GetUnits(),
		Nanos:        -m.GetNanos(),
		CurrencyCode: m.GetCurrencyCode()}
}

// Must panics if the given error is not nil. This can be used with other
// functions like: "m := Must(Sum(a,b))".
func Must(v Money, err error) Money {
	if err != nil {
		panic(err)
	}
	return v
}

// Sum adds two values. Returns an error if one of the values are invalid or
// currency codes are not matching (unless currency code is unspecified for
// both).
func Sum(l, r Money) (Money, error) {
	if !IsValid(l) || !IsValid(r) {
		return Money{}, ErrInvalidValue
	} else if l.GetCurrencyCode() != r.GetCurrencyCode() {
		return Money{}, ErrMismatchingCurrency
	}
	units := l.GetUnits() + r.GetUnits()
	nanos := l.GetNanos() + r.GetNanos()

	if (units == 0 && nanos == 0) || (units > 0 && nanos >= 0) || (units < 0 && nanos <= 0) {
		// same sign <units, nanos>
		units += int64(nanos / nanosMod)
		nanos = nanos % nanosMod
	} else {
		// different sign. nanos guaranteed to not to go over the limit
		if units > 0 {
			units--
			nanos += nanosMod
		} else {
			units++
			nanos -= nanosMod
		}
	}

	return Money{
		Units:        units,
		Nanos:        nanos,
		CurrencyCode: l.GetCurrencyCode()}, nil
}

// MultiplySlow is a slow multiplication operation done through adding the value
// to itself n-1 times.
func MultiplySlow(m Money, n uint32) Money {
	out := m
	for n > 1 {
		out = Must(Sum(out, m))
		n--
	}
	return out
}

type OrderItem struct {
	Item *CartItem
	Cost *Money
}
type CartItem struct {
	ProductId string
	Quantity  int32
}
type orderPrep struct {
	orderItems            []*OrderItem
	cartItems             []*CartItem
	shippingCostLocalized *Money
}

type Money struct {
	CurrencyCode string `json:"currency_code"`
	Units        int64  `json:"units"`
	Nanos        int32  `json:"nanos"`
}

type Product struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Picture     string   `json:"picture"`
	PriceUsd    *Money   `json:"price_usd"`
	Categories  []string `json:"categories"`
}

type Cart struct {
	UserId string      `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Items  []*CartItem `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
}

func GetItem(m *OrderItem) *CartItem {
	if m != nil {
		return m.Item
	}
	return nil
}

func GetQuantity(m *CartItem) int32 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

type MyEvent struct {
	UserId       string         `json:"userId"`
	UserCurrency string         `json:"userCurrency"`
	Address      Address        `json:"address"`
	Email        string         `json:"email"`
	CreditCard   CreditCardInfo `json:"creditCard"`
}

type Address struct {
	StreetAddress string `json:"streetAddress"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	ZipCode       int32  `json:"zipCode"`
}

type CreditCardInfo struct {
	CreditCardNumber          string `json:"creditCardNumber"`
	CreditCardCvv             int32  `json:"creditCardCvv"`
	CreditCardExpirationYear  int32  `json:"creditCardExpirationYear"`
	CreditCardExpirationMonth int32  `json:"creditCardExpirationMonth"`
}

type PlaceOrderRequest struct {
	UserId       string          `json:"userId"`
	UserCurrency string          `json:"userCurrency"`
	Address      *Address        `json:"address"`
	Email        string          `json:"email"`
	CreditCard   *CreditCardInfo `json:"creditCard"`
}

func HandleLambdaEvent(event events.APIGatewayProxyRequest) (string, error) {
	var myEvent PlaceOrderRequest
	fmt.Println(string(event.Body))
	if err := json.Unmarshal([]byte(event.Body), &myEvent); err != nil {
		return "", err
	}
	fmt.Print(string(myEvent.Address.StreetAddress))
	result, _ := PlaceOrder(&myEvent)
	return result, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}

func PlaceOrder(req *PlaceOrderRequest) (string, error) {
	orderID, err := uuid.NewUUID()
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to generate order uuid")
	}

	prep, err := prepareOrderItemsAndShippingQuoteFromCart(req.UserId, req.UserCurrency, req.Address)
	if err != nil {
		return "", status.Errorf(codes.Internal, err.Error())
	}
	total := Money{CurrencyCode: req.UserCurrency,
		Units: 0,
		Nanos: 0}
	total = Must(Sum(total, *prep.shippingCostLocalized))
	for _, it := range prep.orderItems {
		multPrice := MultiplySlow(*it.Cost, uint32(GetQuantity(GetItem(it))))
		total = Must(Sum(total, multPrice))
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

	orderResult := &OrderResult{
		OrderId:            orderID.String(),
		ShippingTrackingId: shippingTrackingID,
		ShippingCost:       prep.shippingCostLocalized,
		ShippingAddress:    req.Address,
		Items:              prep.orderItems,
	}

	resp := &PlaceOrderResponse{Order: orderResult}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	fmt.Println(string(jsonResp))
	fmt.Println(sendOrderConfirmation(req.Email, resp))
	return txID + string(string(jsonResp)) + sendOrderConfirmation(req.Email, resp), nil
}

func prepareOrderItemsAndShippingQuoteFromCart(userID, userCurrency string, address *Address) (orderPrep, error) {
	var out orderPrep
	cartItems, err := getUserCart(userID)
	if err != nil {
		return out, fmt.Errorf("cart failure: %+v", err)
	}
	orderItems, err := prepOrderItems(cartItems, userCurrency)
	if err != nil {
		return out, fmt.Errorf("failed to prepare order: %+v", err)
	}
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
			Units        int64  `json:"units"`
			Nanos        int32  `json:"nanos"`
		} `json:"cost_usd"`
	} `json:"GetQuoteResponse"`
}

func quoteShipping(address *Address, items []*CartItem) (*Money, error) {
	apiURL := "https://tekw7om46d.execute-api.us-east-2.amazonaws.com/default/shipping"
	addressMap := map[string]interface{}{
		"city":          address.City,
		"country":       address.Country,
		"state":         address.State,
		"zipCode":       address.ZipCode,
		"streetAddress": address.StreetAddress,
	}

	itemsMap := make([]map[string]interface{}, len(items))
	for i, item := range items {
		itemMap := map[string]interface{}{
			"product_id": item.ProductId,
			"quantity":   item.Quantity,
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
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
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
		return nil, err
	}
	getQuoteResponse := response.GetQuoteResponse
	//fmt.Printf("Currency Code: %s, Units: %d, Nanos: %d\n", getQuoteResponse.CostUSD.CurrencyCode, getQuoteResponse. CostUSD.Units, getQuoteResponse.CostUSD.Nanos)
	resultMoney := &Money{
		CurrencyCode: getQuoteResponse.CostUSD.CurrencyCode,
		Units:        getQuoteResponse.CostUSD.Units,
		Nanos:        getQuoteResponse.CostUSD.Nanos,
	}

	return resultMoney, nil
}

func getUserCart(userID string) ([]*CartItem, error) {
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
	var cartItems []*CartItem
	for _, productID := range productIDs {
		cartItem := &CartItem{
			ProductId: productID,
			Quantity:  int32(rand.Intn(10)),
		}
		cartItems = append(cartItems, cartItem)
	}

	return cartItems, nil

}

func GetPriceUsd(m *Product) *Money {
	if m != nil {
		return m.PriceUsd
	}
	return nil
}

func GetProductId(m *CartItem) string {
	if m != nil {
		return m.ProductId
	}
	return ""
}
func prepOrderItems(items []*CartItem, userCurrency string) ([]*OrderItem, error) {
	out := make([]*OrderItem, len(items))

	apiURL := "https://o346ng6ah7.execute-api.us-east-2.amazonaws.com/default/productcatalog"

	// JSON data to send in the request body

	for i, item := range items {
		data := map[string]interface{}{
			"type": "get",
			"id":   GetProductId(item),
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

		product := &Product{}

		product.Id = bodyData["id"].(string)
		product.Name = bodyData["name"].(string)
		product.Description = bodyData["description"].(string)
		product.Picture = bodyData["picture"].(string)

		// Assuming "PriceUsd" is a complex structure, you need to further parse it from bodyData.
		priceUsd := &Money{
			CurrencyCode: bodyData["price_usd"].(map[string]interface{})["currency_code"].(string),
			Units:        int64(bodyData["price_usd"].(map[string]interface{})["units"].(float64)),
			Nanos:        int32(bodyData["price_usd"].(map[string]interface{})["nanos"].(float64)),
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
			return nil, fmt.Errorf("failed to convert price of %q to %s", GetProductId(item), userCurrency)
		}
		out[i] = &OrderItem{
			Item: item,
			Cost: price}
	}
	return out, nil
}

func convertCurrency(from *Money, toCurrency string) (*Money, error) {
	apiURL := "https://iks1v50mvk.execute-api.us-east-2.amazonaws.com/default/conversion"
	data := map[string]interface{}{
		"currency_code": from.CurrencyCode,
		"units":         from.Units,
		"nanos":         from.Nanos,
		"to_code":       toCurrency,
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

	pbMoney := &Money{
		CurrencyCode: money.CurrencyCode,
		Units:        int64(money.Units),
		Nanos:        int32(money.Nanos),
	}
	return pbMoney, nil

}

func chargeCard(amount *Money, paymentInfo *CreditCardInfo) (string, error) {
	apiURL := "https://kwrzd713al.execute-api.us-east-2.amazonaws.com/default/payment"
	data := map[string]interface{}{
		"amount": map[string]interface{}{
			"currency_code": amount.CurrencyCode,
			"units":         amount.Units,
			"nanos":         amount.Nanos,
		},
		"credit_card": map[string]interface{}{
			"credit_card_number":           paymentInfo.CreditCardNumber,
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
	return htmlResponse, nil
}

func sendOrderConfirmation(email string, order *PlaceOrderResponse) string {
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

func shipOrder(address *Address, items []*CartItem) (string, error) {
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
		return "failed", err
	}
	responseBody, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return "failed to read response body", err
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
