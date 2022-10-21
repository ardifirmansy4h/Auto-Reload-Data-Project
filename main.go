package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Data struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Status struct {
	Status      Data   `json:"status"`
	WaterStatus string `json:"-"`
	WindStatus  string `json:"-"`
}

var dat Status

func RandWater() int {
	maximum := 100
	random := rand.Intn(maximum)
	return random
}

func RandWind() int {
	maximum := 100
	random := rand.Intn(maximum)
	return random
}
func StatusAlert(data Data) (string, string) {
	var water string
	var wind string

	if data.Water <= 5 {
		water = "aman"
	} else if data.Water >= 6 && data.Water <= 8 {
		water = "siaga"
	} else {
		water = "bahaya"
	}

	if data.Wind <= 6 {
		wind = "aman"
	} else if data.Wind >= 7 && data.Wind <= 15 {
		wind = "siaga"
	} else {
		wind = "bahaya"
	}
	return water, wind
}

func start() {
	ticker := time.NewTicker(15 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case w := <-ticker.C:
				fmt.Println("jam ", w)
				do()
			}
		}

	}()

	time.Sleep(60 * time.Minute)
	ticker.Stop()
	done <- true
	fmt.Println("Stopped service")
}

func do() {
	file, err := os.ReadFile("status.json")
	if err != nil {
		fmt.Println("Error reading status.json: ", err)
		return
	}

	err = json.Unmarshal(file, &dat)
	if err != nil {
		fmt.Println("Error status.json: ", err)
		return
	}
	dat.Status.Water = RandWater()
	dat.Status.Wind = RandWind()
	dat.WaterStatus, dat.WindStatus = StatusAlert(dat.Status)

	newFile, err := json.MarshalIndent(dat, "", "  ")
	if err != nil {
		fmt.Println("Error data to json: ", err)
		return
	}

	err = os.WriteFile("status.json", newFile, 0644)
	if err != nil {
		fmt.Println("Error writing data to files status.json: ", err)
		return
	}

	fmt.Println("Water :", dat.Status.Water, "Status", dat.WaterStatus)
	fmt.Println("Wind  :", dat.Status.Wind, "Status", dat.WindStatus)
}

func StatusCtrl(w http.ResponseWriter, r *http.Request) {
	tmplt, err := template.ParseFiles("index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmplt.Execute(w, dat)
}
func main() {
	go start()
	http.HandleFunc("/", StatusCtrl)
	http.ListenAndServe(":1323", nil)
}
