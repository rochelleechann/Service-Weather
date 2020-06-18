// Package p contains an HTTP Cloud Function.
package p

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	forecast "github.com/mlbright/darksky/v2"
)

type geolocation struct {
	Status  string `json:"status"`
	Results []struct {
		Geometry struct {
			Location struct {
				Lat  float64 `json:"lat"`
				Long float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

func Main(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for the preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	location, ok := r.URL.Query()["location"]

	if !ok || len(location[0]) < 1 {
		log.Printf("Url parameter 'location' is missing.")
		return
	}

	loc := location[0]

	convertEncodedLocation := url.QueryEscape(loc)
	fmt.Println("URL QueryEscape:", convertEncodedLocation)

	// Put together end point based on user input field
	googleBaseURL := "https://maps.googleapis.com/maps/api/geocode/json?address="

	// Removed API Key for now.
	googleEncodingKey := "&key=*****************************"
	googleEncodingEndPoint := googleBaseURL + convertEncodedLocation + googleEncodingKey

	// GET Request to Google Encoding API to retrieve latitude and longitude based on user input
	resp, err := http.Get(googleEncodingEndPoint)
	if err != nil {
		fmt.Println("Error making request to Google Encoding", err)
	}

	defer resp.Body.Close()
	fmt.Println("Response status:", resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading JSON data:", err)
	}

	responseBody := []byte(body)
	data := &geolocation{}
	json.Unmarshal(responseBody, &data)

	// Assign variable to pass data to get weather information
	latitude := data.Results[0].Geometry.Location.Lat
	longitude := data.Results[0].Geometry.Location.Long

	// Call DarkSky Function to Get Weather Info
	weatherInfo := getWeatherInfo(latitude, longitude)

	fmt.Println("Weather Info:", weatherInfo)

	convertJSON, err := json.Marshal(weatherInfo)
	if err != nil {
		fmt.Println("Error converting data to user object")
	}

	// Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(convertJSON)
}

func getWeatherInfo(lat float64, long float64) interface{} {
	log.Println("Called function to get weather information.")

	// Need to figure out where to secure this API key, should not be hard coded
	// Removed API Key for now
	darkSkyKey := "***************************"
	darkSkyKey = strings.TrimSpace(darkSkyKey)

	convertLatitude := fmt.Sprintf("%f", lat)
	convertLongitutde := fmt.Sprintf("%f", long)

	currentTime := time.Now()
	lastWeekTime := currentTime.AddDate(0, 0, -7).Unix()
	time := strconv.FormatInt(lastWeekTime, 10)

	f, err := forecast.Get(darkSkyKey, convertLatitude, convertLongitutde, time, forecast.CA, forecast.English)
	if err != nil {
		fmt.Println("Error could not get weather information", err)
	}

	return f
}
