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
	"os/exec"
	"strings"
	// 	"strings"
	// 	"github.com/aws/aws-lambda-go/lambda"
	// 	"github.com/google/uuid"
)

func checkout(jsonArg string) string {
	//fmt.Println(jsonArg)
	cmd := exec.Command("./checkout/main", jsonArg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		//fmt.Printf("Error executing CartService: %v\n", err)
		return ""
	}
	outputString := string(output)
	//fmt.Println(outputString)
	return outputString
}

func ad() string {
	cmd := exec.Command("./adservice/adservice")
	output, err := cmd.CombinedOutput()
	if err != nil {
		//fmt.Printf("Error executing CartService: %v\n", err)
		return ""
	}
	outputString := string(output)
	//fmt.Println(outputString)
	return outputString
}

func list_product() string {
	cmd := exec.Command("./listproduct/productcatalog")
	output, err := cmd.CombinedOutput()
	if err != nil {
		//.Printf("Error executing CartService: %v\n", err)
		return ""
	}
	outputString := string(output)
	//fmt.Println(outputString)
	return outputString
}

func recommendation_python() string {
	data := map[string]interface{}{
		"type": "recommendation",
		"data": list_product(),
	}
	requestData, _ := json.Marshal(data)
	print(string(requestData))
	httpClient := &http.Client{}
	url := "https://bmk46xska6uzj4hhlwpcnhrdsi0osfnf.lambda-url.us-east-2.on.aws/"
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

func getUserCart(userID string) string {
	request := map[string]interface{}{
		"requestType": "get",
		"UserID":      userID,
	}
	requestData, err := json.Marshal(request)
	if err != nil {
		//fmt.Println("Error marshaling JSON:", err)
		return ""
	}

	jsonArg := string(requestData)
	cmd := exec.Command("./cart/main", jsonArg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		//fmt.Printf("Error executing CartService: %v\n", err)
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
			Quantity:  int32(rand.Intn(10)),
		}
		//fmt.Println(cartItem)
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
