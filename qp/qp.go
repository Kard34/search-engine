package qp

import (
	"fmt"
	"strings"

	wb "github.com/Kard34/search-engine/dataxet/iq-wordbreak"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/stateful"
)

var queryLexer = lexer.Must(stateful.NewSimple([]stateful.Rule{
	//order matters
	{Name: "OpAndNot", Pattern: "\\band not\\b|\\Bและไม่\\B", Action: nil},
	{Name: "OpAnd", Pattern: "\\band\\b|\\Bและ\\B", Action: nil},
	{Name: "OpOr", Pattern: "\\bor\\b|\\Bหรือ\\B", Action: nil},
	{Name: "Phr", Pattern: "\\\"[\u00a0 !\u0023-\u007eก-ํ]*\\\"", Action: nil}, //"anything except["]"
	{Name: "Special", Pattern: "[%#:]", Action: nil},
	{Name: "Str", Pattern: "[\u00a0\u0026\u002b-\u0039\u003b-\u007eก-ํ]+", Action: nil},
	//{Name: "Str", Pattern: "[:+*/_.a-zA-Z0-9ก-์\u00a0%][-:+*/_.a-zA-Z0-9ก-์\u00a0%]*", Action: nil}, //
	{Name: "Sym", Pattern: `-`, Action: nil},
	{Name: "Pt", Pattern: `[()]`, Action: nil},
	{Name: "_", Pattern: `\s+`, Action: nil},
}))

var parser = participle.MustBuild(&Node{},
	participle.Lexer(queryLexer),
	participle.UseLookahead(2),
)

type Node struct {
	Left *NodeAnd `@@`
	Op   *string  `[ _@OpOr`
	Next *Node    ` _@@ ]`
}

type NodeAnd struct {
	Left *NodeAndNot `@@`
	Op   *string     `[ (_@OpAnd)?`
	Next *NodeAnd    ` _?@@ ]` //aaa(bbb) case
}

type NodeAndNot struct {
	Val  *Value      `@@`
	Op   *string     `[ _@OpAndNot` //`[ _(@OpAndNot _|"-")`
	Next *NodeAndNot `_@@ ]`
}

type Value struct {
	Sub    *Node   `"(" _? @@ _? ")"`
	Tag    *string ` | (@Str ":" | @Str "%" | @"#" | @"%")?`
	Phrase *string ` (@Phr`
	Str    *string ` | @Str)`
}

func (v Value) String() string {
	if v.Phrase != nil {
		return (*v.Phrase)[1 : len(*v.Phrase)-1]
	} else if v.Tag != nil {
		// process tag here
		switch strings.ToLower(*v.Tag) {
		case "p":
			return "%source-product-" + *v.Str + "-n"
		case "g":
			return "%source-" + *v.Str + "-n"
		case "h":
			return "%" + *v.Str + "-h"
		case "o":
			return "%" + *v.Str + "-o"
		case "v":
			return "%source-vendor-" + *v.Str + "-n"
		case "#":
			return "#" + *v.Str
		case "%":
			return "%" + *v.Str
		case "cat":
			return "%cat-" + *v.Str + "-n"
		case "fb":
			return "%social-fb-" + *v.Str + "-n"
		case "ln":
			return "%social-ln-" + *v.Str + "-n"
		case "ig":
			return "%social-ig-" + *v.Str + "-n"
		case "tw":
			return "%social-tw-" + *v.Str + "-n"
		case "tt":
			return "%social-tt-" + *v.Str + "-n"
		case "yt":
			return "%social-yt-" + *v.Str + "-n"
		case "social":
			return "%social-" + *v.Str + "-n"
		case "site":
			return "%site-" + *v.Str + "-n"
		case "issd":
			return "%issd-" + *v.Str + "-n"
		default:
			return "%" + *v.Str + "-" + string((*v.Tag))
		}
	} else if v.Str != nil {
		return *v.Str
	} else {
		return ""
	}
}

// return L type
func Parse(tx *wb.TriMain, s string) ([]FlatNode, string, error) {
	dict = tx
	q := &Node{}
	s = standardizeSpaces(s)
	s = strings.ToLower(s)
	err := parser.ParseString(s, q)
	if err != nil {
		return []FlatNode{}, "", err
	}
	// log.Println(q.Op)
	flat := []FlatNode{}
	var hilight string
	_ = Flatten(q, &hilight, &flat)
	// log.Println("hilight:", hilight)
	return flat, "", nil
}

func Run(s string) error {
	// broke := wb.BreakLine(dict, "นา_ยกรัฐมนตรี  .-__--.     ..     โควิด", false)
	// arr := strings.FieldsFunc(broke, Split)
	// for _, w := range arr {
	// 	fmt.Println(w)
	// }
	// s = "นา_ยกรัฐมนตรี and อาจารย์  .-__--.    or ..  โควิด-19"
	flat, _, err := Parse(nil, s)
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range flat {
		fmt.Printf("%+v\n", f)
	}

	// q.JustPrint(1)

	return nil
}
