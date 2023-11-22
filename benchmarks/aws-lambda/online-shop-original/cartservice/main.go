package main

import (
	"math/rand"
	"time"
	"github.com/aws/aws-lambda-go/lambda"
)


func HandleLambdaEvent()([]string,error) {
	stringArray := []string{"66VCHSJNUP", "OLJCESPC7Z","1YMWWN1N4O","L9ECAV7KIM","2ZYFJ3GM2N","0PUK6V6EV0","LS4PSXUNUM","9SIQT8TOJO","6E92ZMYYFZ"}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(stringArray), func(i, j int) {
		stringArray[i], stringArray[j] = stringArray[j], stringArray[i]
	})
	randomIndex := rand.Intn(len(stringArray)) + 1 
	cartItemsList := stringArray[:randomIndex]
	return cartItemsList, nil

}

func main() {
	lambda.Start(HandleLambdaEvent)
}
