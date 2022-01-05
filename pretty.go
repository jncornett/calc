package calc

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/jncornett/calc/lang"
)

func Fprint(w io.Writer, node lang.Node) (err error) {
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	defer func() {
		if ferr := enc.Flush(); ferr != nil {
			if err == nil {
				err = ferr
			}
		}
	}()
	v := lang.VisitorFuncs{
		PushFunc: func(_ interface{}, n lang.Node) error {
			switch t := n.(type) {
			default:
				return enc.EncodeToken(xml.Comment(fmt.Sprintf("unknown node type: %T", n)))
			case lang.Do:
				return enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "do"}})
			case *lang.Assign:
				return enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "assign"}})
			case *lang.Call:
				return enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "call"}})
			case *lang.Object:
				return enc.EncodeElement(t.GoValue, xml.StartElement{Name: xml.Name{Local: "object"}})
			case lang.Ref:
				return enc.EncodeElement(string(t), xml.StartElement{Name: xml.Name{Local: "ref"}})
			case lang.Vector:
				return enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "vector"}})
			}
		},
		PopFunc: func(_ interface{}, n lang.Node) error {
			switch n.(type) {
			default:
				return nil
			case lang.Do:
				return enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "do"}})
			case *lang.Assign:
				return enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "assign"}})
			case *lang.Call:
				return enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "call"}})
			case lang.Vector:
				return enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "vector"}})
			}
		},
	}
	return node.Walk(v)
}
