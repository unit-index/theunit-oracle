package messages

import (
	"encoding/json"
	"errors"

	"github.com/toknowwhy/theunit-oracle/pkg/oracle"
)

var UnitParamsMessageName = "unitParams/v0"

var ErrUnitParamsMalformedMessage = errors.New("malformed UnitParams message")

type UnitParams struct {
	UnitParams *oracle.UnitParams `json:"unitParams"`
	Trace      json.RawMessage    `json:"trace"`
}

func (u *UnitParams) Marshall() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UnitParams) Unmarshall(b []byte) error {
	err := json.Unmarshal(b, u)
	if err != nil {
		return err
	}
	if u.UnitParams == nil {
		return ErrUnitParamsMalformedMessage
	}
	return nil
}

func (u *UnitParams) MarshalBinary() ([]byte, error) {
	return u.Marshall()
}

func (u *UnitParams) UnmarshalBinary(data []byte) error {
	return u.Unmarshall(data)
}
