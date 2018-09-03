package txt_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/txt"
	"github.com/stretchr/testify/assert"
)

func TestFirstN(t *testing.T) {
	table := []struct {
		In  string
		N   int
		Out string
	}{
		{In: "abcd", N: 3, Out: "abc"},
		{In: "abcd", N: 5, Out: "abcd"},
		{In: "abcd", N: 0, Out: ""},
		{In: "abcd", N: -1, Out: ""},
		{In: "aécd", N: 3, Out: "aéc"},
		{In: "aécd", N: 5, Out: "aécd"},
		{In: "aécd", N: 0, Out: ""},
		{In: "aécd", N: -1, Out: ""},
	}
	for i, one := range table {
		assert.Equal(t, one.Out, txt.FirstN(one.In, one.N), "#%d", i)
	}
}

func TestLastN(t *testing.T) {
	table := []struct {
		In  string
		N   int
		Out string
	}{
		{In: "abcd", N: 3, Out: "bcd"},
		{In: "abcd", N: 5, Out: "abcd"},
		{In: "abcd", N: 0, Out: ""},
		{In: "abcd", N: -1, Out: ""},
		{In: "aécd", N: 3, Out: "écd"},
		{In: "aécd", N: 5, Out: "aécd"},
		{In: "aécd", N: 0, Out: ""},
		{In: "aécd", N: -1, Out: ""},
	}
	for i, one := range table {
		assert.Equal(t, one.Out, txt.LastN(one.In, one.N), "#%d", i)
	}
}
