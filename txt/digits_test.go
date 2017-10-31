package txt_test

import (
	"testing"

	"github.com/richardwilkes/gokit/txt"
	"github.com/stretchr/testify/assert"
)

func TestDigitToValue(t *testing.T) {
	checkDigitToValue('5', 5, t)
	checkDigitToValue('٥', 5, t)
	checkDigitToValue('𑁯', 9, t)
	_, err := txt.DigitToValue('a')
	assert.Error(t, err)
}

func checkDigitToValue(ch rune, expected int, t *testing.T) {
	value, err := txt.DigitToValue(ch)
	assert.NoError(t, err)
	assert.EqualValues(t, expected, value)
}
