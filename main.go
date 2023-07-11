package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	}
}

type Location struct {
	Candidates []struct {
		Formatted_address string `json:"formatted_address"`
		Geometry          struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"candidates"`
}

func main() {

	name := "スカイツリー"

	if len(os.Args) >= 2 {
		name = os.Args[1]
	}

	//google map api
	res, err := http.Get(
		"https://maps.googleapis.com/maps/api/place/findplacefromtext/json" +
			"?key=APIKEY&input=" + name +
			"&inputtype=textquery&fields=plus_code,photos,formatted_address,name,geometry")

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var location Location

	err = json.Unmarshal(body, &location)
	if err != nil {
		panic(err)
	}

	lat, lng := location.Candidates[0].Geometry.Location.Lat, location.Candidates[0].Geometry.Location.Lng

	latToStr := strconv.FormatFloat(lat, 'f', -1, 64)
	lngToStr := strconv.FormatFloat(lng, 'f', -1, 64)

	//weather api

	res, err = http.Get("http://api.weatherapi.com/v1/forecast.json?key=APIKEY=" + latToStr + "," + lngToStr + "&days=1&aqi=no&alerts=no")

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather

	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	locationWeather, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Printf(
		"%s, %s:, %.0fC, %s\n",
		locationWeather.Name,
		locationWeather.Country,
		current.TempC,
		current.Condition.Text,
	)

	for _, hour := range hours {

		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		switch hour.Condition.Text {
		case "Partly cloudy":
			hour.Condition.Text = "所により曇"
		case "Cloudy":
			hour.Condition.Text = "曇"
		case "Patchy rain possible":
			hour.Condition.Text = "所により雨"
		case "Overcast":
			hour.Condition.Text = "どんよりとした空模様"
		case "Heavy rain":
			hour.Condition.Text = "大雨"
		case "Light rain shower":
			hour.Condition.Text = "わか雨"
		case "Moderate rain":
			hour.Condition.Text = "雨"
		case "Clear":
			hour.Condition.Text = "快晴"
		case "Mist":
			hour.Condition.Text = "靄"
		case "Fog":
			hour.Condition.Text = "霧"
		case "Sunny":
			hour.Condition.Text = "晴れ"
		case "Light rain":
			hour.Condition.Text = "小雨"
		case "Patchy light drizzle":
			hour.Condition.Text = "所により小雨"
		case "Light drizzle":
			hour.Condition.Text = "小雨"
		case "Moderate or heavy rain shower":
			hour.Condition.Text = "激しい雨nなどに注意"
		case "Heavy rain at times":
			hour.Condition.Text = "時折激しい雨"
		case "Patchy light rain with thunder":
			hour.Condition.Text = "雷を伴う雨"
		case "thunder":
			hour.Condition.Text = "雷"
		case "Thundery outbreaks possible":
			hour.Condition.Text = "雷が発生する可能性あり"
		}

		data := fmt.Sprintf(
			"%s - %.0fC, 降水確率:%.0f, %s\n",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		if hour.ChanceOfRain < 40 {
			fmt.Printf(data)
		} else {
			color.Red(data)
		}

	}
}
