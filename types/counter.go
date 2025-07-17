package types

import (
	"fmt"
	"strings"
)

type Counter struct {
	value  int
	filled bool
	boxed  bool
}

func divides(num int, base []int) bool {
	for _, b := range base {
		if num%b == 0 {
			return true
		}
	}
	return false
}

func NewCounter(base []int) Counter {
	return Counter{value: 1, filled: true, boxed: divides(1, base)}
}

func NewCounterFromInt(value int, base []int) Counter {
	return Counter{value: value, filled: true, boxed: divides(value, base)}
}

func NewCounterFromStr(str string, base []int) (Counter, error) {
	var value int
	if strings.ToLower(str) == "box" {
		return Counter{filled: false, boxed: true}, nil
	}
	_, err := fmt.Sscanf(str, "%d", &value)
	if err != nil {
		return Counter{}, err
	}
	return NewCounterFromInt(value, base), nil
}

func (c Counter) Get() int {
	return c.value
}

func (c Counter) Filled() bool {
	return c.filled
}

func (c Counter) Boxed() bool {
	return c.boxed
}
