package main

import (
	"fmt"
)

type CartItem struct {
	ProductId string
	Quantity  int32
}

func cart() []string {
	stringArray := []string{"66VCHSJNUP", "OLJCESPC7Z"}
	return stringArray
}

func main() {
	fmt.Print(cart())
}
