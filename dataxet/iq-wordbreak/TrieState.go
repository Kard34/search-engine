package iq_wordbreak

import (
	"strings"
)

type TrieState struct {
	root *Trie
	dic  *TriMain
	str  string
}
type Trie struct {
	root           *TriNode
	curr           *TriNode
	prev           *TriNode
	is_crash       bool
	is_wildcard    bool
	text           []string
	start_index    int
	stop_index     int
	text_index     int
	last_word      int
	wordbreak_type string
	point          *lTrekPoint
	error_index    int
	max_word       int
	break_char     string
	output         string
	index_input    int
 
}
type lTrekPoint struct {
	lst []*TrekPoint
}
type TrekPoint struct {
	idx   int
	point int
}

func (t *TrieState) New(dic *TriMain, str string, start int, stop int) {
	t.dic = dic
	t.str = str
	t.root = &Trie{}
	t.root.root = dic.root
	t.root.start_index = start
	t.root.text_index = start
	t.root.stop_index = stop
	t.root.break_char = string(rune(160))
	t.root.text = strings.Split(str, "")
	t.root.point = &lTrekPoint{}
	t.Home()
}

func (t *TrieState) IsType(flag CharFlag) bool {

	return WBTypeIsType(t.root.text[t.root.text_index], flag)

}
func (t *TrieState) IsThai(wasThai bool) bool {
	isDot := t.IsType(CharFlag(Dot))
	isBA := t.IsType(CharFlag(BreakAlpha))
	return (wasThai || !isDot) && isBA
}
func (t *TrieState) StartIndex(i int) int {
	t.root.start_index += i
	return t.root.start_index
}

func (t *TrieState) GetOutput() string {
	t.WriteResult()

	t.root.output = t.root.output + strings.Join(t.root.text[t.root.index_input:], "")
	return t.root.output
}
func (t *TrieState) StopIndex() int {
	return t.root.stop_index
}
func (t *TrieState) AddPoint(a int, b int) {
	t.root.point.lst = append(t.root.point.lst, &TrekPoint{
		idx:   a,
		point: b,
	})
}
func (t *TrieState) ClearPoint() {
	t.root.point.lst = make([]*TrekPoint, 0)
}
func (t *TrieState) WordPointFromStart() (int, int) {

	return t.StartIndex(0), int(PointWord)
}

func (t *TrieState) SetStartIndex(val int) {
	t.root.start_index = val
}
func (t *TrieState) SetTextIndex(i int) {
	t.root.text_index = i
}

func (t *TrieState) TextIndex() int {
	return t.root.text_index
}
func (t *TrieState) MaxWord() int {
	return t.root.max_word
}
func (t *TrieState) CountWord(index int) int {
	// default value = -1

	if index == -1 {
		index = t.Tail()
	}
	NumWord := 0
	for i, c := range t.root.point.lst {
		if i > index {
			break
		}
		if c.point == int(PointWord) {
			NumWord++
		}
	}

	return NumWord
}
func (t *TrieState) Walk() int {
	curr := t.root.curr
	if t.IsEnd() || t.root.is_crash || curr == nil {
		t.root.is_crash = true
		return -11001
	}
	for curr != nil {
		if curr.Charator == t.root.text[t.root.text_index] {
			break
		}
		// fmt.Printf("%s%d|", curr.Charator, curr.Index)
		curr = curr.NextNode
	}
	// fmt.Println()

	t.root.prev = curr
	t.root.text_index++
	t.root.is_crash = (curr == nil)
	if t.root.is_crash {
		t.root.curr = nil
	} else {
		t.root.curr = curr.LinkNode
	}
	return 0
}
func (t *TrieState) IsWord() bool {

	isword := !t.root.is_crash && t.root.prev != nil && t.root.prev.IsWord
	return isword
}
func (t *TrieState) CanLead() bool {
	canlead := t.root.text_index > t.root.stop_index || !t.IsType(CharFlag(UnLeadable))
	return canlead
}
func (t *TrieState) IsCandidate() bool {
 

	b := (t.IsWord() && t.CanLead()) || t.root.is_wildcard
 
	return b
}
func (t *TrieState) SetTailType(point int) {
	l := len(t.root.point.lst)
	if l > 0 {
		t.root.point.lst[l-1].point = point
	}
}
func (t *TrieState) GetTailIndex() int {
	l := len(t.root.point.lst)
	if l > 0 {
		return t.root.point.lst[l-1].idx
	}
	return 0
}

func (t *TrieState) SetMaxWord() {
	t.root.max_word = t.CountWord(-1)
}
func (t *TrieState) TrieBreak() bool {

	// while (!(state.IsEnd() | state.IsCrash | state.IsLast()))
	for {
		IsEnd := t.IsEnd()
		IsLast := t.IsLast()
		if IsEnd || t.root.is_crash || IsLast {
			break
		}
		t.Walk()

		if t.IsCandidate() {
			t.AddPoint(t.root.text_index, int(PointCandidate))
			if t.IsEnd() {
				t.SetTailType(int(PointWord))
				t.SetMaxWord()
				return true
			}
		}
		if t.root.is_wildcard && t.root.text_index <= t.root.stop_index && t.root.text[t.root.text_index] == "*" {

			t.root.text_index++
			t.AddPoint(t.root.text_index, int(PointWord))
			t.SetMaxWord()
			return true
		}
		IsCandidate := t.IsCandidate()
		IsLast = t.IsLast()
		IsNotCandidate := !t.IsCandidate()
		IsEnd = t.IsEnd()
		Tail := t.Tail()

		if (IsCandidate && IsLast) || ((t.root.is_crash || (IsNotCandidate && IsEnd)) && t.root.last_word < Tail) {

			t.SetTailType(int(PointWord))
			t.root.text_index = t.GetTailIndex()
			t.Home()
		}

	}

	t.root.max_word = t.CountWord(-1)

	if t.Tail() >= 0 {
		t.SetTailType(int(PointWord))
		t.root.error_index = t.GetTailIndex()
	} else {
		t.root.error_index = t.root.start_index
	}

	index := t.Tail()
	BackCount := 0
	MaxTrack := 5

	for index >= 0 && BackCount < MaxTrack {

		Point := t.root.point.lst[index]

		if Point.idx <= t.StartIndex(0) {
			break
		}

		if Point.point == int(PointCandidate) {

			BackState := &TrieState{}
			BackState.New(t.dic, t.str, Point.idx, t.StopIndex())
			BackState.root.wordbreak_type = t.root.wordbreak_type

			NoErr := BackState.TrieBreak()
			_ = NoErr
			t.root.max_word = max(t.root.max_word, t.CountWord(index)+1+BackState.CountWord(-1))
			if NoErr || BackState.root.error_index > t.root.error_index {
				t.Merge(BackState, index+1)

				if NoErr {
					return true
				}
				index = t.Tail() + 1
				BackCount = -1
			}
		}
		index--
		BackCount++
	}

	return false

}

func (t *TrieState) Merge(source *TrieState, index int) {

	if index < len(t.root.point.lst) {
		t.RemoveRange(index, len(t.root.point.lst)-index)
		// t.root.point.lst = append(t.root.point.lst[0:index], t.root.point.lst[index+len(t.root.point.lst)-index:]...)
	}
	t.root.point.lst = append(t.root.point.lst, source.root.point.lst...)
	t.root.text_index = source.root.text_index
	t.root.error_index = source.root.error_index
	if index > 0 {
		t.root.point.lst[len(t.root.point.lst)-1].point = int(PointWord)
	}
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func (t *TrieState) IsLast() bool {
	return (t.root.prev != nil) && (t.root.prev.LinkNode == nil)
}
func (t *TrieState) IsEnd() bool {
	return t.root.text_index > t.root.stop_index || !t.IsType(CharFlag(BreakAlpha))
}
func (t *TrieState) Tail() int {
	return len(t.root.point.lst) - 1
}
func (t *TrieState) RemoveRange(i int, j int) {
	// a.lst = append(a.lst[0:i], a.lst[i+j:]...)
	t.root.point.lst = append(t.root.point.lst[0:i], t.root.point.lst[i+j:]...)
}
func (t *TrieState) Insert() {

}

func (t *TrieState) Home() {
	t.root.curr = t.root.root
	t.root.is_crash = false
	t.root.prev = nil
	t.root.last_word = t.Tail()
}

func (t *TrieState) WriteResult() {
	k := 0
	for t.root.index_input < len(t.root.text) {

		for k < len(t.root.point.lst) && t.root.point.lst[k].idx < t.root.index_input {
			k++
		}

		if k < len(t.root.point.lst) {

			curr := t.root.point.lst[k]
			if curr.point == int(PointWord) {

				if t.root.index_input == curr.idx &&
					(t.root.index_input > 0 && WBTypeIsType(t.root.text[t.root.index_input-1], CharFlag(Alpha)) &&
						WBTypeIsType(t.root.text[t.root.index_input], CharFlag(Alpha))) {
					t.root.output += t.root.break_char
				}

			} else {
				k++
			}
		} else {
			break
		}
		t.root.output += t.root.text[t.root.index_input]
		t.root.index_input++
	}
}
func (t *TrieState) InsertPoint(index int, a int, b int) {
	if len(t.root.point.lst) == index { // nil or empty slice or after last element
		t.root.point.lst = append(t.root.point.lst, &TrekPoint{
			idx:   a,
			point: b,
		})
		return
	}
	t.root.point.lst = append(t.root.point.lst[:index+1], t.root.point.lst[index:]...) // index < len(a)
	t.root.point.lst[index] = &TrekPoint{idx: a, point: b}

}
