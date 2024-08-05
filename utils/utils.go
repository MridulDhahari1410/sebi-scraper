package utils

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	"sebi-scrapper/constants"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/sinhashubham95/go-utils/errors"
	"github.com/sinhashubham95/go-utils/log"
)

// GetHandlerName will return the name of the api handler being called.
func GetHandlerName(ctx *gin.Context) string {
	handlers := ctx.HandlerNames()
	l := len(handlers)
	if l == 0 {
		return ctx.Request.RequestURI
	}
	return handlers[l-1]
}

// Copy is a general copy function to copy all k-v into destination map; use only for small maps.
func Copy(dest, src any) {
	dv, sv := reflect.ValueOf(dest), reflect.ValueOf(src)

	for _, k := range sv.MapKeys() {
		dv.SetMapIndex(k, sv.MapIndex(k))
	}
}

// CloseData is used to close data.
func CloseData(data io.ReadCloser) {
	err := data.Close()
	if err != nil {
		log.Error(context.Background()).Err(err).Msg("error closing data")
	}
}

func GetDataAsBytes(data io.ReadCloser) ([]byte, error) {
	by, err := io.ReadAll(data)
	defer CloseData(data)
	return by, err
}

// GetDataAsString is used to get the data as string.
func GetDataAsString(data io.ReadCloser) (string, error) {
	by, err := GetDataAsBytes(data)
	if err != nil {
		return "", err
	}
	return string(by), nil
}

// GetJSONData is used to get the JSON data parsed into a struct, make sure you pass the struct by reference.
func GetJSONData(data io.ReadCloser, val any) error {
	by, err := GetDataAsBytes(data)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(by, val)
}

// MarshalJSON - Marshal returns the JSON encoding of data.
func MarshalJSON(data any) ([]byte, error) {
	return jsoniter.Marshal(data)
}

// UnmarshalJSON -- returns the parsed JSON-encoded data in `val`.
func UnmarshalJSON(data []byte, val any) error {
	return jsoniter.Unmarshal(data, val)
}

// GetRunAtWithOffset is used to get random run at with offset.
func GetRunAtWithOffset(ctx context.Context, at string, offsetInMinutes int) string {
	parts := strings.Split(at, ":")
	if len(parts) != 2 {
		return at
	}
	hr, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Error(ctx).Err(err).Msg("cannot convert hr string to int")
		return at
	}
	mi, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Error(ctx).Err(err).Msg("cannot convert min string to int")
		return at
	}
	num := getRandomNumber(-offsetInMinutes, offsetInMinutes) + mi
	mi = num % 60
	hr = (hr + (num / 60)) % 24
	if hr < 0 {
		hr = hr + 24
	}
	if mi < 0 {
		mi = mi + 60
		if hr == 0 {
			hr = 23
		} else {
			hr--
		}
	}
	fmt.Println(hr, mi)
	return fmt.Sprintf("%02d:%2d", hr, mi)
}

// GetErrorDetails is used to get error details as string.
func GetErrorDetails(e error) string {
	var er errors.Error
	if ok := errors.As(e, &er); ok {
		if s, ok := er.Details.(string); ok {
			return strings.ToLower(s)
		}
	}
	return constants.Empty
}

func getRandomNumber(mi int, ma int) int {
	return rand.Intn(ma-mi+1) + mi
}

func ConvertStringToSliceOfStrings(s string) []string {
	result := strings.Split(s, ",")
	return result
}

