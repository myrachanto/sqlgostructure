package support

import (
	"strconv"
	"time"
)

func Stamper() string {
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[0:5]
	return special
}