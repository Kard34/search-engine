package iq_wordbreak

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"unicode/utf8"
)

type TriMain struct {
	root *TriNode
	idx  int
	word []string
}
type TriNode struct {
	Charator string
	IsWord   bool
	IsLast   bool
	NextNode *TriNode
	LinkNode *TriNode

	Index     int
	IsTop     bool
	IsHead    bool
	IsReach   bool
	FromCount int
	IsCommon  bool
}

func (n *TriNode) NewNextLinkWord(arrWord []string, istart int, idx *int) {

	if n == nil {
		return
	}

	if n.NextNode == nil {

		n.NextNode = &TriNode{
			Charator:  arrWord[istart],
			IsWord:    true,
			IsLast:    true,
			NextNode:  nil,
			LinkNode:  nil,
			Index:     *idx,
			IsTop:     true,
			IsHead:    false,
			IsReach:   false,
			FromCount: 0,
			IsCommon:  false,
		}

		*idx++
		istart++

		fmt.Println("3GO|", n.Print())

		if istart < len(arrWord) {
			n.IsWord = false
			n.NewWord(arrWord, istart, idx)
		}
		return
	} else {

		return
	}

}

func Byte2String874(buff []byte) string {
	var s strings.Builder
	for _, v := range buff {
		if v > 160 {
			s.WriteString(string(rune((int(v-161) + 3585))))
		} else {
			s.WriteString(string(v))
		}
	}
	return s.String()
}

func String2Byte874String(data string) string {
	slen := utf8.RuneCountInString(data)

	slice := make([]byte, slen)
	pos := 0
	for _, char := range data {
		if char > 3584 && char < 3676 {
			slice[pos] = byte((int(char)) - 3585 + 161)
		} else {
			slice[pos] = byte((int(char)))

		}
		pos++
		

	}
	str := hex.EncodeToString(slice)
	return strings.ToUpper(str)
}

func (n *TriNode) Print() string {
	if n == nil {
		return ""
	}
	var sb strings.Builder
	vv := String2Byte874String(n.Charator)
	sb.WriteString(fmt.Sprint("|", n.Index, ":", vv, ":"))
	if n.NextNode != nil {
		sb.WriteString(fmt.Sprint(n.NextNode.Index))
	} else {
		sb.WriteString("N")
	}
	sb.WriteString(":")
	if n.LinkNode != nil {
		sb.WriteString(fmt.Sprint(n.LinkNode.Index))
	} else {
		sb.WriteString("L")
	}
	sb.WriteString(":")
	sb.WriteString(fmt.Sprint(n.FromCount))

	sb.WriteString(":")
	if n.IsWord {
		sb.WriteString("T")
	} else {
		sb.WriteString("F")
	}
	if n.IsLast {
		sb.WriteString("T")
	} else {
		sb.WriteString("F")
	}
	if n.IsTop {
		sb.WriteString("T")
	} else {
		sb.WriteString("F")
	}
	if n.IsHead {
		sb.WriteString("T")
	} else {
		sb.WriteString("F")
	}
	if n.IsReach {
		sb.WriteString("T")
	} else {
		sb.WriteString("F")
	}
	if n.IsCommon {
		sb.WriteString("T")
	} else {
		sb.WriteString("F")
	}

	if n.NextNode != nil {
		sb.WriteString(n.NextNode.Print())
	}
	if n.LinkNode != nil {
		sb.WriteString(n.LinkNode.Print())
	}

	return sb.String()
}
func (t *TriMain) Print() string {
	var sb strings.Builder
	node := t.root
	if node != nil {
		sb.WriteString(node.Print())
	}

	return sb.String()
}
func (n *TriNode) NewWord(arrWord []string, istart int, idx *int) {

	if n == nil {
		return
	}

	if n.LinkNode == nil {

		n.LinkNode = &TriNode{
			Charator:  arrWord[istart],
			IsWord:    true,
			IsLast:    true,
			NextNode:  nil,
			LinkNode:  nil,
			Index:     *idx,
			IsTop:     true,
			IsHead:    false,
			IsReach:   false,
			FromCount: 0,
			IsCommon:  false,
		}

		*idx++
		istart++
		if istart < len(arrWord) {
			n.LinkNode.IsWord = false
			n.LinkNode.NewWord(arrWord, istart, idx)
		}
		return
	} else {

		return
	}

}
func (n *TriNode) FindOnlyNode(arrWord []string, istart int) *TriNode {

	if n == nil {
		return nil
	}

	tri := n
	for {
		if tri.NextNode != nil {
			return nil
		}

		if tri.Charator != arrWord[istart] {
			return nil
		}

		tri = tri.LinkNode
		istart++
		if tri == nil && istart != len(arrWord) {
			return nil
		}

		if istart+1 == len(arrWord) {
			if tri.NextNode == nil && tri.LinkNode == nil && tri.Charator == arrWord[istart] {
				return n
			}
			break
		}

	}
	return nil
}
func (n *TriNode) FindNode(arrWord []string, istart int) *TriNode {
	if n == nil {
		return nil
	}

	if n.Charator == arrWord[istart] {

		if istart+1 == len(arrWord) {
			if n.LinkNode == nil && n.NextNode == nil {
				return n
			}
		} else {

			nt := n.FindOnlyNode(arrWord, istart)
			if nt != nil {
				return nt
			}
		}
	}
	if n.LinkNode != nil {
		tx := n.LinkNode.FindNode(arrWord, istart)
		if tx != nil {
			return tx
		}
	}
	if n.NextNode != nil {
		return n.NextNode.FindNode(arrWord, istart)

	}
	return nil
}

func (t *TriMain) Find(arrWord []string, istart int) *TriNode {
	CurrNode := t.root
	tri := CurrNode.FindNode(arrWord, istart)
	return tri
}
func (t *TriMain) DupEnd(arrWord []string, istart int) *TriNode {
	CurrNode := t.Find(arrWord, istart)

	if CurrNode == nil {
		return nil
	}
	if CurrNode.LinkNode != nil {

		tri := CurrNode
		for tri.LinkNode != nil {
			if tri.IsWord {
				return nil
			}
			tri = tri.LinkNode
		}
		tri = CurrNode
		for tri.LinkNode != nil {
			tri.IsCommon = true
			tri = tri.LinkNode
		}
	} else {

		CurrNode.IsCommon = true
	}
	return CurrNode
}
func extra(word string) (TriNode, int, int) {
	tt := TriNode{}
	lst := strings.Split(word, ":")

	tt.Charator = lst[0]
	i, _ := strconv.Atoi(lst[0])
	tt.Index = i

	ch := new(big.Int)
	ch.SetString(lst[1], 16)

	tt.Charator = Byte2String874(ch.Bytes())
	nl := -1
	if lst[2] != "N" {
		nl, _ = strconv.Atoi(lst[2])
	}
	ll := -1
	if lst[3] != "L" {
		ll, _ = strconv.Atoi(lst[3])
	}
	tt.FromCount, _ = strconv.Atoi(lst[4])

	for idx, bb := range strings.Split(lst[5], "") {
		if idx == 0 {
			if bb == "T" {
				tt.IsWord = true
			} else {
				tt.IsWord = false
			}
		}
		if idx == 1 {
			if bb == "T" {
				tt.IsLast = true
			} else {
				tt.IsLast = false
			}
		}

		if idx == 2 {
			if bb == "T" {
				tt.IsTop = true
			} else {
				tt.IsTop = false
			}
		}
		if idx == 3 {
			if bb == "T" {
				tt.IsHead = true
			} else {
				tt.IsHead = false
			}
		}
		if idx == 4 {
			if bb == "T" {
				tt.IsReach = true
			} else {
				tt.IsReach = false
			}
		}
		if idx == 5 {
			if bb == "T" {
				tt.IsCommon = true
			} else {
				tt.IsCommon = false
			}
		}

	}
	return tt, nl, ll
}
func (nroot *TriNode) LoadLinkNode(nextl int, linkl int) {

	if nroot == nil {
		return
	}

	if nextl != -1 {
		n, nl, ll := extra(Dicword[nextl])
		nroot.NextNode = &TriNode{Charator: n.Charator,
			IsWord:   n.IsWord,
			IsLast:   n.IsLast,
			NextNode: nil,
			LinkNode: nil,

			Index:     n.Index,
			IsTop:     n.IsTop,
			IsHead:    n.IsHead,
			IsReach:   n.IsReach,
			FromCount: 0,
			IsCommon:  n.IsCommon}

		// nroot.NextNode.Load(nl, ll)
		_ = nl
		_ = ll

	}
	if linkl != -1 {
		n, nl, ll := extra(Dicword[linkl])
		nroot.LinkNode = &TriNode{Charator: n.Charator,
			IsWord:   n.IsWord,
			IsLast:   n.IsLast,
			NextNode: nil,
			LinkNode: nil,

			Index:     n.Index,
			IsTop:     n.IsTop,
			IsHead:    n.IsHead,
			IsReach:   n.IsReach,
			FromCount: 0,
			IsCommon:  n.IsCommon}
		// nroot.LinkNode.Load(nl, ll)
		_ = nl
		_ = ll
	}

}
func (nroot *TriNode) Load() {

	if nroot == nil {
		return
	}
	if nroot.NextNode != nil && nroot.LinkNode != nil {
		return
	}
	_, nl, ll := extra(Dicword[nroot.Index])
	// fmt.Println(nroot.Index, nl, ll)
	nroot.LoadLinkNode(nl, ll)
}

var Dicword []string

func (t *TriMain) Load(word []string) *TriMain {
	Dicword = word
	idx := 0
	n, nl, ll := extra(Dicword[idx])
	t.word = word
	t.root = &TriNode{
		Charator: n.Charator,
		IsWord:   n.IsWord,
		IsLast:   n.IsLast,
		NextNode: nil,
		LinkNode: nil,

		Index:     n.Index,
		IsTop:     n.IsTop,
		IsHead:    n.IsHead,
		IsReach:   n.IsReach,
		FromCount: 0,
		IsCommon:  n.IsCommon,
	}

	rr := t.root
	sk := NewStack()

	for rr != nil || sk.Length() > 0 {

		for rr != nil {

			sk.Push(rr)
			rr = rr.NextNode
		}
		if sk.Length() == 0 {
			break
		}
		rr = sk.Pop()

		rr.Load()
		if rr.NextNode != nil {
			sk.Push(rr.NextNode)
		}
		rr = rr.LinkNode

	}
	// _ = rr
	_ = nl
	_ = ll

	return t
}
func (t *TriMain) CheckFinish() *TriNode { return nil }
func (t *TriMain) Add(word string, i int) *TriMain {
	arrWord := strings.Split(word, "")

	if t.root == nil {

		istart := 0
		idx := 0
		t.root = &TriNode{
			Charator:  arrWord[istart],
			IsWord:    true,
			IsLast:    true,
			NextNode:  nil,
			LinkNode:  nil,
			Index:     idx,
			IsTop:     true,
			IsHead:    false,
			IsReach:   false,
			FromCount: 0,
			IsCommon:  false,
		}
		idx++
		istart++
		t.root.IsWord = false
		t.root.NewWord(arrWord, istart, &idx)
		t.idx = idx
		n := t.root
		for n != nil {
			n = n.LinkNode
		}
		return t
	} else {
		CurrNode := t.root
		PrevNode := t.root
		_ = CurrNode
		_ = PrevNode
		PrevNode = nil

		for i, ch := range arrWord {

			for CurrNode != nil {
				if CurrNode.Charator == arrWord[i] {
					break
				}
				PrevNode = CurrNode
				CurrNode = CurrNode.NextNode
			}
			if CurrNode == nil {
				PrevNode.IsLast = false
				idx := t.idx
				PrevNode.NewWord(arrWord, i, &idx)
				fmt.Println("6GO|", t.root.Print())

				t.idx = idx
				PrevNode = PrevNode.NextNode
				if i+1 == len(arrWord) {
					PrevNode.IsWord = true
					return t
				}
				n := i

				for { // while(true)

					// //=======
					// idx = t.idx
					PrevNode.LinkNode = t.DupEnd(arrWord, n)

					if PrevNode.LinkNode != nil {
						if PrevNode.Index == PrevNode.LinkNode.Index {
							PrevNode.LinkNode = nil
						}
					}

					if PrevNode.LinkNode != nil {
						return t
					}
					PrevNode.IsLast = false
					//PrevNo.e.LinkNode =wWordN(ByteWord, n)

					idx = t.idx
					PrevNode.NewWord(arrWord, n, &idx)
					t.idx = idx
					n++
					PrevNode = PrevNode.LinkNode
					if n == len(arrWord) {
						PrevNode.IsWord = true
						return t
					}
					//=======

				}

			}
			if i+1 == len(arrWord) {
				CurrNode.IsWord = true
			}

			if CurrNode.LinkNode != nil {
				if CurrNode.LinkNode.FromCount > 1 {
					CurrNode.FromCount = 0
				}
				if CurrNode.LinkNode.IsCommon {
					CurrNode.FromCount = 0
				}

				CurrNode = CurrNode.LinkNode
			} else {
				idx := t.idx
				fmt.Println("2GO|", t.root.Print())
				CurrNode.NewWord(arrWord, i+1, &idx)
				t.idx = idx
				fmt.Println("3GO|", t.root.Print())

				return t
			}

			_ = ch

		}
	}
	return t
}

func NewTriNode() TriNode {
	return TriNode{
		Charator: "",
		IsWord:   false,
		IsLast:   false,
		NextNode: nil,
		LinkNode: nil,

		Index:     -1,
		IsTop:     true,
		IsHead:    false,
		IsReach:   false,
		FromCount: 0,
		IsCommon:  false,
	}

}

type Pnode struct {
	char string
	next *Pnode
	link *Pnode
	flag int
}
type PTree struct {
	root *Pnode
	link int
}

func (t *PTree) InsertNext(data Pnode) *PTree {
	if t.root == nil {
		t.root = &Pnode{char: data.char, next: nil, link: nil, flag: data.flag}
	} else {
		t.root.InsertNext(data)
	}
	t.link++
	return t

}
func (n *Pnode) InsertNext(data Pnode) {
	if n == nil {
		return
	}
	if n.next == nil {

		n.next = &Pnode{char: data.char, next: nil, link: nil, flag: data.flag}
	} else {
		n.next.InsertNext(data)

	}

}
