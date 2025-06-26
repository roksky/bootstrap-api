package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"strings"
	"unicode"
)

func ConvertSliceToReference[T interface{}](slice []T) []*T {
	result := make([]*T, len(slice))
	for i, value := range slice {
		result[i] = &value
	}
	return result
}

func ReadJsonAsType[T interface{}](rc io.ReadCloser) ([]T, error) {
	var items []T
	jsonData, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func ParamAsUUId(ctx *gin.Context, key string) (uuid.UUID, error) {
	keyVal := ctx.Param(key)
	if keyVal == "" {
		return uuid.Nil, errors.New(fmt.Sprintf("%s is required", keyVal))
	}
	// parse as uuid
	uuidVal, err := uuid.Parse(keyVal)
	if err != nil {
		return uuid.Nil, errors.New(fmt.Sprintf("%s is not a valid uuid", keyVal))
	}

	return uuidVal, nil
}

func QueryAsUUId(ctx *gin.Context, key string) *uuid.UUID {
	keyVal := ctx.Query(key)
	if keyVal == "" {
		return nil
	}
	// parse as uuid
	uuidVal, err := uuid.Parse(keyVal)
	if err != nil {
		return nil
	}

	return &uuidVal
}

func GetSortString(input string) string {
	parts := strings.Split(input, ",")
	var transformedParts []string

	for _, part := range parts {
		subParts := strings.Split(part, ":")
		if len(subParts) > 1 {
			transformedParts = append(transformedParts, fmt.Sprintf("%s %s", CamelToSnake(subParts[0]), subParts[1]))
		} else {
			transformedParts = append(transformedParts, CamelToSnake(subParts[0]))
		}
	}
	return strings.TrimSpace(strings.Join(transformedParts, ", "))
}

// CamelToSnake converts a camel case string to snake case.
func CamelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
