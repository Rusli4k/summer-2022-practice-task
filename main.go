package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
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

	// func main() {
	//     // ... запит даних від користувача
	//     result, err := FindTrains(departureStation, arrivalStation, criteria))
	//     // ... обробка помилки
	//     // ... друк result
	// }
	var (
		departureStation string
		arrivalStation   string
		criteria         string
	)

	// fmt.Println("Please, input departure station: (press ENTER after input)")
	// fmt.Scanln(&departureStation)

	// fmt.Println("and where are you going: (press ENTER after input)")
	// fmt.Scanln(&arrivalStation)

	// fmt.Println("result may be sorts by price/arrival-time/departure-time: (press ENTER after input)")
	// fmt.Scanln(&criteria)
	departureStation = "1929"
	arrivalStation = "1921"
	criteria = "arrival-time"

	//result, err := FindTrains(departureStation, arrivalStation, criteria)
	result, err := FindTrains(departureStation, arrivalStation, criteria)
	checkForErrors(err)
	for _, v := range result {
		fmt.Printf("%+v\n", v)
	}
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	var ( //IMHO we can use variable of the struct Train, but one repeat is better than one dependency
		departureStationID int
		arrivalStationID   int
	)
	var ts Trains
	var ans Trains

	const (
		criteria1 = "price"
		criteria2 = "arrival-time"
		criteria3 = "departure-time"
	)
	ts.UnmarshalTrains("data.json") //reading data about trains from file

	//checks for invalid input of departureStation
	if departureStation == "" {
		return nil, errors.New("empty departure station")
	}
	res, err := strconv.ParseUint(departureStation, 10, 0)
	if err != nil {
		return nil, errors.New("bad departure station input")
	}
	departureStationID = int(res) //type convert without error check because int>uint

	//checks for invalid input of arrivalStation
	if arrivalStation == "" {
		return nil, errors.New("empty arrival station")
	}
	res, err = strconv.ParseUint(arrivalStation, 10, 0)
	if err != nil {
		return nil, errors.New("bad arrival station input")
	}
	arrivalStationID = int(res) //type convert without error check because int>uint

	//check for invalid input of criteria
	if criteria != criteria1 && criteria != criteria2 && criteria != criteria3 {
		return nil, errors.New("unsupported criteria")
	}

	for _, v := range ts {
		if v.ArrivalStationID == arrivalStationID && v.DepartureStationID == departureStationID {
			ans = append(ans, v)
		}
	}
	switch criteria {
	case criteria1:
		sort.SliceStable(ans, func(i, j int) bool {
			return ans[i].Price < ans[j].Price
		})
	case criteria2:
		sort.SliceStable(ans, func(i, j int) bool {
			return ans[i].ArrivalTime.Before(ans[j].ArrivalTime)
		})
	case criteria3:
		sort.SliceStable(ans, func(i, j int) bool {
			return ans[i].DepartureTime.Before(ans[j].DepartureTime)
		})
	}

	if len(ans) >= 3 {
		return ans[:3], nil
	}
	return ans, nil
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
