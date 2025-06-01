package timeofday

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	time := 10*Hour + 30*Minute + PM

	if time.HumanReadable() != "10:30 pm" {
		fmt.Println(time.HumanReadable())
		t.Fail()
	}

	time = 12*Hour + PM

	if time.HumanReadable() != "12:00 pm" {
		fmt.Println(time.HumanReadable())
		t.Fail()
	}
}

func TestTimeToSecondsInDay(t *testing.T) {
	example, err := time.Parse(time.DateTime, "2025-06-01 07:55:13")
	assert.NoError(t, err)

	assert.Equal(t, TimeToSecondsInDay(example), 25_200+3_300+13)

}
