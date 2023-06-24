package browser

import (
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"strings"

	"rogchap.com/v8go"

	"github.com/fogleman/gg"
	"golang.org/x/net/html"
)

type Browser interface {
	NavigateToURL(url string) (*Node, error)
	RenderDOMTree(n *Node, x, y float64, dc *gg.Context, js *v8go.Context, baseURL string) float64
}

type browserImpl struct{}

func NewBrowser() Browser {
	return &browserImpl{}
}

func (b *browserImpl) RenderDOMTree(n *Node, x, y float64, dc *gg.Context, js *v8go.Context, baseURL string) float64 {
	const defaultFontSize = 20.0
	const lineHeight = 80.0

	if x == 0 {
		x = 10
	}

	if n.Type == html.TextNode {
		text := n.Data

		err := dc.LoadFontFace("/System/Library/Fonts/NewYork.ttf", defaultFontSize)
		if err != nil {
			log.Fatalf("LoadFontFace: %v", err)
		}

		r, g, b := 0, 0, 0
		if colorStr, ok := n.Styles["color"]; ok {
			switch colorStr {
			case "red":
				r, g, b = 255, 0, 0
			case "blue":
				r, g, b = 0, 0, 255
			}
		}
		dc.SetColor(color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255})

		words := strings.Fields(text)
		spaceWidth, _ := dc.MeasureString(" ")

		for _, word := range words {
			wordWidth, _ := dc.MeasureString(word)
			if x+wordWidth > 1080 {
				x = 10
				y += lineHeight
			}
			dc.DrawString(word, x, y)
			x += wordWidth + spaceWidth
		}
	}

	if n.Type == html.ElementNode && n.Data == "script" {
		if len(n.Children) > 0 && n.Children[0].Type == html.TextNode {
			jsCode := n.Children[0].Data
			_, err := js.RunScript(jsCode, "inline")
			if err != nil {
				log.Printf("JavaScript error: %v", err)
			}
		}
	}

	for _, child := range n.Children {
		y = b.RenderDOMTree(child, x, y, dc, js, baseURL)
	}

	return y
}

func (b *browserImpl) NavigateToURL(url string) (*Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	domTree := buildDOMTree(doc)
	return domTree, nil
}

func buildDOMTree(n *html.Node) *Node {
	node := &Node{
		Type:       n.Type,
		Data:       n.Data,
		Attributes: n.Attr,
		Styles:     extractStyles(n),
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		node.Children = append(node.Children, buildDOMTree(c))
	}

	return node
}

func extractStyles(n *html.Node) map[string]string {
	styles := make(map[string]string)
	for _, attr := range n.Attr {
		if attr.Key == "style" {
			styleAttributes := strings.Split(attr.Val, ";")
			for _, styleAttr := range styleAttributes {
				if style := strings.Split(styleAttr, ":"); len(style) == 2 {
					key := strings.TrimSpace(style[0])
					value := strings.TrimSpace(style[1])
					styles[key] = value
				}
			}
		}
	}
	return styles
}
