package helper

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/roksky/bootstrap-api/types"
)

func AsPtr[T interface{}](b T) *T {
	return &b
}

func StringToBool(input string) (*bool, error) {
	if input == "true" {
		return AsPtr(true), nil
	} else if input == "false" {
		return AsPtr(false), nil
	} else {
		return nil, errors.New("invalid boolean value")
	}
}

func StringToInt64(input string) (*int64, error) {
	val, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return nil, err
	}
	return AsPtr(val), nil
}

func StringToFloat64(input string) (*float64, error) {
	val, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return nil, err
	}
	return AsPtr(val), nil
}

func StringToDate(input string) (*time.Time, error) {
	val, err := time.Parse(time.RFC3339, input)
	if err != nil {
		return nil, err
	}
	return AsPtr(val), nil
}

func StringToDateOnly(input string) (*types.DateOnly, error) {
	val, err := time.Parse("2006-01-02", input)
	if err != nil {
		return nil, err
	}
	return AsPtr(types.DateOnly{
		Time: val,
	}), nil
}

func StringToTimeOnly(input string) (*types.TimeOnly, error) {
	val, err := time.Parse("15:04", input)
	if err != nil {
		return nil, err
	}
	return AsPtr(types.TimeOnly{
		Time: val,
	}), nil
}

func JsonRawMessageToJson(message json.RawMessage) (map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
