package restaurant

import "fmt"

type DayOfWeek int

const (
	Sunday    DayOfWeek = iota // 0
	Monday                     // 1
	Tuesday                    // 2
	Wednesday                  // 3
	Thursday                   // 4
	Friday                     // 5
	Saturday                   // 6
)

func (d DayOfWeek) String() string {
	switch d {
	case Sunday:
		return "Sunday"
	case Monday:
		return "Monday"
	case Tuesday:
		return "Tuesday"
	case Wednesday:
		return "Wednesday"
	case Thursday:
		return "Thursday"
	case Friday:
		return "Friday"
	case Saturday:
		return "Saturday"
	default:
		return "Unknown"
	}
}

func ParseDayOfWeek(i int32) (DayOfWeek, error) {
	if i < 0 || i > 6 {
		return Sunday, fmt.Errorf("invalid day_of_week: %d", i)
	}
	return DayOfWeek(i), nil
}
