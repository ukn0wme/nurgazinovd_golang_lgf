package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidDurationFormat = errors.New("invalid duration format")

type Duration int32

func (r Duration) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d seconds", r)
	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil
}

func (r *Duration) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidDurationFormat
	}
	parts := strings.Split(unquotedJSONValue, " ")
	if len(parts) != 2 || parts[1] != "seconds" {
		return ErrInvalidDurationFormat
	}
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidDurationFormat
	}
	*r = Duration(i)
	return nil
}
