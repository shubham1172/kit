/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dapr/kit/config"
	"github.com/dapr/kit/ptr"
)

type testConfig struct { //nolint: maligned
	Int         int            `mapstructure:"int"`
	IntPtr      *int           `mapstructure:"intPtr"`
	Int64       int64          `mapstructure:"int64"`
	Int64Ptr    *int64         `mapstructure:"int64Ptr"`
	Int32       int32          `mapstructure:"int32"`
	Int32Ptr    *int32         `mapstructure:"int32Ptr"`
	Int16       int16          `mapstructure:"int16"`
	Int16Ptr    *int16         `mapstructure:"int16Ptr"`
	Int8        int8           `mapstructure:"int8"`
	Int8Ptr     *int8          `mapstructure:"int8Ptr"`
	Uint        uint           `mapstructure:"uint"`
	UintPtr     *uint          `mapstructure:"uintPtr"`
	Uint64      uint64         `mapstructure:"uint64"`
	Uint64Ptr   *uint64        `mapstructure:"uint64Ptr"`
	Uint32      uint32         `mapstructure:"uint32"`
	Uint32Ptr   *uint32        `mapstructure:"uint32Ptr"`
	Uint16      uint16         `mapstructure:"uint16"`
	Uint16Ptr   *uint16        `mapstructure:"uint16Ptr"`
	Byte        byte           `mapstructure:"byte"`
	BytePtr     *byte          `mapstructure:"bytePtr"`
	Float64     float64        `mapstructure:"float64"`
	Float64Ptr  *float64       `mapstructure:"float64Ptr"`
	Float32     float32        `mapstructure:"float32"`
	Float32Ptr  *float32       `mapstructure:"float32Ptr"`
	Bool        bool           `mapstructure:"bool"`
	BoolPtr     *bool          `mapstructure:"boolPtr"`
	Duration    time.Duration  `mapstructure:"duration"`
	DurationPtr *time.Duration `mapstructure:"durationPtr"`
	Time        time.Time      `mapstructure:"time"`
	TimePtr     *time.Time     `mapstructure:"timePtr"`
	String      string         `mapstructure:"string"`
	StringPtr   *string        `mapstructure:"stringPtr"`
	Decoded     Decoded        `mapstructure:"decoded"`
	DecodedPtr  *Decoded       `mapstructure:"decodedPtr"`
	Nested      nested         `mapstructure:"nested"`
	NestedPtr   *nested        `mapstructure:"nestedPtr"`
}

type nested struct {
	Integer int64  `mapstructure:"integer"`
	String  string `mapstructure:"string"`
}

type Decoded int

func (u *Decoded) DecodeString(text string) error {
	if text == "unlimited" {
		*u = -1

		return nil
	}

	val, err := strconv.Atoi(text)
	if err != nil {
		return err
	}

	*u = Decoded(val)

	return nil
}

func TestDecode(t *testing.T) {
	timeVal := getTimeVal()
	tests := map[string]interface{}{
		"primitive values": map[string]interface{}{
			"int":         -9999,
			"intPtr":      ptr.Of(-9999),
			"int64":       -1234,
			"int64Ptr":    ptr.Of(-12345),
			"int32":       -5678,
			"int32Ptr":    ptr.Of(-5678),
			"int16":       -9012,
			"int16Ptr":    ptr.Of(-9012),
			"int8":        -128,
			"int8Ptr":     ptr.Of(-128),
			"uint":        9999,
			"uintPtr":     ptr.Of(9999),
			"uint64":      1234,
			"uint64Ptr":   ptr.Of(1234),
			"uint32":      5678,
			"uint32Ptr":   ptr.Of(5678),
			"uint16":      9012,
			"uint16Ptr":   ptr.Of(9012),
			"byte":        255,
			"bytePtr":     ptr.Of(255),
			"float64":     1234.5,
			"float64Ptr":  ptr.Of(1234.5),
			"float32":     6789.5,
			"float32Ptr":  ptr.Of(6789.5),
			"bool":        true,
			"boolPtr":     ptr.Of(true),
			"duration":    5 * time.Second,
			"durationPtr": ptr.Of(5 * time.Second),
			"time":        timeVal,
			"timePtr":     ptr.Of(timeVal),
			"string":      1234,
			"stringPtr":   ptr.Of("1234"),
			"decoded":     "unlimited",
			"decodedPtr":  "unlimited",
			"nested": map[string]interface{}{
				"integer": 1234,
				"string":  5678,
			},
			"nestedPtr": map[string]interface{}{
				"integer": 1234,
				"string":  5678,
			},
		},
		"string values": map[string]interface{}{
			"int":         "-9999",
			"intPtr":      "-9999",
			"int64":       "-1234",
			"int64Ptr":    "-12345",
			"int32":       "-5678",
			"int32Ptr":    "-5678",
			"int16":       "-9012",
			"int16Ptr":    "-9012",
			"int8":        "-128",
			"int8Ptr":     "-128",
			"uint":        "9999",
			"uintPtr":     "9999",
			"uint64":      "1234",
			"uint64Ptr":   "1234",
			"uint32":      "5678",
			"uint32Ptr":   "5678",
			"uint16":      "9012",
			"uint16Ptr":   "9012",
			"byte":        "255",
			"bytePtr":     "255",
			"float64":     "1234.5",
			"float64Ptr":  "1234.5",
			"float32":     "6789.5",
			"float32Ptr":  "6789.5",
			"bool":        "true",
			"boolPtr":     "true",
			"duration":    "5000",
			"durationPtr": "5s",
			"time":        "2021-01-02T15:04:05-07:00",
			"timePtr":     "2021-01-02T15:04:05-07:00",
			"string":      "1234",
			"stringPtr":   "1234",
			"decoded":     "unlimited",
			"decodedPtr":  "unlimited",
			"nested": map[string]string{
				"integer": "1234",
				"string":  "5678",
			},
			"nestedPtr": map[string]string{
				"integer": "1234",
				"string":  "5678",
			},
		},
	}

	expected := getExpected()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var actual testConfig
			err := config.Decode(tc, &actual)
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestDecodeErrors(t *testing.T) {
	var actual testConfig
	err := config.Decode(map[string]interface{}{
		"int":         "-badval",
		"intPtr":      "-badval",
		"int64":       "-badval",
		"int64Ptr":    "-badval",
		"int32":       "-badval",
		"int32Ptr":    "-badval",
		"int16":       "-badval",
		"int16Ptr":    "-badval",
		"int8":        "-badval",
		"int8Ptr":     "-badval",
		"uint":        "-9999",
		"uintPtr":     "-9999",
		"uint64":      "-1234",
		"uint64Ptr":   "-1234",
		"uint32":      "-5678",
		"uint32Ptr":   "-5678",
		"uint16":      "-9012",
		"uint16Ptr":   "-9012",
		"byte":        "-1",
		"bytePtr":     "-1",
		"float64":     "badval.5",
		"float64Ptr":  "badval.5",
		"float32":     "badval.5",
		"float32Ptr":  "badval.5",
		"bool":        "badval",
		"boolPtr":     "badval",
		"duration":    "badval",
		"durationPtr": "badval",
		"time":        "badval",
		"timePtr":     "badval",
		"decoded":     "badval",
		"decodedPtr":  "badval",
		"string":      1234,
		"stringPtr":   1234,
	}, &actual)
	if assert.Error(t, err) {
		errMsg := err.Error()
		expectedNumErrors := 32
		expectedPrefix := " error(s) decoding:"
		assert.True(t, strings.HasPrefix(errMsg, fmt.Sprintf("%d%s", expectedNumErrors, expectedPrefix)), errMsg)
		prefixIndex := strings.Index(errMsg, expectedPrefix)
		if assert.True(t, prefixIndex != -1) {
			errMsg = errMsg[prefixIndex+len(expectedPrefix):]
			errMsg = strings.TrimSpace(errMsg)
			errors := strings.Split(errMsg, "\n")
			errorSet := make(map[string]struct{}, len(errors))
			for _, e := range errors {
				errorSet[e] = struct{}{}
			}
			expectedErrors := []string{
				"* error decoding 'int': invalid int \"-badval\"",
				"* error decoding 'intPtr': invalid int \"-badval\"",
				"* error decoding 'int16': invalid int16 \"-badval\"",
				"* error decoding 'int16Ptr': invalid int16 \"-badval\"",
				"* error decoding 'int32': invalid int32 \"-badval\"",
				"* error decoding 'int32Ptr': invalid int32 \"-badval\"",
				"* error decoding 'int64': invalid int64 \"-badval\"",
				"* error decoding 'int64Ptr': invalid int64 \"-badval\"",
				"* error decoding 'int8': invalid int8 \"-badval\"",
				"* error decoding 'int8Ptr': invalid int8 \"-badval\"",
				"* error decoding 'uint': invalid uint \"-9999\"",
				"* error decoding 'uintPtr': invalid uint \"-9999\"",
				"* error decoding 'uint64': invalid uint64 \"-1234\"",
				"* error decoding 'uint64Ptr': invalid uint64 \"-1234\"",
				"* error decoding 'uint32': invalid uint32 \"-5678\"",
				"* error decoding 'uint32Ptr': invalid uint32 \"-5678\"",
				"* error decoding 'uint16': invalid uint16 \"-9012\"",
				"* error decoding 'uint16Ptr': invalid uint16 \"-9012\"",
				"* error decoding 'byte': invalid uint8 \"-1\"",
				"* error decoding 'bytePtr': invalid uint8 \"-1\"",
				"* error decoding 'float32': invalid float32 \"badval.5\"",
				"* error decoding 'float32Ptr': invalid float32 \"badval.5\"",
				"* error decoding 'float64': invalid float64 \"badval.5\"",
				"* error decoding 'float64Ptr': invalid float64 \"badval.5\"",
				"* error decoding 'duration': invalid duration \"badval\"",
				"* error decoding 'durationPtr': invalid duration \"badval\"",
				"* error decoding 'time': invalid time \"badval\"",
				"* error decoding 'timePtr': invalid time \"badval\"",
				"* error decoding 'decoded': invalid Decoded \"badval\": strconv.Atoi: parsing \"badval\": invalid syntax",
				"* error decoding 'decodedPtr': invalid Decoded \"badval\": strconv.Atoi: parsing \"badval\": invalid syntax",
				"* error decoding 'bool': invalid bool \"badval\"",
				"* error decoding 'boolPtr': invalid bool \"badval\"",
			}
			for _, expectedError := range expectedErrors {
				assert.Contains(t, errors, expectedError)
				delete(errorSet, expectedError)
			}
			assert.Empty(t, errorSet)
		}
	}
}

func getTimeVal() time.Time {
	timeVal, _ := time.Parse(time.RFC3339, "2021-01-02T15:04:05-07:00")

	return timeVal
}

func getExpected() testConfig {
	timeVal := getTimeVal()

	return testConfig{
		Int:         -9999,
		IntPtr:      ptr.Of(-9999),
		Int64:       -1234,
		Int64Ptr:    ptr.Of(int64(-12345)),
		Int32:       -5678,
		Int32Ptr:    ptr.Of(int32(-5678)),
		Int16:       -9012,
		Int16Ptr:    ptr.Of(int16(-9012)),
		Int8:        -128,
		Int8Ptr:     ptr.Of(int8(-128)),
		Uint:        9999,
		UintPtr:     ptr.Of(uint(9999)),
		Uint64:      1234,
		Uint64Ptr:   ptr.Of(uint64(1234)),
		Uint32:      5678,
		Uint32Ptr:   ptr.Of(uint32(5678)),
		Uint16:      9012,
		Uint16Ptr:   ptr.Of(uint16(9012)),
		Byte:        255,
		BytePtr:     ptr.Of(byte(255)),
		Float64:     1234.5,
		Float64Ptr:  ptr.Of(1234.5),
		Float32:     6789.5,
		Float32Ptr:  ptr.Of(float32(6789.5)),
		Bool:        true,
		BoolPtr:     ptr.Of(true),
		Duration:    5 * time.Second,
		DurationPtr: ptr.Of(5 * time.Second),
		Time:        timeVal,
		TimePtr:     ptr.Of(timeVal),
		String:      "1234",
		StringPtr:   ptr.Of("1234"),
		Decoded:     -1,
		DecodedPtr:  ptr.Of(Decoded(-1)),
		Nested: nested{
			Integer: 1234,
			String:  "5678",
		},
		NestedPtr: &nested{
			Integer: 1234,
			String:  "5678",
		},
	}
}
