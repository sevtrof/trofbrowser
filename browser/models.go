package browser

import (
	"golang.org/x/net/html"
)

type Node struct {
	Type       html.NodeType
	Data       string
	Attributes []html.Attribute
	Children   []*Node
	Styles     map[string]string
	X, Y       float64
}

type Color struct {
	R, G, B uint8
	A       uint8
}

type Font struct {
	Path  string
	Size  float64
	Color Color
}
