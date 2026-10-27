package dayofweek

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDaysOfWeek(t *testing.T) {
	assert.Equal(t, 127, int(Daily))
	assert.Equal(t, 1+64, int(Weekends))
	assert.Equal(t, 2+4+8+16+32, int(Weekdays))
}
