package iq_wordbreak

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func readLinesDict(path string) ([]string, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	var lines []string
	info := ""
	scanner := bufio.NewScanner(file)
	ByteOrderMarkAsString := string('\uFEFF')

	for scanner.Scan() {
		oneline := strings.TrimSpace(scanner.Text())
		oneline = strings.TrimPrefix(oneline, ByteOrderMarkAsString)

		if len(oneline) < 1 {
			continue
		}
		if strings.HasPrefix(oneline, "//") {
			info = info + oneline + "\n"
		} else {
			lines = append(lines, scanner.Text())
		}
	}
	//check dict
	for i, s := range lines {

		lst := strings.Split(s, ":")
		lk, _ := strconv.Atoi(lst[0])
		if i != lk {
			var lerr []string
			return lerr, "", fmt.Errorf("FILE DICT ERROR")
		}
	}
	return lines, info, scanner.Err()
}

func B2I(b bool) int {
	if b {
		return 1
	}
	return 0
}
func Char2Int(str string) int {

	for _, char := range str {

		n := int(char)
		if n == 32 {
			return 1 // space
		} else if n >= 48 && n <= 57 {
			return 2 // number eng
		} else if n >= 65 && n <= 90 {
			return 3 // alp eng
		} else if n >= 97 && n <= 122 {
			return 3 // alp eng
		} else if n >= 1 && n <= 127 {
			return 7 // other eng
		} else if n >= 3585 && n <= 3630 {
			return 4 // alp th
		} else if n >= 3632 && n <= 3642 {
			return 4 // alp th
		} else if n >= 3648 && n <= 3662 {
			return 4 // alp th
		} else if n >= 3664 && n <= 3673 {
			return 5 // num th
		} else if n >= 3585 && n <= 3679 {
			return 6 // other th
		} else {
			break
		}
	}

	return 0
}

func ClearHashTag(wstr string, cutcodepage bool) string {
	var s string
	for k1, wstr1 := range strings.Split(wstr, "\n") {
		if k1 != 0 {
			s += "\n"
		}
		for k2, wstr2 := range strings.Split(wstr1, "\r") {
			if k2 != 0 {
				s += "\r"
			}
			for k3, wstr3 := range strings.Split(wstr2, "\t") {
				if k3 != 0 {
					s += "\t"
				}
				for k4, wstr4 := range strings.Split(wstr3, " ") {
					if k4 != 0 {
						s += " "
					}
					if len(wstr4) > 1 && wstr4[0] == '#' {
						for _, wstr5 := range strings.Fields(wstr4) {
							s += wstr5
						}
					} else {
						if cutcodepage {
							s += ProcessChangeCodePage(wstr4)
						} else {
							s += wstr4

						}

					}

				}

			}
		}

	}
	return s
}

func GetType(str string, _ string) string {
	// if str == "." {
	// 	return ntype
	// }
	for i, char := range str {
		intc := int(char)
		if intc >= 3585 && intc <= 3675 {
			return "T"
		}
		if i > 0 {
			break
		}
	}
	return "E"
}
func ProcessChangeCodePage(strinput string) string {
	stroutput := ""
	for _, wstr := range strings.Split(strinput, "\u00a0") {

		if utf8.RuneCountInString(wstr) == len(wstr) {
			stroutput = stroutput + wstr + "\u00a0"
			continue
		}

		state := "U"
		statenow := "U"
		str := ""
		for _, ch := range strings.Split(wstr, "") {
			statenow = GetType(ch, state)

			if state == "T" && statenow == "E" {
				str = str + "\u00a0"
			}
			if state == "E" && statenow == "T" {
				str = str + "\u00a0"
			}
			str = str + ch
			state = statenow
		}

		_ = state
		_ = statenow

		stroutput = stroutput + str + "\u00a0"

	}
	return strings.TrimSpace(stroutput)

}

// false, true, true)
// recover bool, breakonly bool, cutcodepage bool
func BreakLineQuery(t *TriMain, line string) string {
	wstr := ""

	for i, l := range strings.Split(line, "\"") {
		_ = i
		for j, l2 := range strings.Split(l, " ") {
			_ = j
			if len(l2) == 0 {
				continue
			}
			str := BreakLineOrg(t, l2, true)
			str = ProcessChangeCodePage(str)
			if len(l2) != len(str) && i%2 == 0 {
				str = "\" " + str + " \""
			}
			wstr = wstr + " " + str
		}
		wstr = wstr + " \""
	}

	wstr = strings.TrimSuffix(wstr, " \"")
	// return wstr
	return ClearHashTag(wstr, true)
}

func BreakLine(t *TriMain, line string, recover bool, breakonly bool, cutcodepage bool) string {

	if breakonly {
		wstr := BreakLineOrg(t, line, recover)
		return ClearHashTag(wstr, cutcodepage)
	} else {
		var strb strings.Builder
		last := 1
		now := 0
		for i, data := range strings.Split(line, "") {

			now = Char2Int(data)
			if last != now {
				if now != 1 && last != 1 {
					strb.WriteString(string(rune(160)))
				}
			}
			strb.WriteString(data)
			last = now
			_ = i
		}
		return BreakLineOrg(t, strb.String(), recover)
	}

}

func BreakLineOrg(t *TriMain, line string, recover bool) string {
	// char[] InputStream = source.ToCharArray();
	//         TrieState state = new TrieState(dict.rootNode, InputStream, 0, InputStream.Length - 1);
	state := &TrieState{}
	state.New(t, line, 0, len(strings.Split(line, ""))-1)
	wasThai := false
	waitForPoint := false

	for state.StartIndex(0) <= state.StopIndex() {
		isThai := state.IsThai(wasThai)
		if isThai != wasThai {
			state.AddPoint(state.WordPointFromStart())
			waitForPoint = false
		}
		if isThai {
			state.WriteResult()
			state.ClearPoint()
			state.Home()
			prevTail := state.Tail()
			prevWord := state.CountWord(-1)
			noErr := state.TrieBreak()
			if state.Tail() > prevTail {

				if waitForPoint {
					if recover && (state.MaxWord() < prevWord+2) && !noErr {
						//if B2I(recover)&state.MaxWord() < prevWord+2&B2I(!noErr) {
						state.RemoveRange(prevTail+1, state.Tail()-prevTail)
						state.StartIndex(1)
					} else {
						a, b := state.WordPointFromStart()
						state.InsertPoint(prevTail+1, a, b)
					}
				} else {
					state.SetStartIndex(state.GetTailIndex())
				}
			} else {
				state.StartIndex(1)
			}

			waitForPoint = (state.Tail() == prevTail)

		} else {
			state.SetStartIndex(state.TextIndex() + 1)
		}

		wasThai = isThai
		state.SetTextIndex(state.StartIndex(0))
	}

	result := state.GetOutput()

	return result
}

func SplitLinesDict(s string) ([]string, string) {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	info := ""
	for sc.Scan() {
		word := sc.Text()
		if len(word) < 15 {
			continue
		}
		if strings.HasPrefix(word, "//") {
			info = info + word + "\n"
		} else {
			lines = append(lines, word)
		}

	}
	return lines, info
}
func LoadDictString(dicData string, data string) (string, *TriMain) {

	rootNode := &TriMain{}

	// t1 := time.Now()
	datadic, info := SplitLinesDict(dicData)
	_ = info
	if len(datadic) < 100 {
		return "", nil
	}
	rootNode.Load(datadic)
	// _ = rootNode
	// t2 := time.Now()
	// fmt.Println("Dict Time : ", (t2.Sub(t1)).String())
	// b := Break(rootNode, data)
	b := BreakLine(rootNode, data, true, true, true)

	return b, rootNode

}

func AssignNodeID(rootNode TriNode, nodeId int) int {
	starNode := rootNode
	if (starNode.Index == -1) && (!starNode.IsTop) {
		return nodeId
	}
	CurrNode := starNode
	for CurrNode.Index != -1 {
		if CurrNode.Index == -1 {
			CurrNode.Index = nodeId
			nodeId++
		}
		CurrNode = *CurrNode.NextNode
	}
	return nodeId
}

func Dequeue(queue []TriNode) ([]TriNode, TriNode) {
	if len(queue) == 0 {
		return queue, NewTriNode()
	}
	element := queue[0]
	return queue[1:], element
}
func SetLinkFlash(QueueNode []TriNode) error {
	if len(QueueNode) == 0 {
		return nil
	}
	var StartNode TriNode
	QueueNode, StartNode = Dequeue(QueueNode)
	CurrNode := StartNode
	for CurrNode.Index != -1 {
		if CheckTriNode(CurrNode, StartNode) {
			CurrNode.IsHead = true
		} else {
			CurrNode.IsTop = false
		}
		CurrNode.FromCount++
		CurrNode.IsReach = true
		CurrNode = *CurrNode.NextNode

	}
	CurrNode = StartNode
	TNode := CurrNode
	for CurrNode.Index != -1 {
		TNode = *CurrNode.LinkNode
		QueueNode = append(QueueNode, TNode)
		CurrNode = *CurrNode.NextNode
	}
	_ = TNode
	return nil
}
func CheckTriNode(Node1 TriNode, Node2 TriNode) bool {
	return true
}
func AddNode(word string, rootNode TriNode, idx int) (TriNode, int, error) {
	if rootNode.Index == -1 {
		rootNode, idx, _ = NewWordNode(word, idx, 0)
		return rootNode, idx, nil
	}
	CurrNode := rootNode
	PrevNode := NewTriNode()

	for i, ch := range word {

		for CurrNode.Index != -1 {
			if CurrNode.Charator == (string)(ch) {
				break
			}
			PrevNode = CurrNode
			CurrNode = *CurrNode.NextNode
		}

		if CurrNode.Index == -1 {
			PrevNode.IsLast = false
			*PrevNode.NextNode, _ = NewWord((string)(ch), i)
			PrevNode = *PrevNode.LinkNode
			if i+1 == len(word) {
				PrevNode.IsWord = true
				return rootNode, idx, nil
			}

			for n, ch2 := range word {
				if n < i {
					continue
				}

				LinkNode, _ := DupEnd(rootNode, word, n)
				PrevNode.LinkNode = &LinkNode

				_ = ch2
			}

		}

		if i+1 == len(word) {
			CurrNode.IsWord = true
		}
		if CurrNode.LinkNode == nil {

			*CurrNode.LinkNode, idx, _ = NewWordNode(word, i+1, idx)
			return rootNode, idx, nil

		} else {

			if CurrNode.LinkNode.FromCount > 1 {
				*CurrNode.LinkNode = Dupe(CurrNode)
			}
			if CurrNode.LinkNode.IsCommon {
				*CurrNode.LinkNode = Dupe(CurrNode)
			}
			CurrNode = *CurrNode.LinkNode

		}

	}

	return rootNode, idx, nil
}

func DupEnd(rootNode TriNode, word string, nstart int) (TriNode, error) {
	CurrNode, _ := Find(rootNode, word, nstart)
	return CurrNode, nil
}

func Find(rootNode TriNode, word string, nstart int) (TriNode, error) {
	return rootNode, nil
	// return FindNode( rootNode,word, nstart), nil
}

func ResetNode(rootNote TriNode) (TriNode, error) {
	node := rootNote
	for node.Index != -1 {
		node.FromCount = 0
		node.Index = -1
		node.IsHead = false
		node.IsTop = true
		node = *node.NextNode
	}
	node = rootNote
	for node.Index != -1 {
		ResetNode(*node.LinkNode)
		node = *node.NextNode

	}
	return rootNote, nil
}
func Dupe(node TriNode) TriNode {
	return NewTriNode()
}
func NewWord(word string, i int) (TriNode, error) {
	cnode, _, _ := NewWordNode(word, i, i)
	return cnode, nil
}
func NewWordNode(word string, istart int, idx int) (TriNode, int, error) {

	var newNode *TriNode
	var firstNode *TriNode
	var fromNode *TriNode
	node := NewTriNode()
	for i, ch := range word {
		if i < istart {
			continue
		}
		node = NewTriNode()
		newNode = &node

		idx++
		newNode.Charator = string(ch)
		if firstNode == nil {
			firstNode = newNode
		}
		if fromNode != nil {
			fromNode.LinkNode = newNode
			fromNode.IsWord = false
		}
		fromNode = newNode
	}

	return *firstNode, idx, nil

}
