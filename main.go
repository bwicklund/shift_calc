package main

import (
	"encoding/json"
	"shift_calc/services"
)

func main() {
	shifts, err := services.LoadShifts("data/dataset.json")
	if err != nil {
		panic(err)
	}

	// Lets format shifts as JSON and print it out for now...
	shiftsJSON, err := json.Marshal(shifts)
	if err != nil {
		panic(err)
	}
	println(string(shiftsJSON))
}
