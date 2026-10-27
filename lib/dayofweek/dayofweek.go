package dayofweek

type DaysOfWeek uint8

const (
	Sunday DaysOfWeek = 1 << iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

var (
	Weekdays DaysOfWeek
	Weekends DaysOfWeek
	Daily    DaysOfWeek
)

var stringsToDays map[string]DaysOfWeek

func init() {
	stringsToDays = map[string]DaysOfWeek{
		"Sunday":    Sunday,
		"Monday":    Monday,
		"Tuesday":   Tuesday,
		"Wednesday": Wednesday,
		"Thursday":  Thursday,
		"Friday":    Friday,
		"Saturday":  Saturday,
	}
	Weekdays = Monday | Tuesday | Wednesday | Thursday | Friday
	Weekends = Saturday | Sunday
	Daily = Weekdays | Weekends
}

func DaysOfWeekFromStringArray(week []string) DaysOfWeek {
	var result = DaysOfWeek(0)
	for _, day := range week {
		u, ok := stringsToDays[day]
		if ok {
			result = result | u
		}
	}
	return result
}
