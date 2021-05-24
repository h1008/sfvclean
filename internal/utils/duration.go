package utils

import (
	"errors"
	"strconv"
)

type Duration struct {
	Years  int
	Months int
	Days   int
}

var ErrInvalidDuration = errors.New("invalid duration")

func ParseDuration(s string) (Duration, error) {
	l := len(s) - 1
	if l <= 0 {
		return Duration{}, ErrInvalidDuration
	}

	count, err := strconv.Atoi(s[:l])
	if err != nil {
		return Duration{}, err
	}

	if count < 0 {
		return Duration{}, ErrInvalidDuration
	}

	switch s[l] {
	case 'Y':
		return Duration{Years: count}, nil
	case 'M':
		return Duration{Months: count}, nil
	case 'D':
		return Duration{Days: count}, nil
	default:
		return Duration{}, ErrInvalidDuration
	}
}
