package qp

import (
	"log"
	"regexp"
	"strings"

	wb "github.com/Kard34/search-engine/dataxet/iq-wordbreak"
)

func init() {
	// datadic, _ := os.ReadFile("godic.txt")
	// _, dict = wb.LoadDictString(string(datadic), "")

}

var dict *wb.TriMain

type FlatNode struct {
	Idx int    `json:"index"`
	Lt  int    `json:"left"`
	Rt  int    `json:"right"`
	Val string `json:"value"`
}

func Flatten(n ND, hilight *string, arr *[]FlatNode) int {
	firstIdx := len(*arr)
	// log.Println(firstIdx)

	opIdx := firstIdx
	if n.hasNext() { // has " Op Right"
		// fmt.Println(next, nil, next == nil, next != nil)
		if n.getOp() == "and not" {
			*arr = append(*arr, FlatNode{Idx: len(*arr), Lt: -1, Rt: -1, Val: "and"})
		} else {
			*arr = append(*arr, FlatNode{Idx: len(*arr), Lt: -1, Rt: -1, Val: n.getOp()})
		}
	}

	lt := len(*arr)
	if n.isLeaf() { // >>L<<. . .
		// check phrase (Word Break) here
		subwords := []string{n.val()}

		if n.breakable() {
			if dict != nil {
				broke := wb.BreakLine(dict, n.val(), false, false, false)
				subwords = strings.FieldsFunc(broke, Split)
			} else {
				subwords = strings.FieldsFunc(n.val(), Split)
			}
		}
		var i int
		for i = 0; i < len(subwords)-1; i++ {
			if isLetter(subwords[i]) {
				*arr = append(*arr, FlatNode{Idx: len(*arr), Lt: len(*arr) + 1, Rt: len(*arr) + 2, Val: "phrase2"})
				*arr = append(*arr, FlatNode{Idx: len(*arr), Lt: -1, Rt: -1, Val: subwords[i]})
				*hilight += subwords[i]
			} else {
				if len(*arr) > 1 {
					(*arr)[len(*arr)-2].Val = "phrase3"
				}
			}
		}
		if n.breakable() {
			if isLetter(subwords[i]) {
				*arr = append(*arr, FlatNode{Idx: len(*arr), Lt: -1, Rt: -1, Val: subwords[i]})
				*hilight += subwords[i]
			} else {
				if len(*arr) > 1 {
					(*arr)[len(*arr)-1].Idx = (*arr)[len(*arr)-2].Idx
					(*arr)[len(*arr)-2] = (*arr)[len(*arr)-1]
					*arr = (*arr)[:len(*arr)-1]
				}
			}
		} else {
			*arr = append(*arr, FlatNode{Idx: len(*arr), Lt: -1, Rt: -1, Val: subwords[i]})
			*hilight += subwords[i]
		}

	} else { // >>(L)<<. . .
		lt = Flatten(n.getLeft(), hilight, arr)
	}

	if n.hasNext() { // . . .>>& B<<
		if opIdx == len(*arr) {
			log.Println(*arr)
			return lt
		}
		(*arr)[opIdx].Lt = lt
		if n.getOp() == "and not" {
			(*arr)[opIdx].Rt = len(*arr)
			opIdx = len(*arr)
			*arr = append(*arr, FlatNode{Idx: len(*arr), Lt: -1, Rt: -1, Val: "not"})
		}
		next := n.getNext()
		rt := Flatten(next, hilight, arr)
		(*arr)[opIdx].Rt = rt
		// (*arr)[opIdx].Rt = Flatten(*n.Next, arr)

	}
	// log.Println(hilight)
	*hilight += n.val()
	return firstIdx
}

func Split(r rune) bool {
	return r == '\u00a0' || r == ' '
}

var isLetter = regexp.MustCompile(`^[0-9a-zA-Zก-ฮะ-ํ]+$`).MatchString

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// func MakeOrderStr(qa []queryparser.ArrayNode, i int, midPhrase bool) string {
// 	if qa[i].Lt == -1 && qa[i].Rt == -1 {
// 		return " " + qa[i].Val + " "
// 	}
// 	if qa[i].Val == "phrase2" || qa[i].Val == "phrase3" {
// 		var s string
// 		if !midPhrase {
// 			midPhrase = true
// 			s = `"`
// 		}
// 		s += MakeOrderStr(qa, int(qa[i].Lt), midPhrase)
// 		if qa[i].Val == "phrase3" {
// 			s += "-"
// 		}
// 		s += MakeOrderStr(qa, int(qa[i].Rt), midPhrase)
// 		if qa[qa[i].Rt].Val != "phrase2" && qa[qa[i].Rt].Val != "phrase3" {
// 			s += `"`
// 		}
// 		return s
// 	} else {
// 		var s string
// 		if qa[i].Val != "not" {
// 			s = "("
// 		}
// 		if qa[i].Lt != -1 {
// 			s += MakeOrderStr(qa, int(qa[i].Lt), false)
// 		}
// 		s += " " + qa[i].Val + " "
// 		if qa[i].Rt != -1 {
// 			s += MakeOrderStr(qa, int(qa[i].Rt), false)
// 		}
// 		if qa[i].Val != "not" {
// 			s += ")"
// 		}
// 		return s
// 	}
// }
