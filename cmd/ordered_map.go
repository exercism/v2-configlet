package cmd

import (
	"bytes"
	"encoding/json"
	"sort"
)

// generate json from a map given a particular ordering
// https://stackoverflow.com/a/30673838/504550

// KeyVal is the key-value component of an OrderedMap
type KeyVal struct {
	Key string
	Val interface{}
}

// OrderedMap is a newtype which JSON-encodes an ordered map
type OrderedMap []KeyVal

// MarshalJSON implements the json.Marshaler interface
func (omap OrderedMap) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("{")
	for i, kv := range omap {
		if i != 0 {
			buf.WriteString(",")
		}
		// marshal key
		key, err := json.Marshal(kv.Key)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteString(":")
		// marshal value
		val, err := json.Marshal(kv.Val)
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}

	buf.WriteString("}")
	return buf.Bytes(), nil
}

// WithOrdering transforms an arbitrary string-keyed map into an OrderedMap.
//
// All keys contained in `order` are emitted in that order. If a given key is
// not present in `dict`, it is omitted without error.
//
// After all ordered keys are emitted, all remaining keys are
// emitted alphabetically.
func WithOrdering(dict map[string]interface{}, order ...string) OrderedMap {
	var om OrderedMap
	for _, key := range order {
		val, isPresent := dict[key]
		if isPresent {
			om = append(om, KeyVal{Key: key, Val: val})
		}
		delete(dict, key)
	}

	var unorderedKeys []string
	for key := range dict {
		unorderedKeys = append(unorderedKeys, key)
	}
	sort.Strings(unorderedKeys)
	for _, key := range unorderedKeys {
		val := dict[key]
		om = append(om, KeyVal{Key: key, Val: val})
	}

	return om
}
