package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Ad struct {
	RedirectURL string
	Text        string
}

type AdServiceLambda struct {
	AdsMap map[string][]Ad
}

func NewAdServiceLambda() *AdServiceLambda {
	return &AdServiceLambda{
		AdsMap: createAdsMap(),
	}
}

func createAdsMap() map[string][]Ad {
	adsMap := make(map[string][]Ad)

	hairdryer := Ad{
		RedirectURL: "/product/2ZYFJ3GM2N",
		Text:        "Hairdryer for sale. 50% off.",
	}

	tankTop := Ad{
		RedirectURL: "/product/66VCHSJNUP",
		Text:        "Tank top for sale. 20% off.",
	}

	candleHolder := Ad{
		RedirectURL: "/product/0PUK6V6EV0",
		Text:        "Candle holder for sale. 30% off.",
	}

	bambooGlassJar := Ad{
		RedirectURL: "/product/9SIQT8TOJO",
		Text:        "Bamboo glass jar for sale. 10% off.",
	}

	watch := Ad{
		RedirectURL: "/product/1YMWWN1N4O",
		Text:        "Watch for sale. Buy one, get second kit for free",
	}

	mug := Ad{
		RedirectURL: "/product/6E92ZMYYFZ",
		Text:        "Mug for sale. Buy two, get third one for free",
	}

	loafers := Ad{
		RedirectURL: "/product/L9ECAV7KIM",
		Text:        "Loafers for sale. Buy one, get second one for free",
	}

	adsMap["clothing"] = []Ad{tankTop}
	adsMap["accessories"] = []Ad{watch}
	adsMap["footwear"] = []Ad{loafers}
	adsMap["hair"] = []Ad{hairdryer}
	adsMap["decor"] = []Ad{candleHolder}
	adsMap["kitchen"] = []Ad{bambooGlassJar, mug}

	return adsMap
}

func (lambda *AdServiceLambda) HandleRequest(input interface{}) string {
	randomAd := lambda.getRandomAd()
	adText := randomAd.Text
	fmt.Println("Random Ad:", adText)
	return adText
}

func (lambda *AdServiceLambda) getRandomAd() Ad {
	rand.Seed(time.Now().UnixNano())
	category := getRandomCategory()
	ads := lambda.AdsMap[category]
	randomIndex := rand.Intn(len(ads))
	return ads[randomIndex]
}

func getRandomCategory() string {
	categories := []string{"clothing", "accessories", "footwear", "hair", "decor", "kitchen"}
	randomIndex := rand.Intn(len(categories))
	return categories[randomIndex]
}

func main() {
	lambda := NewAdServiceLambda()
	lambda.HandleRequest(nil)
}
