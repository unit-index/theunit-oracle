package origins

import (
	"encoding/json"
	"strconv"
	"time"
)

type stringAsFloat64 float64

func (s *stringAsFloat64) UnmarshalJSON(bytes []byte) error {
	var ss string
	if err := json.Unmarshal(bytes, &ss); err != nil {
		return err
	}
	f, err := strconv.ParseFloat(ss, 64)
	if err != nil {
		return err
	}
	*s = stringAsFloat64(f)
	return nil
}
func (s *stringAsFloat64) val() float64 {
	return float64(*s)
}

type firstStringFromSliceAsFloat64 float64

func (s *firstStringFromSliceAsFloat64) UnmarshalJSON(bytes []byte) error {
	var ss []string
	if err := json.Unmarshal(bytes, &ss); err != nil {
		return err
	}
	f, err := strconv.ParseFloat(ss[0], 64)
	if err != nil {
		return err
	}
	*s = firstStringFromSliceAsFloat64(f)
	return nil
}

func (s *firstStringFromSliceAsFloat64) val() float64 {
	return float64(*s)
}

//nolint:unused
type stringAsUnixTimestamp time.Time

//nolint:unused
func (s *stringAsUnixTimestamp) UnmarshalJSON(bytes []byte) error {
	var ss string
	if err := json.Unmarshal(bytes, &ss); err != nil {
		return err
	}
	i, err := strconv.ParseInt(ss, 10, 64)
	if err != nil {
		return err
	}
	*s = stringAsUnixTimestamp(time.Unix(i, 0))
	return nil
}

//nolint:unused
func (s *stringAsUnixTimestamp) val() time.Time {
	return time.Time(*s)
}

//nolint:unused
type stringAsInt64 int64

//nolint:unused
func (s *stringAsInt64) UnmarshalJSON(bytes []byte) error {
	var ss string
	if err := json.Unmarshal(bytes, &ss); err != nil {
		return err
	}
	i, err := strconv.ParseInt(ss, 10, 64)
	if err != nil {
		return err
	}
	*s = stringAsInt64(i)
	return nil
}

//nolint:unused
func (s *stringAsInt64) val() int64 {
	return int64(*s)
}

type intAsUnixTimestamp time.Time

func (s *intAsUnixTimestamp) UnmarshalJSON(bytes []byte) error {
	var i int64
	if err := json.Unmarshal(bytes, &i); err != nil {
		return err
	}
	*s = intAsUnixTimestamp(time.Unix(i, 0))
	return nil
}
func (s *intAsUnixTimestamp) val() time.Time {
	return time.Time(*s)
}

type intAsUnixTimestampMs time.Time

func (s *intAsUnixTimestampMs) UnmarshalJSON(bytes []byte) error {
	var i int64
	if err := json.Unmarshal(bytes, &i); err != nil {
		return err
	}
	*s = intAsUnixTimestampMs(time.Unix(i/1000, 0))
	return nil
}
func (s *intAsUnixTimestampMs) val() time.Time {
	return time.Time(*s)
}
