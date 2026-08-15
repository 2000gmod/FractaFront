package fir

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"reflect"
)

func hashSlice(types []Type) string {
	h := fnv.New64a()

	var lenBuf [8]byte
	binary.LittleEndian.PutUint64(lenBuf[:], uint64(len(types)))
	h.Write(lenBuf[:])

	for _, f := range types {
		ptr := reflect.ValueOf(f).Pointer()
		var ptrBuf [8]byte
		binary.LittleEndian.PutUint64(ptrBuf[:], uint64(ptr))
		h.Write(ptrBuf[:])
	}
	return fmt.Sprintf("%016x", h.Sum64())
}
