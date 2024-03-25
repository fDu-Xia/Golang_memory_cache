package utils

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	B = 1 << (iota * 10)
	KB
	MB
	GB
)

func ParseSize(size string) int64 {
	defer func() {
		if p := recover(); p != nil {
			log.Println("ParseSize 仅支持 B、KB、MB、GB")
		}
	}()

	re := regexp.MustCompile(`([0-9]+)(\w+)`)
	match := re.FindStringSubmatch(size)
	sizeNum, _ := strconv.ParseInt(match[1], 10, 64)
	unit := strings.ToUpper(match[2])
	switch unit {
	case "B":
		return sizeNum * B
	case "KB":
		return sizeNum * KB
	case "MB":
		return sizeNum * MB
	case "GB":
		return sizeNum * GB
	default:
		return 100 * MB
	}
}
