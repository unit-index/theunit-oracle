package marshal

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/toknowwhy/theunit-oracle/internal/gofer/marshal/testutil"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
)

func TestJSON_Nodes(t *testing.T) {
	var err error
	b := &bytes.Buffer{}
	m := newJSON(false)

	ab := gofer.Pair{Base: "A", Quote: "B"}
	cd := gofer.Pair{Base: "C", Quote: "D"}
	ns := testutil.Models(ab, cd)

	err = m.Write(b, ns[ab])
	assert.NoError(t, err)

	err = m.Write(b, ns[cd])
	assert.NoError(t, err)

	err = m.Flush()
	assert.NoError(t, err)

	expected := `["A/B", "C/D"]`

	assert.JSONEq(t, expected, b.String())
}

func TestNDJSON_Nodes(t *testing.T) {
	var err error
	b := &bytes.Buffer{}
	m := newJSON(true)

	ab := gofer.Pair{Base: "A", Quote: "B"}
	cd := gofer.Pair{Base: "C", Quote: "D"}
	ns := testutil.Models(ab, cd)

	err = m.Write(b, ns[ab])
	assert.NoError(t, err)

	err = m.Write(b, ns[cd])
	assert.NoError(t, err)

	err = m.Flush()
	assert.NoError(t, err)

	result := bytes.Split(b.Bytes(), []byte("\n"))

	assert.JSONEq(t, `"A/B"`, string(result[0]))
	assert.JSONEq(t, `"C/D"`, string(result[1]))
}

func TestJSON_Prices(t *testing.T) {
	var err error
	b := &bytes.Buffer{}
	m := newJSON(false)

	ab := gofer.Pair{Base: "A", Quote: "B"}
	ts := testutil.Prices(ab)

	err = m.Write(b, ts[ab])
	assert.NoError(t, err)

	err = m.Flush()
	assert.NoError(t, err)

	expected := `
		[
		   {
			  "type":"aggregator",
			  "base":"A",
			  "quote":"B",
			  "price":10,
			  "bid":10,
			  "ask":10,
			  "vol24h":0,
			  "ts":"1970-01-01T00:00:10Z",
			  "params":{
				 "method":"median",
				 "minimumSuccessfulSources":"1"
			  },
			  "prices":[
				 {
					"type":"origin",
					"base":"A",
					"quote":"B",
					"price":10,
					"bid":10,
					"ask":10,
					"vol24h":10,
					"ts":"1970-01-01T00:00:10Z",
					"params":{
					   "origin":"a"
					}
				 },
				 {
					"type":"aggregator",
					"base":"A",
					"quote":"B",
					"price":10,
					"bid":10,
					"ask":10,
					"vol24h":10,
					"ts":"1970-01-01T00:00:10Z",
					"params":{
					   "method":"indirect"
					},
					"prices":[
					   {
						  "type":"origin",
						  "base":"A",
						  "quote":"B",
						  "price":10,
						  "bid":10,
						  "ask":10,
						  "vol24h":10,
						  "ts":"1970-01-01T00:00:10Z",
						  "params":{
							 "origin":"a"
						  }
					   }
					]
				 },
				 {
					"type":"aggregator",
					"base":"A",
					"quote":"B",
					"price":10,
					"bid":10,
					"ask":10,
					"vol24h":0,
					"ts":"1970-01-01T00:00:10Z",
					"params":{
					   "method":"median",
					   "minimumSuccessfulSources":"1"
					},
					"prices":[
					   {
						  "type":"origin",
						  "base":"A",
						  "quote":"B",
						  "price":10,
						  "bid":10,
						  "ask":10,
						  "vol24h":10,
						  "ts":"1970-01-01T00:00:10Z",
						  "params":{
							 "origin":"a"
						  }
					   },
					   {
						  "type":"origin",
						  "base":"A",
						  "quote":"B",
						  "price":20,
						  "bid":20,
						  "ask":20,
						  "vol24h":20,
						  "ts":"1970-01-01T00:00:20Z",
						  "params":{
							 "origin":"b"
						  },
						  "error":"something"
					   }
					]
				 }
			  ]
		   }
		]
	`

	assert.JSONEq(t, expected, b.String())
}
