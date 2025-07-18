package types

import "log"

type Sequence struct {
	base     []int
	contents []Counter
}

func isAfterMutliple(num int, base []int) bool {
	divisor := num - 1
	return divides(divisor, base)
}

func (s *Sequence) Add(value Counter) bool {
	last := s.contents[len(s.contents)-1]
	pass1 := last.boxed && isAfterMutliple(value.value, s.base)
	pass2 := last.filled && value.filled && value.Sub(&last).value == 1
	if pass1 || pass2 {
		s.contents = append(s.contents, value)
		return true
	} else {
		return false
	}
}

func BuildSequenceFromBase(base []int) (Sequence, error) {
	if len(base) <= 0 {
		return Sequence{}, ErrEmptyBase
	}
	return Sequence{base, make([]Counter, 0)}, nil
}

func NewGenerator(base []int) (func() Sequence, func() Counter, error) {
	if len(base) <= 0 {
		return func() Sequence {
				return Sequence{}
			}, func() Counter {
				return Counter{}
			}, ErrEmptyBase
	}
	genSequence := func() Sequence {
		value, err := BuildSequenceFromBase(base)
		if err != nil {
			log.Fatalln(err)
		}
		return value
	}
	genCounter := func() Counter {
		value, err := BuildCounterFromBase(base)
		if err != nil {
			log.Fatalln(err)
		}
		return value
	}
	return genSequence, genCounter, nil
}
