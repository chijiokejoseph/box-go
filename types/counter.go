package types

import (
	"errors"
	"strconv"
	"strings"
)

type Counter struct {
	value  int
	filled bool
	boxed  bool
	base   []int
}

var ErrEmptyBase = errors.New("slice of integers should not be empty")

func divides(num int, base []int) bool {
	for _, b := range base {
		if num%b == 0 {
			return true
		}
	}
	return false
}

func BuildCounterFromBase(base []int) (Counter, error) {
	if len(base) <= 0 {
		return Counter{}, ErrEmptyBase
	}
	return Counter{value: 1, filled: true, boxed: divides(1, base), base: base}, nil
}

func BuildCounterFromInt(value int, base []int) (Counter, error) {
	if len(base) <= 0 {
		return Counter{}, ErrEmptyBase
	}
	return Counter{value: value, filled: true, boxed: divides(value, base), base: base}, nil
}

func BuildCounterFromStr(str string, base []int) (Counter, error) {
	if strings.ToLower(str) == "box" {
		return Counter{filled: false, boxed: true, base: base}, nil
	}
	value, err := strconv.Atoi(str)
	if err != nil {
		return Counter{}, err
	}
	return BuildCounterFromInt(value, base)
}

func (c *Counter) Get() int {
	return c.value
}

func (c *Counter) Filled() bool {
	return c.filled
}

func (c *Counter) Boxed() bool {
	return c.boxed
}

func (c *Counter) Add(other *Counter) Counter {
	var (
		boxed  bool
		filled bool
		value  int
	)
	base := c.base
	switch {
	case c.filled && !other.filled:
		{
			value, boxed, filled = c.value, c.boxed, true
		}

	case !c.filled && other.filled:
		{
			value, boxed, filled = other.value, other.boxed, true
		}
	case c.filled && other.filled:
		{
			value, boxed, filled = c.value+other.value, c.boxed && other.boxed, true
		}
	default:
		{
			value, boxed, filled = 0, c.boxed && other.boxed, false
		}
	}
	return Counter{value, filled, boxed, base}
}

func (c *Counter) Sub(other *Counter) Counter {
	var (
		boxed  bool
		filled bool
		value  int
	)
	base := c.base
	switch {
	case c.filled && !other.filled:
		{
			value, boxed, filled = c.value, c.boxed, true
		}

	case !c.filled && other.filled:
		{
			value, boxed, filled = other.value, other.boxed, true
		}
	case c.filled && other.filled:
		{
			value, boxed, filled = c.value-other.value, c.boxed && other.boxed, true
		}
	default:
		{
			value, boxed, filled = 0, c.boxed && other.boxed, false
		}
	}
	if value < 0 {
		value = 0
	}
	return Counter{value, filled, boxed, base}
}
