# Shift Calculator

This program processes a dataset of employee work shifts, calculates regular and overtime hours, detects invalid overlapping shifts and skips counting the hours, handles CST/CDT timezone changeovers and correctly handles shifts that cross week boundaries.

## Installation (MacOS)

1. Prerequisites

Make sure you have Go (I am targeting 1.24) installed. You can check your version by running:

```BASH
go version
```

If you don't have Go installed, install it via Homebrew:

```BASH
brew install go
```

2. Clone the Repository

```BASH
git clone <repository-url>
cd shift_calc
```

3. Install Dependencies

```BASH
go mod tidy
```

## Usage

Running the Program

```BASH
go run main.go
```

This will process the shifts from data/dataset.json and print the output in JSON format.

Example Output:

```JSON
[
  {
        "EmployeeID": 41779236,
        "StartOfWeek": "2021-08-22",
        "RegularHours": 40,
        "OvertimeHours": 10,
        "InvalidShifts": []
    },
    {
        "EmployeeID": 41779236,
        "StartOfWeek": "2021-08-29",
        "RegularHours": 12.5,
        "OvertimeHours": 0,
        "InvalidShifts": []
    },
    {
        "EmployeeID": 41779237,
        "StartOfWeek": "2021-08-29",
        "RegularHours": 25,
        "OvertimeHours": 0,
        "InvalidShifts": []
    },
    {
        "EmployeeID": 41779237,
        "StartOfWeek": "2021-08-22",
        "RegularHours": 37.5,
        "OvertimeHours": 0,
        "InvalidShifts": []
    },
    {
        "EmployeeID": 41784027,
        "StartOfWeek": "2021-08-29",
        "RegularHours": 12,
        "OvertimeHours": 0,
        "InvalidShifts": []
    },
    {
        "EmployeeID": 41784027,
        "StartOfWeek": "2021-08-22",
        "RegularHours": 31.5,
        "OvertimeHours": 0,
        "InvalidShifts": []
    }
]
```

## Running Tests

The project includes tests for shift processing and time utilities.

Run all tests with:

```
go test ./...
```

## How it Works

Reads data/dataset.json: Loads and validates shift data.

Loops through shifts:

-Calculates Regular vs. Overtime Hours:

-Detects Overlapping Shifts: Flags shifts that conflict for the same employee

-Detects Shifts that span a week boundary and splits the hours between the weeks

Outputs the JSON

## Future Improvements

If given more time, I would:

1. Remove repetitive code!
2. Split up code for more testability
3. Add ALOT more unit tests to cover all edge cases, I ran out of time so my tests are just quick and dirty
4. Optimize overlapping shift detection (just a loop of all shifts, fine for this dataset, this should have been a map of shift data that I could search by EmployeeID to make it much quicker)
5. Better app structure, its been awhile since I wrote Go so I am not sure of all the conventions
6. Clean up my comments, they are all over the place to give you an idea of what I was thinking.
7. There are probably a ton of spelling mistakes in there
8. Probably looked for more help from libraries, I don't really import anything other than `validator\v10`
9. Used JSONSchema to validate incoming and outgoing JSON against the specifications.
10. Handled CST/CDT converstion in a custom Unmarshal func for shift

## Author

Bryon Wicklund bryon.wicklund@gmail.com
