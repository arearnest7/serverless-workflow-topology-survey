package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	// 	"strings"
	// 	"github.com/aws/aws-lambda-go/lambda"
	// 	"github.com/google/uuid"
)

// type OrderResult struct {
// 	OrderId            string
// 	ShippingTrackingId string
// 	ShippingCost       *Money
// 	ShippingAddress    *Address
// 	Items              []*OrderItem
// }

// type PlaceOrderResponse struct {
// 	Order *OrderResult
// }

// const (
// 	nanosMin = -999999999
// 	nanosMax = +999999999
// 	nanosMod = 1000000000
// )

// func (m *Money) GetCurrencyCode() string {
// 	if m != nil {
// 		return m.CurrencyCode
// 	}
// 	return ""
// }

// func (m *Money) GetUnits() int64 {
// 	if m != nil {
// 		return m.Units
// 	}
// 	return 0
// }

// func (m *Money) GetNanos() int32 {
// 	if m != nil {
// 		return m.Nanos
// 	}
// 	return 0
// }

// var (
// 	ErrInvalidValue        = errors.New("one of the specified money values is invalid")
// 	ErrMismatchingCurrency = errors.New("mismatching currency codes")
// )

// // IsValid checks if specified value has a valid units/nanos signs and ranges.
// func IsValid(m Money) bool {
// 	return signMatches(m) && validNanos(m.GetNanos())
// }

// func signMatches(m Money) bool {
// 	return m.GetNanos() == 0 || m.GetUnits() == 0 || (m.GetNanos() < 0) == (m.GetUnits() < 0)
// }

// func validNanos(nanos int32) bool { return nanosMin <= nanos && nanos <= nanosMax }

// // IsZero returns true if the specified money value is equal to zero.
// func IsZero(m Money) bool { return m.GetUnits() == 0 && m.GetNanos() == 0 }

// // IsPositive returns true if the specified money value is valid and is
// // positive.
// func IsPositive(m Money) bool {
// 	return IsValid(m) && m.GetUnits() > 0 || (m.GetUnits() == 0 && m.GetNanos() > 0)
// }

// // IsNegative returns true if the specified money value is valid and is
// // negative.
// func IsNegative(m Money) bool {
// 	return IsValid(m) && m.GetUnits() < 0 || (m.GetUnits() == 0 && m.GetNanos() < 0)
// }

// // AreSameCurrency returns true if values l and r have a currency code and
// // they are the same values.
// func AreSameCurrency(l, r Money) bool {
// 	return l.GetCurrencyCode() == r.GetCurrencyCode() && l.GetCurrencyCode() != ""
// }

// // AreEquals returns true if values l and r are the equal, including the
// // currency. This does not check validity of the provided values.
// func AreEquals(l, r Money) bool {
// 	return l.GetCurrencyCode() == r.GetCurrencyCode() &&
// 		l.GetUnits() == r.GetUnits() && l.GetNanos() == r.GetNanos()
// }

// // Negate returns the same amount with the sign negated.
// func Negate(m Money) Money {
// 	return Money{
// 		Units:        -m.GetUnits(),
// 		Nanos:        -m.GetNanos(),
// 		CurrencyCode: m.GetCurrencyCode()}
// }

// // Must panics if the given error is not nil. This can be used with other
// // functions like: "m := Must(Sum(a,b))".
// func Must(v Money, err error) Money {
// 	if err != nil {
// 		panic(err)
// 	}
// 	return v
// }

// // Sum adds two values. Returns an error if one of the values are invalid or
// // currency codes are not matching (unless currency code is unspecified for
// // both).
// func Sum(l, r Money) (Money, error) {
// 	if !IsValid(l) || !IsValid(r) {
// 		return Money{}, ErrInvalidValue
// 	} else if l.GetCurrencyCode() != r.GetCurrencyCode() {
// 		return Money{}, ErrMismatchingCurrency
// 	}
// 	units := l.GetUnits() + r.GetUnits()
// 	nanos := l.GetNanos() + r.GetNanos()

// 	if (units == 0 && nanos == 0) || (units > 0 && nanos >= 0) || (units < 0 && nanos <= 0) {
// 		// same sign <units, nanos>
// 		units += int64(nanos / nanosMod)
// 		nanos = nanos % nanosMod
// 	} else {
// 		// different sign. nanos guaranteed to not to go over the limit
// 		if units > 0 {
// 			units--
// 			nanos += nanosMod
// 		} else {
// 			units++
// 			nanos -= nanosMod
// 		}
// 	}

// 	return Money{
// 		Units:        units,
// 		Nanos:        nanos,
// 		CurrencyCode: l.GetCurrencyCode()}, nil
// }

// // MultiplySlow is a slow multiplication operation done through adding the value
// // to itself n-1 times.
// func MultiplySlow(m Money, n uint32) Money {
// 	out := m
// 	for n > 1 {
// 		out = Must(Sum(out, m))
// 		n--
// 	}
// 	return out
// }

// type Cart struct {
// 	UserId string      `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
// 	Items  []*CartItem `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
// }

// func PlaceOrder(req *PlaceOrderRequest) (string, error) {
// 	orderID, _ := uuid.NewUUID()
// 	//fmt.Print(orderID.String())
// 	address := Address{
// 		StreetAddress: req.Address.StreetAddress,
// 		City:          req.Address.City,
// 		State:         req.Address.State,
// 		Country:       req.Address.Country,
// 		ZipCode:       req.Address.ZipCode,
// 	}

// 	prep, _ := prepareOrderItemsAndShippingQuoteFromCart(req.UserId, req.UserCurrency, address)
// 	//shippingCostJSON, err := json.Marshal(prep.shippingCostLocalized)
// 	// if err != nil {
// 	// 	fmt.Println("Error marshaling Money to JSON:", err)
// 	// 	return
// 	// }
// 	//fmt.Printf("Shipping Cost (USD) as JSON: %s\n", string(shippingCostJSON))
// 	total := Money{CurrencyCode: req.UserCurrency,
// 		Units: 0,
// 		Nanos: 0}
// 	shippingCost := Money{
// 		CurrencyCode: prep.shippingCostLocalized.CurrencyCode,
// 		Units:        prep.shippingCostLocalized.Units,
// 		Nanos:        prep.shippingCostLocalized.Nanos,
// 	}
// 	// moneyStr := fmt.Sprintf("{ currency_code: '%s', units: %d, nanos: %d }", shippingCost.CurrencyCode, shippingCost.Units,shippingCost.Nanos)
// 	// fmt.Printf("Total amount: %s\n", moneyStr)
// 	total, _ = Sum(total, shippingCost)

// 	for _, it := range prep.orderItems {
// 		multPrice := MultiplySlow(*it.Cost, uint32(it.Item.Quantity))
// 		total = Must(Sum(total, multPrice))
// 	}
// 	// moneyStr := fmt.Sprintf("{ currency_code: '%s', units: %d, nanos: %d }", total.CurrencyCode, total.Units, total.Nanos)
// 	// fmt.Printf("Total amount: %s\n", moneyStr)
// 	creditCard := &CreditCard{
// 		CreditCardNumber:          req.CreditCard.CreditCardNumber,
// 		CreditCardCvv:             int32(req.CreditCard.CreditCardCvv),
// 		CreditCardExpirationYear:  int32(req.CreditCard.CreditCardExpirationYear),
// 		CreditCardExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
// 	}
// 	// fmt.Println("Credit card")
// 	// fmt.Print(creditCard)
// 	txID, err := chargeCard(&total, creditCard)
// 	fmt.Println(txID)
// 	if err != nil {
// 		return "", err
// 	}
// 	shippingTrackingID, err := shipOrder(&address, prep.cartItems)
// 	if err != nil {
// 		return "", err
// 	}
// 	//fmt.Println(shippingTrackingID)
// 	orderResult := &OrderResult{
// 		OrderId:            orderID.String(),
// 		ShippingTrackingId: shippingTrackingID,
// 		ShippingCost:       prep.shippingCostLocalized,
// 		ShippingAddress:    &address,
// 		Items:              prep.orderItems,
// 	}

// 	resp := &PlaceOrderResponse{Order: orderResult}
// 	jsonResp, err := json.Marshal(resp)
// 	if err != nil {
// 		fmt.Println("Error marshaling JSON:", err)
// 		return "", err
// 	}
// 	fmt.Println(string(jsonResp))
// 	//fmt.Println(sendOrderConfirmation(req.Email, resp))
// 	return sendOrderConfirmation(req.Email, resp),nil

// }

// type OrderItem struct {
// 	Item *CartItem
// 	Cost *Money
// }
// type CartItem struct {
// 	ProductId string
// 	Quantity  int32
// }
// type orderPrep struct {
// 	orderItems            []*OrderItem
// 	cartItems             []*CartItem
// 	shippingCostLocalized *Money
// }

// type Money struct {
// 	CurrencyCode string `json:"currency_code"`
// 	Units        int64  `json:"units"`
// 	Nanos        int32  `json:"nanos"`
// }

// type Product struct {
// 	Id          string `json:"id"`
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	Picture     string `json:"picture"`
// 	PriceUsd    *Money `json:"price_usd"`
// }

// func (c *CartItem) String() string {
// 	return fmt.Sprintf("ProductID: %s, Quantity: %d", c.ProductId, c.Quantity)
// }

// func prepareOrderItemsAndShippingQuoteFromCart(userID string, userCurrency string, address Address) (orderPrep, error) {
// 	var out orderPrep
// 	cartItems := getUserCart(userID)
// 	// for _, item := range cartItems {
// 	// 	println(item.String())
// 	// }
// 	orderItems, err := prepOrderItems(cartItems, userCurrency)
// 	if err != nil {
// 		return out, fmt.Errorf("prepaer: %+v", err)
// 	}
// 	// orderItemsJSON, _ := orderItemsToString(orderItems)
// 	// if err != nil {
// 	// 	fmt.Println("Error converting OrderItems to string:", err)
// 	// 	return nil,err
// 	// }
// 	// fmt.Println("OrderItems as a JSON string:")
// 	// fmt.Println(orderItemsJSON)
// 	shippingUSD, err := quoteShipping(&address, cartItems)
// 	if err != nil {
// 		return out, fmt.Errorf("shipping quote failure: %+v", err)
// 	}
// 	// fmt.Println("Shipping in USD")
// 	// fmt.Println(shippingUSD)
// 	shippingPrice, err := convertCurrency(shippingUSD, userCurrency)
// 	if err != nil {
// 		return out, fmt.Errorf("failed to convert shipping cost to currency: %+v", err)
// 	}
// 	// fmt.Println("ShippingPrice")
// 	// fmt.Println(shippingPrice)
// 	out.shippingCostLocalized = shippingPrice
// 	out.cartItems = cartItems
// 	out.orderItems = orderItems
// 	// fmt.Println("Order Prep Output")

// 	// fmt.Println("Cart Items:")
// 	// for _, item := range out.cartItems {
// 	//     fmt.Printf("ProductID: %s, Quantity: %d", item.ProductId, item.Quantity)
// 	// }

// 	// fmt.Println("Shipping Cost Localized: %v", out.shippingCostLocalized)
// 	return out, nil
// }

// type QuoteResponse struct {
// 	ShipOrderResponse struct {
// 		TrackingID string `json:"tracking_id"`
// 	} `json:"ShipOrderResponse"`
// 	GetQuoteResponse struct {
// 		CostUSD struct {
// 			CurrencyCode string `json:"currency_code"`
// 			Units        int64  `json:"units"`
// 			Nanos        int32  `json:"nanos"`
// 		} `json:"cost_usd"`
// 	} `json:"GetQuoteResponse"`
// }

// func quoteShipping(address *Address, items []*CartItem) (*Money, error) {
// 	addressMap := map[string]interface{}{
// 		"city":          address.City,
// 		"country":       address.Country,
// 		"state":         address.State,
// 		"zipCode":       address.ZipCode,
// 		"streetAddress": address.StreetAddress,
// 	}

// 	itemsMap := make([]map[string]interface{}, len(items))
// 	for i, item := range items {
// 		itemMap := map[string]interface{}{
// 			"product_id": item.ProductId,
// 			"quantity":   item.Quantity,
// 		}
// 		itemsMap[i] = itemMap
// 	}

// 	data := map[string]interface{}{
// 		"Address": addressMap,
// 		"Items":   itemsMap,
// 	}
// 	jsonBody, _ := json.Marshal(data)
// 	//fmt.Println(string(jsonBody))
// 	cmd := exec.Command("./shipping/shipping", string(jsonBody))
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Printf("Error executing CartService: %v\n", err)
// 		//return nil
// 	}
// 	//fmt.Printf("Shipping: %s\n", output)
// 	var response QuoteResponse
// 	err = json.Unmarshal([]byte(output), &response)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		//return nil,err
// 	}
// 	getQuoteResponse := response.GetQuoteResponse
// 	//fmt.Printf("Currency Code: %s, Units: %d, Nanos: %d\n", getQuoteResponse.CostUSD.CurrencyCode, getQuoteResponse. CostUSD.Units, getQuoteResponse.CostUSD.Nanos)
// 	resultMoney := &Money{
// 		CurrencyCode: getQuoteResponse.CostUSD.CurrencyCode,
// 		Units:        getQuoteResponse.CostUSD.Units,
// 		Nanos:        getQuoteResponse.CostUSD.Nanos,
// 	}
// 	return resultMoney, nil
// }

// func extractProductIDs(output string) []string {
// 	// Split the string by space and remove empty strings
// 	parts := strings.Fields(output)

// 	// Filter out any unwanted characters such as '[' and ']'
// 	var productIDs []string
// 	for _, part := range parts {
// 		cleanedPart := strings.Trim(part, "[]")
// 		productIDs = append(productIDs, cleanedPart)
// 	}

// 	return productIDs
// }

// func getUserCart(userID string) []*CartItem {
// 	request := map[string]interface{}{
// 		"requestType": "get",
// 		"UserID":      userID,
// 	}
// 	requestData, err := json.Marshal(request)
// 	if err != nil {
// 		fmt.Println("Error marshaling JSON:", err)
// 		return nil
// 	}

// 	jsonArg := string(requestData)
// 	cmd := exec.Command("./cart/cart", jsonArg)
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Printf("Error executing CartService: %v\n", err)
// 		return nil
// 	}
// 	//fmt.Printf("CartService output: %s\n", string(output))

// 	// Extract product IDs from the output
// 	outputString := string(output)
// 	productIDs := extractProductIDs(outputString)
// 	// for _, productID := range productIDs {
// 	// 	fmt.Println(productID)
// 	// }
// 	var cartItems []*CartItem
// 	for _, productID := range productIDs {
// 		cartItem := &CartItem{
// 			ProductId: productID,
// 			Quantity:  3,
// 		}
// 		cartItems = append(cartItems, cartItem)
// 	}
// 	return cartItems
// }

// // func emptyUserCart(var userID) {
// // 	resp := os.system("./CartService " + json.dumps({requestType: "empty", UserId: userID}))
// // 	return nil
// // }

// func GetProductId(m *CartItem) string {
// 	if m != nil {
// 		return m.ProductId
// 	}
// 	return ""
// }

// // func orderItemsToString(orderItems []*OrderItem) (string, error) {
// // 	orderItemsJSON, err := json.Marshal(orderItems)
// // 	if err != nil {
// // 		return "", err
// // 	}

// // 	return string(orderItemsJSON), nil
// // }

// func prepOrderItems(items []*CartItem, userCurrency string) ([]*OrderItem, error) {
// 	out := make([]*OrderItem, len(items))
// 	for i, item := range items {
// 		data := GetProductId(item)
// 		//fmt.Println(string(data))
// 		cmd := exec.Command("./productcatalog/productcatalog", string(data))
// 		output, err := cmd.CombinedOutput()
// 		if err != nil {
// 			fmt.Printf("Error executing ProductCatalog: %v\n", err)
// 			//return nil,err
// 		}
// 		product := &Product{}
// 		err = json.Unmarshal(output, product)
// 		if err != nil {
// 			fmt.Println("Error marshaling PriceUsd:", err)
// 			return nil, err
// 		}
// 		// Print the JSON representation of the PriceUsd field
// 		// fmt.Println("productcatalog")
// 		// fmt.Println("Product Price (USD):", string(output))
// 		price, err := convertCurrency(product.PriceUsd, userCurrency)
// 		//fmt.Println("ConvertCurrencyOutput")
// 		//fmt.Println(price)
// 		if err != nil {
// 			return nil, err
// 		}
// 		out[i] = &OrderItem{
// 			Item: item,
// 			Cost: price}
// 	}
// 	// orderItemsJSON, err := orderItemsToString(out)
// 	// if err != nil {
// 	// 	fmt.Println("Error converting OrderItems to string:", err)
// 	// 	return nil,err
// 	// }

// 	// fmt.Println("OrderItems as a JSON string:")
// 	// fmt.Println(orderItemsJSON)
// 	return out, nil
// }

// func convertCurrency(from *Money, toCurrency string) (*Money, error) {
// 	data := map[string]interface{}{
// 		"type": "currency",
// 		"body": map[string]interface{}{
// 			"currency_code": "USD",
// 			"nanos":         990000000,
// 			"requestType":   "convert",
// 			"to_code":       "EUR",
// 			"units":         19,
// 		},
// 	}
// 	requestData, err := json.Marshal(data)
// 	if err != nil {
// 		fmt.Println("Error marshaling JSON:", err)
// 		return nil, err
// 	}
// 	//fmt.Println("Conversion NodeJS")
// 	//fmt.Println(string(requestData))
// 	//fmt.Println(string(requestData))
// 	// cmd := exec.Command("node", "currency.js", string(requestData))
// 	// output, err := cmd.CombinedOutput()
// 	// if err != nil {
// 	// 	fmt.Printf("Error executing Currency: %v\n", err)
// 	// 	return nil, err
// 	// }
// 	httpClient := &http.Client{}
// 	//fmt.Println(string(requestData))
// 	url := "https://ofyvg75rl4vlqxhfjwc6adf73y0johle.lambda-url.us-east-2.on.aws/"
// 	buf := bytes.NewBuffer(requestData)
// 	result, err := http.NewRequest("POST", url, buf)
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return nil, err
// 	}
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get product")
// 	}
// 	result.Header.Set("Content-Type", "application/json")
// 	resp, err := httpClient.Do(result)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	responseBody, err := ioutil.ReadAll(resp.Body)
// 	//fmt.Print(string(responseBody))
// 	var money struct {
// 		Units        int    `json:"units"`
// 		Nanos        int    `json:"nanos"`
// 		CurrencyCode string `json:"currency_code"`
// 	}

// 	if err := json.Unmarshal([]byte(responseBody), &money); err != nil {
// 		return nil, err
// 	}

// 	moneyResult := &Money{
// 		CurrencyCode: money.CurrencyCode,
// 		Units:        int64(money.Units),
// 		Nanos:        int32(money.Nanos),
// 	}
// 	//fmt.Print(moneyResult)
// 	return moneyResult, nil

// }

// type CreditCard struct {
// 	CreditCardNumber          string `json:"creditCardNumber"`
// 	CreditCardCvv             int32  `json:"creditCardCvv"`
// 	CreditCardExpirationYear  int32  `json:"creditCardExpirationYear"`
// 	CreditCardExpirationMonth int32  `json:"creditCardExpirationMonth"`
// }

// func chargeCard(amount *Money, paymentInfo *CreditCard) (string, error) {
// 	data := map[string]interface{}{
// 		"type": "payment",
// 		"body": map[string]interface{}{
// 			"amount": map[string]interface{}{
// 				"currency_code": amount.CurrencyCode,
// 				"units":         amount.Units,
// 				"nanos":         amount.Nanos,
// 			},
// 			"credit_card": map[string]interface{}{
// 				"credit_card_number":           paymentInfo.CreditCardNumber,
// 				"credit_card_expiration_month": paymentInfo.CreditCardExpirationMonth,
// 				"credit_card_expiration_year":  paymentInfo.CreditCardExpirationYear,
// 			},
// 		},
// 	}
// 	requestData, err := json.Marshal(data)
// 	if err != nil {
// 		fmt.Println("Error marshaling JSON:", err)
// 		return "", err
// 	}
// 	httpClient := &http.Client{}
// 	//fmt.Println(string(requestData))
// 	url := "https://ofyvg75rl4vlqxhfjwc6adf73y0johle.lambda-url.us-east-2.on.aws/"
// 	buf := bytes.NewBuffer(requestData)
// 	result, err := http.NewRequest("POST", url, buf)
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return "", err
// 	}
// 	result.Header.Set("Content-Type", "application/json")
// 	resp, err := httpClient.Do(result)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()
// 	responseBody, _ := ioutil.ReadAll(resp.Body)
// 	return string(responseBody), nil

// }

// func sendOrderConfirmation(email string, order *PlaceOrderResponse) (string){
// 	data := map[string]interface{}{
// 		"type": "email",
// 		"data": map[string]interface{}{
// 			"email": email,
// 			"order": map[string]interface{}{
// 				"order_id":             order.Order.OrderId,
// 				"shipping_tracking_id": order.Order.ShippingTrackingId,
// 				"shipping_cost": map[string]interface{}{
// 					"units":         order.Order.ShippingCost.Units,
// 					"nanos":         order.Order.ShippingCost.Nanos,
// 					"currency_code": order.Order.ShippingCost.CurrencyCode,
// 				},
// 				"shipping_address": map[string]interface{}{
// 					"street_address_1": order.Order.ShippingAddress.StreetAddress,
// 					"street_address_2": "",
// 					"city":             order.Order.ShippingAddress.City,
// 					"country":          order.Order.ShippingAddress.Country,
// 					"zip_code":         order.Order.ShippingAddress.ZipCode,
// 				},
// 				"items": order.Order.Items,
// 			},
// 		},
// 	}
// 	requestData, _ := json.Marshal(data)
// 	fmt.Println(string(requestData))
// 	httpClient := &http.Client{}
// 	//fmt.Println(string(requestData))
// 	url := "https://bmk46xska6uzj4hhlwpcnhrdsi0osfnf.lambda-url.us-east-2.on.aws/"
// 	buf := bytes.NewBuffer(requestData)
// 	result, err := http.NewRequest("POST", url, buf)
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 	}
// 	result.Header.Set("Content-Type", "application/json")
// 	resp, _ := httpClient.Do(result)
// 	defer resp.Body.Close()

// 	responseBody, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response body:", err)
// 	}

// 	htmlResponse := string(responseBody)
// 	//fmt.Println(htmlResponse)
// 	return htmlResponse
// }

// func shipOrder(address *Address, items []*CartItem) (string, error) {
// 	addressMap := map[string]interface{}{
// 		"city":          address.City,
// 		"country":       address.Country,
// 		"state":         address.State,
// 		"zipCode":       address.ZipCode,
// 		"streetAddress": address.StreetAddress,
// 	}

// 	itemsMap := make([]map[string]interface{}, len(items))
// 	for i, item := range items {
// 		itemMap := map[string]interface{}{
// 			"product_id": item.ProductId,
// 			"quantity":   item.Quantity,
// 		}
// 		itemsMap[i] = itemMap
// 	}

// 	data := map[string]interface{}{
// 		"Address": addressMap,
// 		"Items":   itemsMap,
// 	}
// 	jsonBody, _ := json.Marshal(data)
// 	//fmt.Println(string(jsonBody))
// 	cmd := exec.Command("./shipping/shipping", string(jsonBody))
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Printf("Error executing CartService: %v\n", err)
// 		//return nil
// 	}
// 	//fmt.Println(string(output))
// 	var ship struct {
// 		ShipOrderResponse struct {
// 			TrackingID string `json:"TrackingId"`
// 		} `json:"ShipOrderResponse"`
// 	}

// 	if err := json.Unmarshal([]byte(output), &ship); err != nil {
// 		return "failed to unmarshal response", err
// 	}
// 	//fmt.Println(ship)
// 	trackingID := ship.ShipOrderResponse.TrackingID
// 	return trackingID, nil
// }

// type PlaceOrderRequest struct {
// 	UserId       string          `json:"userId"`
// 	UserCurrency string          `json:"userCurrency"`
// 	Address      *Address        `json:"address"`
// 	Email        string          `json:"email"`
// 	CreditCard   *CreditCardInfo `json:"creditCard"`
// }

// type MyEvent struct {
// 	UserId       string          `json:"userId"`
// 	UserCurrency string          `json:"userCurrency"`
// 	Address      *Address        `json:"address"`
// 	Email        string          `json:"email"`
// 	CreditCard   *CreditCardInfo `json:"creditCard"`
// }

// type Address struct {
// 	StreetAddress string `json:"streetAddress"`
// 	City          string `json:"city"`
// 	State         string `json:"state"`
// 	Country       string `json:"country"`
// 	ZipCode       int32  `json:"zipCode"`
// }

// type CreditCardInfo struct {
// 	CreditCardNumber          string `json:"creditCardNumber"`
// 	CreditCardCvv             int32  `json:"creditCardCvv"`
// 	CreditCardExpirationYear  int32  `json:"creditCardExpirationYear"`
// 	CreditCardExpirationMonth int32  `json:"creditCardExpirationMonth"`
// }

func checkout(jsonArg string) string {
	fmt.Println(jsonArg)
	cmd := exec.Command("./checkout/checkout", jsonArg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing CartService: %v\n", err)
		return ""
	}
	outputString := string(output)
	fmt.Println(outputString)
	return outputString
}

func ad() string {
	cmd := exec.Command("./adservice/adservice")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing CartService: %v\n", err)
		return ""
	}
	outputString := string(output)
	fmt.Println(outputString)
	return outputString
}

func list_product() string {
	cmd := exec.Command("./listproduct/productcatalog")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing CartService: %v\n", err)
		return ""
	}
	outputString := string(output)
	fmt.Println(outputString)
	return outputString
}

func recommendation_python() string {
	data := map[string]interface{}{
		"type": "recommendation",
		"data" : list_product(),
	}
	requestData, _ := json.Marshal(data)
	print(string(requestData))
	httpClient := &http.Client{}
	url := "https://bmk46xska6uzj4hhlwpcnhrdsi0osfnf.lambda-url.us-east-2.on.aws/"
	buf := bytes.NewBuffer(requestData)
	result, err := http.NewRequest("POST", url, buf)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "nil"
	}
	if err != nil {
		return ""
	}
	result.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(result)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	return string(responseBody)
}
type CartItem struct {
	ProductId string
	Quantity  int32
}

func extractProductIDs(output string) []string {
	// Split the string by space and remove empty strings
	parts := strings.Fields(output)

	// Filter out any unwanted characters such as '[' and ']'
	var productIDs []string
	for _, part := range parts {
		cleanedPart := strings.Trim(part, "[]")
		productIDs = append(productIDs, cleanedPart)
	}

	return productIDs
}
func (c *CartItem) String() string {
	return fmt.Sprintf("ProductID: %s, Quantity: %d", c.ProductId, c.Quantity)
}

func getUserCart(userID string) (string) {
	request := map[string]interface{}{
		"requestType": "get",
		"UserID":      userID,
	}
	requestData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return ""
	}

	jsonArg := string(requestData)
	cmd := exec.Command("./cart/cart", jsonArg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing CartService: %v\n", err)
		return ""
	}
	//fmt.Printf("CartService output: %s\n", string(output))

	// Extract product IDs from the output
	outputString := string(output)
	productIDs := extractProductIDs(outputString)
	// for _, productID := range productIDs {
	// 	fmt.Println(productID)
	// }
	var cartItems []*CartItem
	for _, productID := range productIDs {
		cartItem := &CartItem{
			ProductId: productID,
			Quantity:  3,
		}
		cartItems = append(cartItems, cartItem)
	}
	string_to_return := ""
	for _, item := range cartItems {
		string_to_return += "\n"
		string_to_return += item.String()
	}
	return string_to_return
}


type MyEvent struct {
	Call string `json:"call"`
	Body string `json:"body"`
}

func HandleLambdaEvent(myEvent MyEvent) (string, error) {
	fmt.Println(myEvent)
	type_call := myEvent.Call
	request_body := myEvent.Body
	fmt.Println(type_call)
	fmt.Println(request_body)
	switch type_call {
	case "checkout":
		result := checkout(request_body)
		return result, nil
	case "ad":
		result := ad()
		return result, nil
	case "list":
		result := list_product()
		return result, nil
	case "recommend":
		result := recommendation_python()
		return result, nil
	case "cart":
		result := getUserCart("1234")
		return result,nil
	default:
		fmt.Println("Invalid")
		return "Invalid", nil
	}

}

func main() {
	lambda.Start(HandleLambdaEvent)
}

// func main() {
//     type_call := os.Args[1]
// 	request_body := os.Args[2]
//     switch type_call {
//     case "checkout":
//         checkout(string(request_body))
//     case "ad":
//         ad()
// 	case "list":
// 		list_product()
// 	case "recommend":
// 		recommendation_python(request_body)
//     default:
//         fmt.Println("It's a regular day.")
//     }
// }
