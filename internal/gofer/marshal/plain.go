package marshal

import (
	"fmt"
	"io"
	"strings"

	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
)

type plainItem struct {
	writer io.Writer
	item   []byte
}

type plain struct {
	items []plainItem
}

func newPlain() *plain {
	return &plain{}
}

// Write implements the Marshaller interface.
func (p *plain) Write(writer io.Writer, item interface{}) error {
	var i []byte
	switch typedItem := item.(type) {
	case *gofer.Price:
		i = p.handlePrice(typedItem)
	case *gofer.Model:
		i = p.handleModel(typedItem)
	case error:
		i = []byte(fmt.Sprintf("Error: %s", typedItem.Error()))
	default:
		return fmt.Errorf("unsupported data type")
	}

	p.items = append(p.items, plainItem{writer: writer, item: i})
	return nil
}

// Flush implements the Marshaller interface.
func (p *plain) Flush() error {
	var err error
	for _, i := range p.items {
		_, err = i.writer.Write(i.item)
		if err != nil {
			return err
		}
		_, err = i.writer.Write([]byte{'\n'})
		if err != nil {
			return err
		}
	}
	return nil
}

func (*plain) handlePrice(price *gofer.Price) []byte {
	if price.Error != "" {
		return []byte(fmt.Sprintf("%s - %s", price.Pair, strings.TrimSpace(price.Error)))
	}
	return []byte(fmt.Sprintf("%s %f", price.Pair, price.Price))
}

func (*plain) handleModel(node *gofer.Model) []byte {
	return []byte(node.Pair.String())
}
