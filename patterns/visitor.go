package patterns

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type Visitor func(Shape)
type Shape interface {
	Accept(Visitor)
}

type Circle struct {
	R int
}

type Rectangle struct {
	W int
	H int
}

func (c *Circle) Accept(v Visitor) {
	v(c)
}

func (c *Rectangle) Accept(v Visitor) {
	v(c)
}

func JsonVisitor(s Shape) {
	bytes, err := json.Marshal(s)
	if err == nil {
		fmt.Println(string(bytes))
	}
}

func XmlVisitor(s Shape) {
	bytes, err := xml.Marshal(s)
	if err == nil {
		fmt.Println(string(bytes))
	}
}
