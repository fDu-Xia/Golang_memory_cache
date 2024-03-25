package cache

import (
	"unsafe"
)

func GetValSize(v Value) int64 {
	return int64(unsafe.Sizeof(v))
}
