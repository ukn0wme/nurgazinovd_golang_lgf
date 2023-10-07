package data

import (
	"fmt"
	"strconv"
)

type Duration int32

func (r Duration) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d seconds", r)
	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil
}
