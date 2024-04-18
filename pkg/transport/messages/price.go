package messages

import (
	"encoding/json"
	"errors"

	"github.com/toknowwhy/theunit-oracle/pkg/oracle"
)

var PriceMessageName = "price/v0"

var ErrPriceMalformedMessage = errors.New("malformed price message")

type Price struct {
	Price *oracle.Price   `json:"price"`
	Trace json.RawMessage `json:"trace"`
}

func (p *Price) Marshall() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Price) Unmarshall(b []byte) error {
	err := json.Unmarshal(b, p)
	if err != nil {
		return err
	}
	if p.Price == nil {
		return ErrPriceMalformedMessage
	}
	return nil
}

func (p *Price) MarshalBinary() ([]byte, error) {
	return p.Marshall()
}

func (p *Price) UnmarshalBinary(data []byte) error {
	return p.Unmarshall(data)
}
