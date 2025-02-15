package main

import (
	"encoding/json"
	"fmt"
	"shift_calc/models"
	"shift_calc/services"
)

func main() {
	// Load the shifts from the dataset
	shifts, err := services.LoadShifts("data/dataset.json")
	if err != nil {
		panic(err)
	}

	// Condense the shift data
	var employeeSummary []models.EmployeeReport = services.CondenseShiftData(shifts)

	// marshal and print the employee hours report
	employeeSummaryJSON, err := json.MarshalIndent(employeeSummary, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(employeeSummaryJSON))
}
