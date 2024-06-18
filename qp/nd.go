package qp

type ND interface {
	val() string
	getOp() string
	hasNext() bool
	hasLeft() bool
	getLeft() ND
	getNext() ND
	isLeaf() bool
	breakable() bool
}

func (n Node) getOp() string {
	return "or"
}
func (n NodeAnd) getOp() string {
	return "and"
}
func (n NodeAndNot) getOp() string {
	return "and not"
}

func (n Node) hasLeft() bool {
	return n.Left != nil
}
func (n NodeAnd) hasLeft() bool {
	return n.Left != nil
}
func (n NodeAndNot) hasLeft() bool {
	return n.Val.Sub != nil
}

func (n Node) hasNext() bool {
	return n.Next != nil
}
func (n NodeAnd) hasNext() bool {
	return n.Next != nil
}
func (n NodeAndNot) hasNext() bool {
	return n.Next != nil
}

func (n Node) getLeft() ND {
	return n.Left
}
func (n NodeAnd) getLeft() ND {
	return n.Left
}
func (n NodeAndNot) getLeft() ND {
	return n.Val.Sub
}

func (n Node) getNext() ND {
	return n.Next
}
func (n NodeAnd) getNext() ND {
	return n.Next
}
func (n NodeAndNot) getNext() ND {
	return n.Next
}
func (n Node) isLeaf() bool {
	return false
}
func (n NodeAnd) isLeaf() bool {
	return false
}
func (n NodeAndNot) isLeaf() bool {
	return n.Val.Sub == nil
}

func (n Node) val() string {
	return n.Left.Left.Val.String()
}
func (n NodeAnd) val() string {
	return n.Left.Val.String()
}
func (n NodeAndNot) val() string {
	return n.Val.String()
}

func (n Node) breakable() bool {
	return false //never called, come after isLeaf
}
func (n NodeAnd) breakable() bool {
	return false //never called, come after isLeaf
}
func (n NodeAndNot) breakable() bool {
	return n.Val.Phrase != nil || (n.Val.Tag == nil && n.Val.Str != nil && string((*n.Val.Str)[0]) != "%")
}
