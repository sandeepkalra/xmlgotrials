package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

func Check(e error) {
	if e != nil {
		fmt.Println(e.Error())
		panic(e)
	}
}

type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:"-"`
	Content []byte     `xml:",innerxml"`
	Nodes   []Node     `xml:",any"`
}

const (
	dataStr = `
	<?xml version = "1.0" encoding = "utf-8"?>
	<?xml  xslplane.1.xml="12" ?>
	<?xml-stylesheet type = "text/xsl"  href = "xslplane.1.xsl" ?>
	<plane>
	<year> 1977 </year>
	<make> Cessna </make>
	<model> Skyhawk </model>
	<color> Light blue and white </color>
	</plane>
	<outer>
		<person1>One</person1>
		<person11>One One</person11>
	</outer>
	`
)

func find(path string) {
	data := []byte(dataStr)
	dec := xml.NewDecoder(bytes.NewReader(data))
	for {
		tok, e := dec.Token()
		if e == io.EOF {
			break
		} else {
			Check(e)
		}

		/* Locate the top level Element Node */
		switch tok.(type) {
		case xml.CharData:
			fmt.Println("Content>", string(tok.(xml.CharData)), len(string(tok.(xml.CharData))), "<-")
		case xml.Comment:
			fmt.Println("Comment>", string(tok.(xml.Comment)), len(string(tok.(xml.Comment))), "<-")
		case xml.Directive:
			fmt.Println("Directive>", string(tok.(xml.Directive)), len(string(tok.(xml.Directive))), "<-")
		case xml.ProcInst:
			fmt.Println("ProcInst>", tok.(xml.ProcInst).Target, len(string(tok.(xml.ProcInst).Target)), "<-")
		case xml.StartElement:
			fmt.Println("StartElement>", tok.(xml.StartElement).Name, "<-")
		}
		xpath := strings.Split(path, "/")

		if tok, ok := tok.(xml.StartElement); ok {
			// are we interested in this node ?
			if xpath[0] == tok.Name.Local {
				fmt.Println("--node-->", tok.Name.Local)
				var start xml.StartElement
				var node Node
				err := dec.DecodeElement(&node, &start)
				Check(err)

				// are we interested in this sub-node ?
				for _, n := range node.Nodes {
					fmt.Println("= node,len, looking-for,len => ", n.XMLName.Local, len(n.XMLName.Local), xpath[1], len(xpath[1]))
					if strings.Compare(n.XMLName.Local, xpath[1]) == 0 {
						fmt.Println(" ::: node found >>>", string(n.XMLName.Local))
						fmt.Println(" ::: node Val >>>", string(n.Content))
						fmt.Println(" ::: node content >>> ", string(n.Content))
						fmt.Println(" ::: node attribute count >>> ", len(n.Attrs))
						fmt.Println(" ::: node sub-node count >>> ", len(n.Nodes))
					}
				}
			}
		}
		// else {
		// fmt.Println("Not found")
		// }
	}
}

// func main() {
// 	find("outer/person11") ;; after you run this, grep for 'node' to check result
// }
