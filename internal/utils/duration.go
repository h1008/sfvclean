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

func ParseDuration(str string) (Duration, error) {
	l := len(str) - 1
	if l <= 0 {
		return Duration{}, ErrInvalidDuration
	}

	count, err := strconv.Atoi(str[:l])
	if err != nil {
		return Duration{}, err
	}

	if count < 0 {
		return Duration{}, ErrInvalidDuration
	}

	switch str[l] {
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
