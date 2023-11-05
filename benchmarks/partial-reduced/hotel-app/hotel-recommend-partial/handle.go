package function

import (
	"context"
	"fmt"
	"net/http"
	"math"
	"os"
	"encoding/json"
	"io/ioutil"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/sirupsen/logrus"

	"github.com/hailocab/go-geoindex"
)

type RequestBody struct {
        request string "json:\"request\""
        Lat float64 "json:\"Lat,omitempty\""
        Lon float64 "json:\"Lon,omitempty\""
        HotelId string "json:\"HotelId,omitempty\""
        HotelIds []string "json:\"HotelIds,omitempty\""
        RoomNumber int "json:\"RoomNumber,omitempty\""
        CustomerName string "json:\"CustomerName,omitempty\""
        Username string "json:\"Username,omitempty\""
        Password string "json:\"Password,omitempty\""
        Require string "json:\"Require,omitempty\""
        InDate string "json:\"InDate,omitempty\""
        OutDate string "json:\"OutDate,omitempty\""
}

type Hotel struct {
	ID     bson.ObjectId `bson:"_id"`
	HId    string        `bson:"hotelid"`
	HLat   float64       `bson:"lat"`
	HLon   float64       `bson:"lon"`
	HRate  float64       `bson:"rate"`
	HPrice float64       `bson:"price"`
}

// loadRecommendations loads hotel recommendations from mongodb.
func loadRecommendations(session *mgo.Session) map[string]Hotel {
	// session, err := mgo.Dial("mongodb-recommendation")
	// if err != nil {
	// 	panic(err)
	// }
	// defer session.Close()
	s := session.Copy()
	defer s.Close()

	c := s.DB("recommendation-db").C("recommendation")

	// unmarshal json profiles
	var hotels []Hotel
	err := c.Find(bson.M{}).All(&hotels)
	if err != nil {
		log.Println("Failed get hotels data: ", err)
	}

	profiles := make(map[string]Hotel)
	for _, hotel := range hotels {
		profiles[hotel.HId] = hotel
	}

	return profiles
}

// GiveRecommendation returns recommendations within a given requirement.
func GetRecommendations(req RequestBody) string {
	MongoSession, _ := mgo.Dial(os.Getenv("HOTEL_APP_DATABASE"))
	var hotels map[string]Hotel
	hotels = loadRecommendations(MongoSession)

	res := make([]string, 0)
	// fmt.Printf("GetRecommendations:\n")
	// fmt.Printf("%+v\n", s.hotels)

	require := req.Require
	if require == "dis" {
		p1 := &geoindex.GeoPoint{
			Pid:  "",
			Plat: req.Lat,
			Plon: req.Lon,
		}
		min := math.MaxFloat64
		for _, hotel := range hotels {
			tmp := float64(geoindex.Distance(p1, &geoindex.GeoPoint{
				Pid:  "",
				Plat: hotel.HLat,
				Plon: hotel.HLon,
			})) / 1000
			if tmp < min {
				min = tmp
			}
		}
		for _, hotel := range hotels {
			tmp := float64(geoindex.Distance(p1, &geoindex.GeoPoint{
				Pid:  "",
				Plat: hotel.HLat,
				Plon: hotel.HLon,
			})) / 1000
			if tmp == min {
				res = append(res, hotel.HId)
			}
		}
	} else if require == "rate" {
		max := 0.0
		for _, hotel := range hotels {
			if hotel.HRate > max {
				max = hotel.HRate
			}
		}
		for _, hotel := range hotels {
			if hotel.HRate == max {
				res = append(res, hotel.HId)
			}
		}
	} else if require == "price" {
		min := math.MaxFloat64
		for _, hotel := range hotels {
			if hotel.HPrice < min {
				min = hotel.HPrice
			}
		}
		for _, hotel := range hotels {
			if hotel.HPrice == min {
				res = append(res, hotel.HId)
			}
		}
	} else {
		log.Println("Wrong parameter: ", require)
	}

	ret, _ := json.Marshal(res)
        return string(ret)
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
        body, _ := ioutil.ReadAll(req.Body)
        var body_u *RequestBody
        json.Unmarshal(body, &body_u)
        defer req.Body.Close()
	fmt.Fprintf(res, GetRecommendations(*body_u)) // echo to caller
}
