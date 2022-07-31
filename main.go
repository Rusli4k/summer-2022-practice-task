package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type Trains []Train

type Train struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float32
	ArrivalTime        time.Time
	DepartureTime      time.Time
}

func checkForErrors(err error) {
	if err != nil {
		fmt.Println("error is detected:", err)
	}
}

func main() { //func getTrains
	jsonfile, err := os.Open("data.json")
	defer jsonfile.Close()
	checkForErrors(err)
	data, err := io.ReadAll(jsonfile)
	checkForErrors(err)
	var trains Trains
	err = json.Unmarshal(data, &trains)
	checkForErrors(err)
	for i, v := range trains {
		fmt.Println(i, v)
		if i == 5 {
			break
		}
	}
	fmt.Println("hello")
}

func (t *Train) UnmarshalJSON(data []byte) error {
	type Alias Train
	aux := &struct {
		ArrivalTime   string
		DepartureTime string
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	t.ArrivalTime, err = time.Parse("15:04:05", aux.ArrivalTime)
	if err != nil {
		return err
	}

	t.DepartureTime, err = time.Parse("15:04:05", aux.DepartureTime)
	if err != nil {
		return err
	}

	return nil
}
