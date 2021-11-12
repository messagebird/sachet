package sms77api

import (
	"reflect"
	"strconv"
)

func pickMapByKey(needle interface{}, haystack interface{}) (interface{}, interface{}) {
	mapIter := reflect.ValueOf(haystack).MapRange()

	for mapIter.Next() {
		if needle == mapIter.Key() {
			return needle, mapIter.Value()
		}
	}

	return nil, nil
}

func inArray(needle interface{}, haystack interface{}) bool {
	slice := reflect.ValueOf(haystack)
	c := slice.Len()

	for i := 0; i < c; i++ {
		if needle == slice.Index(i).Interface() {
			return true
		}
	}

	return false
}

func toUint(id string, bitSize int) uint64 {
	n, err := strconv.ParseUint(id, 10, bitSize)

	if nil == err {
		return n
	}

	return 0
}
