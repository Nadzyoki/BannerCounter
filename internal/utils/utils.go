package utils

import (
	"errors"
	"strings"
	"time"
)

const (
	IncorrectFormatKey  = "incorrect format key"
	IncorrectFormatDate = "incorrect format date"
)

func SplitKey(key string) (string, time.Time, error) {
	parts := strings.Split(key, "/")
	if len(parts) == 2 {
		tm, err := time.Parse("2006-01-02 15:04", parts[1])
		if err != nil {
			return "", time.Time{}, errors.New(IncorrectFormatDate)
		}
		return parts[0], tm, nil
	}
	return "", time.Time{}, errors.New(IncorrectFormatKey)
}

func CreateKey(id string, tm time.Time) string {
	return id + "/" + tm.Format("2006-01-02 15:04")
}
