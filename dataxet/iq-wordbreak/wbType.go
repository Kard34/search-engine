package iq_wordbreak

func WBType(char int) int {
	if char < 0 {
		return 65536
	}
	if char < 46 {
		return 2
	}
	if char < 47 {
		return 4098
	}
	if char < 65 {
		return 2
	}
	if char < 91 {
		return 8194
	}
	if char < 97 {
		return 2
	}
	if char < 123 {
		return 8194
	}
	if char < 128 {
		return 2
	}
	if char < 3585 {
		return 65536
	}
	if char < 3631 {
		return 6
	}
	if char < 3632 {
		return 2
	}
	if char < 3633 {
		return 18
	}
	if char < 3634 {
		return 33
	}
	if char < 3636 {
		return 18
	}
	if char < 3640 {
		return 33
	}
	if char < 3642 {
		return 67
	}
	if char < 3643 {
		return 3
	}
	if char < 3648 {
		return 2
	}
	if char < 3653 {
		return 10
	}
	if char < 3654 {
		return 18
	}
	if char < 3655 {
		return 4
	}
	if char < 3656 {
		return 1281
	}
	if char < 3660 {
		return 1664
	}
	if char < 3661 {
		return 1792
	}
	if char < 3663 {
		return 1
	}
	if char < 3664 {
		return 2
	}
	if char < 3674 {
		return 2050
	}
	if char < 3680 {
		return 2
	}
	return 65536
}
func WBTypeIsType(str string, check CharFlag) bool {

	for _, char := range str {

		intc := int(char)
		wbb := WBType(intc)
		ff := CharFlag(wbb)
		n := CharFlag(ff & check)
		result1 := ff & check
		_ = result1
		if n > 0 {
			return true
		} else {
			break
		}
	}
	return false
}

type PointType int

const (
	PointSentinel  PointType = 1
	PointCandidate PointType = 2
	PointWord      PointType = 3
)
