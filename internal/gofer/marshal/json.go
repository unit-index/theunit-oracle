package marshal

import (
	encodingJSON "encoding/json"
	"fmt"
	"github.com/toknowwhy/theunit-oracle/pkg/unit"
	"io"
	"time"

	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
)

type jsonItem struct {
	writer io.Writer
	item   interface{}
}

type json struct {
	ndjson bool
	items  []jsonItem
}

func newJSON(ndjson bool) *json {
	return &json{
		ndjson: ndjson,
	}
}

// Write implements the Marshaller interface.
func (j *json) Write(writer io.Writer, item interface{}) error {
	var i interface{}
	switch typedItem := item.(type) {
	case *gofer.Price:
		i = j.handlePrice(typedItem)
	case *gofer.Model:
		i = j.handleModel(typedItem)
	case *unit.UnitPerMonthParams:
		i = j.handleUnitParams(typedItem)
	case error:
		i = j.handleError(typedItem)
	default:
		return fmt.Errorf("unsupported data type")
	}

	j.items = append(j.items, jsonItem{writer: writer, item: i})
	return nil
}

// Flush implements the Marshaller interface.
func (j *json) Flush() error {
	if j.ndjson {
		for _, item := range j.items {
			bts, err := encodingJSON.Marshal(item.item)
			if err != nil {
				return err
			}
			_, err = item.writer.Write(bts)
			if err != nil {
				return err
			}
			_, err = item.writer.Write([]byte{'\n'})
			if err != nil {
				return err
			}
		}
	} else {
		items := map[io.Writer][]interface{}{}
		for _, i := range j.items {
			items[i.writer] = append(items[i.writer], i.item)
		}
		for w, is := range items {
			bts, err := encodingJSON.Marshal(is)
			if err != nil {
				return err
			}
			_, err = w.Write(bts)
			if err != nil {
				return err
			}
			_, err = w.Write([]byte{'\n'})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (*json) handlePrice(price *gofer.Price) interface{} {
	return jsonPriceFromGoferPrice(price)
}

func (*json) handleModel(node *gofer.Model) interface{} {
	return node.Pair.String()
}

func (*json) handleError(err error) interface{} {
	return struct {
		Error string `json:"error"`
	}{Error: err.Error()}
}

func (*json) handleUnitParams(params *unit.UnitPerMonthParams) interface{} {
	return jsonUnitParams{
		CSupply:       params.CSupply,
		LastPrice:     params.LastPrice,
		LastMarketCap: params.LastMarketCap,
	}
}

type jsonUnitParams struct {
	CSupply       float64 `json:"CSupply"`
	LastMarketCap float64 `json:"LastMarketCap"`
	LastPrice     float64 `json:"LastPrice"`
}

type jsonPrice struct {
	Type       string            `json:"type"`
	Base       string            `json:"base"`
	Quote      string            `json:"quote"`
	Price      float64           `json:"price"`
	Bid        float64           `json:"bid"`
	Ask        float64           `json:"ask"`
	Volume24h  float64           `json:"vol24h"`
	Timestamp  time.Time         `json:"ts"`
	Parameters map[string]string `json:"params,omitempty"`
	Prices     []jsonPrice       `json:"prices,omitempty"`
	Error      string            `json:"error,omitempty"`
}

func jsonPriceFromGoferPrice(t *gofer.Price) jsonPrice {
	var prices []jsonPrice
	for _, c := range t.Prices {
		prices = append(prices, jsonPriceFromGoferPrice(c))
	}
	return jsonPrice{
		Type:       t.Type,
		Base:       t.Pair.Base,
		Quote:      t.Pair.Quote,
		Price:      t.Price,
		Bid:        t.Bid,
		Ask:        t.Ask,
		Volume24h:  t.Volume24h,
		Timestamp:  t.Time.In(time.UTC),
		Parameters: t.Parameters,
		Prices:     prices,
		Error:      t.Error,
	}
}
