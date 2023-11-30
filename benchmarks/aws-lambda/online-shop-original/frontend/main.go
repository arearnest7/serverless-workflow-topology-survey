package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	// 	"strings"
	// 	"github.com/aws/aws-lambda-go/lambda"
	// 	"github.com/google/uuid"
)

func checkout(jsonArg string) string {
	//fmt.Println(jsonArg)
	httpClient := &http.Client{}
	url := "https://de6kmyppcoefrusdfmabvsrczi0ucywk.lambda-url.us-east-2.on.aws/"
	buf := bytes.NewBuffer([]byte(jsonArg))
	result, err := http.NewRequest("POST", url, buf)
	if err != nil {
		//fmt.Println("Error creating request:", err)
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
	outputString := string(responseBody)
	return outputString
}

func ad() string {
	httpClient := &http.Client{}
	url := "https://mdgjmy76f67xpk4c2ojauf4aru0hrxxj.lambda-url.us-east-2.on.aws/"
	buf := bytes.NewBuffer(nil)
	result, err := http.NewRequest("POST", url, buf)
	if err != nil {
		//fmt.Println("Error creating request:", err)
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
	outputString := string(responseBody)
	return outputString

}

func list_product() string {
	httpClient := &http.Client{}
	url := "https://hmcaeu5brplevg6qpwax4aolqy0wydpv.lambda-url.us-east-2.on.aws/"
	buf := bytes.NewBuffer(nil)
	result, err := http.NewRequest("POST", url, buf)
	if err != nil {
		//fmt.Println("Error creating request:", err)
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
	outputString := string(responseBody)
	return outputString
}

func recommendation_python() string {
	data := map[string]interface{}{
		"data": list_product(),
	}
	requestData, _ := json.Marshal(data)
	//print(string(requestData))
	httpClient := &http.Client{}
	url := "https://s3rzleoznfrs2fvtvklv5kygjm0vekpx.lambda-url.us-east-2.on.aws/"
	buf := bytes.NewBuffer(requestData)
	result, err := http.NewRequest("POST", url, buf)
	if err != nil {
		//fmt.Println("Error creating request:", err)
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
	re := regexp.MustCompile(`"([^"]+)"`)
	matches := re.FindAllStringSubmatch(output, -1)

	var productIDs []string
	for _, match := range matches {
		productIDs = append(productIDs, match[1])
	}

	//fmt.Println(len(productIDs))
	return productIDs
}
func (c *CartItem) String() string {
	return fmt.Sprintf("ProductID: %s, Quantity: %d", c.ProductId, c.Quantity)
}

func getUserCart(userID string) string {
	request := map[string]interface{}{
		"userID": userID,
	}
	requestData, err := json.Marshal(request)
	if err != nil {
		//fmt.Println("Error marshaling JSON:", err)
		return ""
	}
	//fmt.Println(string(requestData))
	httpClient := &http.Client{}
	url := "https://c60ekgpu4l.execute-api.us-east-2.amazonaws.com/default/cartservice"
	buf := bytes.NewBuffer(requestData)
	result, err := http.NewRequest("POST", url, buf)
	if err != nil {
		//fmt.Println("Error creating request:", err)
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
	outputString := string(responseBody)
	//mt.Println(outputString)
	productIDs := extractProductIDs(outputString)
	//fmt.Println(productIDs)
	var cartItems []*CartItem
	for _, productID := range productIDs {
		cartItem := &CartItem{
			ProductId: productID,
			Quantity:  int32(rand.Intn(10)),
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

func HandleLambdaEvent(event events.APIGatewayProxyRequest) (string, error) {
	var myEvent MyEvent
	if err := json.Unmarshal([]byte(event.Body), &myEvent); err != nil {
		//fmt.Println("Error parsing event body:", err)
		return "",err
	}
	//fmt.Println(myEvent)
	type_call := myEvent.Call
	request_body := myEvent.Body
	//fmt.Println(type_call)
	//fmt.Println(request_body)
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
		return result, nil
	default:
		//fmt.Println("Invalid")
		return "Invalid", nil
	}

}

func main() {
	lambda.Start(HandleLambdaEvent)
}

