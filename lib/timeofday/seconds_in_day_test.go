package timeofday

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHumanReadable(t *testing.T) {
	assert.Equal(t, "10:30 pm", (10*Hour + 30*Minute + PM).HumanReadable())
	assert.Equal(t, "12:00 pm", (12 * Hour).HumanReadable())
	assert.Equal(t, "10:30 am (next day)", (10*Hour + 30*Minute + NEXT_DAY).HumanReadable())
	assert.Equal(t, "10:30 pm (next day)", (10*Hour + 30*Minute + PM + NEXT_DAY).HumanReadable())
}

func TestTimeToSecondsInDay(t *testing.T) {
	example, err := time.Parse(time.DateTime, "2025-06-01 07:55:13")
	assert.NoError(t, err)

	assert.Equal(t, SecondsInDay(25_200+3_300+13), TimeToSecondsInDay(example))
}
