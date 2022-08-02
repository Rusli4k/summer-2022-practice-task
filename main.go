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

func main() { //func getTrains
	var (
		departureStation string
		arrivalStation   string
		criteria         string
	)
	fmt.Println("Please, input departure station:")
	fmt.Scanln(&departureStation)
	fmt.Println("and where are you going:")
	fmt.Scanln(&arrivalStation)
	fmt.Println("result may be sorts by:")
	fmt.Scanln(&criteria)
	var ts Trains
	ts.UnmarshalTrains("data.json")

	for i, v := range ts {
		fmt.Println(i, v)
		if i == 5 {
			break
		}
	}
}

func (t *Trains) UnmarshalTrains(f string) { // func for unmarshaling from file to slice of structs Train.  Take name of file
	jfile, err := os.Open("data.json")
	defer jfile.Close()
	checkForErrors(err)
	data, err := io.ReadAll(jfile)
	checkForErrors(err)
	err = json.Unmarshal(data, &t)
	checkForErrors(err)
}

func (t *Train) UnmarshalJSON(data []byte) error { //addon for customizing Unmarshal for non standart time format
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

func checkForErrors(err error) { //when you repeat some code more then three times, make function  - "Rule of Three"
	if err != nil {
		fmt.Println("error was detected:", err)
	}
}
