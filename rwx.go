package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/beevik/etree"
)

//Node is node
type Node struct {
	UID       string   `json:"id"`
	Name      string   `json:"name"`
	Comment   *string  `json:"comment,omitempty"`
	Attribute *string  `json:"attribute,omitempty"`
	Value     *string  `json:"value,omitempty"`
	Nodes     []string `json:"nodes"`
}

//InsertReq is InsertReq
type InsertReq struct {
	XPath string   `json:"xpath"`
	Nodes []string `json:"nodes"`
}

//DeleteReq is DeleteReq
type DeleteReq struct {
	XPath string `json:"xpath"`
}

//UpdateReq is UpdateReq
type UpdateReq struct {
	XPath     string  `json:"xpath"`
	Value     *string `json:"value,omitempty"`
	Replace   *string `json:"replace,omitempty"`
	Attribute *string `json:"attribute,omitempty"`
}

//JSONData is JSONData
type JSONData struct {
	Nodes   []Node      `json:"nodes"`
	Inserts []InsertReq `json:"insert"`
	Deletes []DeleteReq `json:"delete"`
	Updates []UpdateReq `json:"update"`
}

var nodes map[string]Node

func readJSONFile(jsonFile string) *JSONData {
	var j JSONData
	f, _ := ioutil.ReadFile(jsonFile)
	data := []byte(f)
	if e := json.Unmarshal(data, &j); e != nil {
		panic(e)
	}
	return &j
}

func makeNodes(root *etree.Element, n Node) {
	if root == nil || len(n.Name) == 0 {
		return
	}

	newElem := root.CreateElement(n.Name)
	if n.Attribute != nil {
		keyVals := strings.Split(*n.Attribute, "=")
		if len(keyVals) > 1 {
			newElem.CreateAttr(keyVals[0], keyVals[1])
		}
	}

	if n.Value != nil {
		newElem.CreateText(*n.Value)
	}
	if n.Comment != nil {
		newElem.CreateComment(*n.Comment)
	}

	if len(n.Nodes) != 0 {
		for _, s := range n.Nodes {
			childNode := nodes[s]
			makeNodes(newElem, childNode)
		}
	}
}

func insert(j []InsertReq, xmlFile string) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(xmlFile); err != nil {
		panic(err)
	}

	for _, inserts := range j {
		fmt.Println(" insert req: ", inserts.XPath, inserts.Nodes)
		for _, n := range inserts.Nodes {
			pathOfElem := etree.MustCompilePath(inserts.XPath)
			elem := doc.FindElementPath(pathOfElem)
			node := nodes[n]
			makeNodes(elem, node)
		}
	}
	doc.Indent(2)
	doc.WriteTo(os.Stdout)
	doc.WriteToFile(xmlFile)
}

func update(j []UpdateReq, xmlFile string) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(xmlFile); err != nil {
		panic(err)
	}

	for _, update := range j {
		pathOfElem := etree.MustCompilePath(update.XPath)
		elem := doc.FindElementPath(pathOfElem)
		if elem != nil {
			if update.Attribute != nil {
				keyVal := strings.Split(*update.Attribute, "=")
				if len(keyVal) > 1 {
					elem.RemoveAttr(keyVal[0])
					elem.CreateAttr(keyVal[0], keyVal[1])
				}
			}
			if update.Value != nil {
				elem.SetText(*update.Value)
			}
			if update.Replace != nil {
				vals := strings.Split(*update.Replace, "=")
				if len(vals) > 1 {
					text := elem.Text()
					if strings.Contains(text, vals[0]) {
						newText := strings.Replace(text, vals[0], vals[1], -1)
						elem.SetText(newText)
					}
				}
			}
		}
	}
	doc.Indent(2)
	doc.WriteTo(os.Stdout)
	// doc.WriteToFile(xmlFile)
}

func delete(j []DeleteReq, xmlFile string) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(xmlFile); err != nil {
		panic(err)
	}
	doc.Indent(2)

	for _, updates := range j {
		pathOfElem := etree.MustCompilePath(updates.XPath)
		elem := doc.FindElementPath(pathOfElem)
		if elem != nil {
			parent := elem.Parent()
			if parent != nil {
				fmt.Println("removing : ", elem.GetPath(), elem.Index())
				parent.RemoveChildAt(elem.Index())
			}

		}
	}
	doc.Indent(2)
	// doc.WriteTo(os.Stdout)
	doc.WriteToFile(xmlFile)
}

func main() {
	jData := readJSONFile("j.json")
	nodes = make(map[string]Node, 0)
	for _, node := range jData.Nodes {
		nodes[node.UID] = node
	}

	// *Inserts *
	// insert(jData.Inserts, "./sample.xml")

	// *Updates *
	update(jData.Updates, "./sample.xml")

	// // *Deletes *
	// delete(jData.Deletes, "./sample.xml")
}
