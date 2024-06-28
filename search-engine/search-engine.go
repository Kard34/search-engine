package search_engine

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Kard34/search-engine/ftime"
	"github.com/Kard34/search-engine/qp"

	_ "github.com/mattn/go-sqlite3"
)

type Treenode struct {
	Value string
	Left  *Treenode
	Right *Treenode
}

type ChunkData struct {
	Index         int
	Position      int
	Allocate      int
	CountDocument int
	StartPosition int
	CountPosition int
}

type responseData struct {
	Date     string `json: "displaytime"`
	DocID    string `json: "id"`
	Headline string `json: "headline'`
}

var (
	Path     = "./Index/"
	Filename = "20240129"

	Fidx *os.File
	Db   *sql.DB

	CkData map[string]ChunkData
)

func SearchFile(query string, node_list []qp.FlatNode, limit int, offset int, timex, timey ftime.CTime, filename string) (listdata []responseData) {
	db, err := sql.Open("sqlite3", filename+".sqlite")
	if err != nil {
		fmt.Println(err)
		return
	}
	Db = db
	defer Db.Close()

	fidx, err := os.Open(filename + ".idx")
	if err != nil {
		fmt.Println(err)
		return
	}
	Fidx = fidx
	defer Fidx.Close()

	Load(query)
	Root := MakeTree(node_list)
	listdata = Search(Root, limit, offset, timex, timey)

	return
}

func Search(tree *Treenode, limit int, offset int, timex, timey ftime.CTime) (listdata []responseData) {
	Timex := ParseStr(TimeToStr(timex))
	Timey := ParseStr(TimeToStr(timey))
	Buffx := docInvert(Timex.Year(), int(Timex.Month()), Timex.Day(), Timex.Hour(), 0)
	Buffx = append(Buffx, []byte{0, 0, 0}...)
	t := 0
	if Timey.Minute() > 0 || Timey.Second() > 0 {
		t = 1
	}
	Buffy := docInvert(Timey.Year(), int(Timey.Month()), Timey.Day(), Timey.Hour()+t, 0)
	Buffy = append(Buffy, []byte{0, 0, 0}...)
	ID_List := SearchData(tree, Buffx, Buffy)
	placeholders := make([]string, len(ID_List))
	args := make([]interface{}, len(ID_List)+4)
	// xxx := make([]string, len(ID_List))
	for i, id := range ID_List {
		placeholders[i] = "?"
		args[i] = id
		// xxx[i] = "'" + strconv.Itoa(int(id)) + "'"
	}
	args[len(args)-4] = timex.UnixMilli()
	args[len(args)-3] = timey.UnixMilli()
	args[len(args)-2] = limit
	args[len(args)-1] = offset
	x := `
	SELECT DOCID, TIME64, HEADLINE 
	FROM HDL 
	WHERE INVDOCID IN` + `(` + strings.Join(placeholders, ",") + `)
	AND TIME64 BETWEEN ? AND ?
	ORDER BY TIME64
	LIMIT ? OFFSET ?`
	rows, err := Db.Query(x, args...)
	checkERROR(err)
	defer rows.Close()
	for rows.Next() {
		var DOCID string
		var TIME64 int64
		var HEADLINE string
		err := rows.Scan(&DOCID, &TIME64, &HEADLINE)
		checkERROR(err)

		DisplayTime := time.UnixMilli(int64(TIME64))
		listdata = append(listdata, responseData{DisplayTime.UTC().Format("2006-01-02T15:04:05"), DOCID, HEADLINE})

		if len(listdata) >= limit {
			break
		}
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}
	// sort.Slice(listdata, func(i, j int) bool {
	// 	tx, _ := time.Parse("2006-01-02T15:04:05", listdata[i].Date)
	// 	ty, _ := time.Parse("2006-01-02T15:04:05", listdata[j].Date)
	// 	return tx.Before(ty)
	// })
	return
}

func SearchData(tree *Treenode, buffx, buffy []byte) (invdocid_list []uint64) {
	Chunkdata, Buff := SearchMatching(tree, buffx, buffy)
	for i := 0; i < Chunkdata.CountDocument; i++ {
		Buff8 := make([]byte, 0)
		Buff8 = append(Buff8, Buff[i*10:(i*10)+5]...)
		Buff8 = append(Buff8, []byte{0, 0, 0}...)
		INVDOCID := binary.LittleEndian.Uint64(Buff8)
		invdocid_list = append(invdocid_list, INVDOCID)
	}
	return
}

func SearchMatching(tree *Treenode, buffx, buffy []byte) (chunkdata ChunkData, buff []byte) {
	var Chunk1 ChunkData
	var Buff1 []byte
	var Chunk2 ChunkData
	var Buff2 []byte
	if tree.Left != nil {
		Chunk1, Buff1 = SearchMatching(tree.Left, buffx, buffy)
	}
	if tree.Right != nil {
		Chunk2, Buff2 = SearchMatching(tree.Right, buffx, buffy)
	}
	if tree.Left == nil && tree.Right == nil {
		chunkdata, buff = LoadWord(tree.Value, buffx, buffy)
	} else {
		chunkdata, buff = Match(Chunk1, Chunk2, Buff1, Buff2, tree.Value)
	}
	return
}

func LoadWord(word string, buffx, buffy []byte) (chunkdata ChunkData, buff []byte) {
	Lpos1, Found1 := BinaryChunkBuff(buffx, word)
	Lpos2, Found2 := BinaryChunkBuff(buffy, word)
	_ = Found1
	if Found2 {
		Lpos2++
	}
	Allocate := 16
	CountDocument := 0
	CountPosition := 0
	StartPoint := int32(CkData[word].Position) + 16
	INVDOCID_LIST := make([]byte, 0)
	PositionPoint := make([]int32, 2)

	for i := Lpos1; i < Lpos2; i++ {
		Buff1 := make([]byte, 10)
		Fidx.Seek(int64(StartPoint+int32(i*10)), io.SeekStart)
		Fidx.Read(Buff1)
		INVID := Buff1[0:5]
		INDEX := []byte{byte(CountPosition & 255), byte((CountPosition >> 8) & 255), byte((CountPosition >> 16) & 255)}
		LENGTH := Buff1[8:10]
		LengthValue := int(binary.LittleEndian.Uint16(LENGTH))
		Buff10 := make([]byte, 0)
		Buff10 = append(Buff10, INVID...)
		Buff10 = append(Buff10, INDEX...)
		Buff10 = append(Buff10, LENGTH...)
		Allocate += 10
		CountDocument++
		CountPosition += LengthValue
		INVDOCID_LIST = append(INVDOCID_LIST, Buff10...)
		if i == Lpos1 {
			Temp := Buff1[5:8]
			Temp = append(Temp, []byte{0}...)
			PositionPoint[0] = int32(binary.LittleEndian.Uint32(Temp))
		} else if i == Lpos2-1 {
			Temp := Buff1[5:8]
			Temp = append(Temp, []byte{0}...)
			PositionPoint[1] = int32(binary.LittleEndian.Uint32(Temp) + uint32(LengthValue))
		}
	}
	Buff := make([]byte, (PositionPoint[1]-PositionPoint[0])*2)
	Fidx.Seek(int64(StartPoint+int32(CkData[word].StartPosition)+(PositionPoint[0]*2)), io.SeekStart)
	Fidx.Read(Buff)
	buff = append(buff, INVDOCID_LIST...)
	buff = append(buff, Buff...)
	chunkdata = ChunkData{CkData[word].Index, CkData[word].Position, Allocate + (CountPosition * 2), CountDocument, Allocate - 16, CountPosition}
	return
}

func BinaryChunkBuff(buffsearch []byte, word string) (lpos int, found bool) {
	LposLo := 0
	LposHi := CkData[word].CountDocument - 1
	CompareResult := -1
	lpos = 0
	StartPoint := int32(CkData[word].Position) + 16
	for LposLo <= LposHi {
		lpos = (LposLo + LposHi) / 2
		Buff := make([]byte, 10)
		Fidx.Seek(int64(StartPoint+int32(lpos*10)), io.SeekStart)
		Fidx.Read(Buff)
		CompareResult = bytes.Compare(buffsearch[0:5], Buff[0:5])
		if CompareResult < 0 {
			LposHi = lpos - 1
		} else if CompareResult > 0 {
			LposLo = lpos + 1
		} else {
			break
		}
	}
	if CompareResult > 0 {
		lpos += 1
	}
	found = CompareResult == 0
	return
}

func Match(cho1, cho2 ChunkData, buffw1, buffw2 []byte, op string) (cho ChunkData, buff []byte) {
	idx := 0
	jdx := 0
	cho = ChunkData{-1, -1, 0, 0, 0, 0}
	buffdoc := make([]byte, 0)
	buffpos := make([]byte, 0)
	buff0 := make([]byte, 5)
	nCompareResult := -1
	start_doc_pos := 0
	len_doc_post := 0
	buff3 := make([]byte, 4)
	buff2 := make([]byte, 2)
	for idx < cho1.CountDocument && jdx < cho2.CountDocument {
		b1 := buffw1[idx*10 : (idx*10)+10]
		b2 := buffw2[jdx*10 : (jdx*10)+10]
		nCompareResult = bytes.Compare(b1[0:5], b2[0:5])
		if nCompareResult < 0 {
			if op == "or" {
				buffdoc = append(buffdoc, b1[0:5]...)
				buffdoc = append(buffdoc, buff0...)
			}
			idx++
		} else if nCompareResult > 0 {
			if op == "or" {
				buffdoc = append(buffdoc, b2[0:5]...)
				buffdoc = append(buffdoc, buff0...)
			}
			jdx++
		} else {
			if op == "or" {
				buffdoc = append(buffdoc, b1[0:5]...)
				buffdoc = append(buffdoc, buff0...)
			} else if op == "and" {
				buffdoc = append(buffdoc, b1[0:5]...)
				buffdoc = append(buffdoc, buff0...)
			} else {
				diff := 3
				if op == "phrase2" {
					diff = 2
				}

				st1, len1 := invposition(b1)
				st2, len2 := invposition(b2)
				if cho1.StartPosition+((st1+len1)*2) > len(buffw1) {
					fmt.Println("error")
				}
				if cho2.StartPosition+((st2+len2)*2) > len(buffw2) {
					fmt.Println("error")
				}

				bo1 := buffw1[cho1.StartPosition+(st1*2) : cho1.StartPosition+((st1+len1)*2)]
				bo2 := buffw2[cho2.StartPosition+(st2*2) : cho2.StartPosition+((st2+len2)*2)]

				pos := comparepharse(bo1, bo2, diff)

				if len(pos) > 0 {
					buffpos = append(buffpos, pos...)
					buffdoc = append(buffdoc, b1[0:5]...)

					binary.LittleEndian.PutUint32(buff3, uint32(start_doc_pos))
					buffdoc = append(buffdoc, buff3[0:3]...)
					binary.LittleEndian.PutUint16(buff2, uint16(len(pos)/2))
					buffdoc = append(buffdoc, buff2...)
					start_check := int(binary.LittleEndian.Uint32(buff3))
					len_check := int(binary.LittleEndian.Uint16(buff2))
					if start_check != start_doc_pos || len_check != len(pos)/2 {
						fmt.Print("error")
					}
					_ = len_doc_post
					cho.CountDocument++
					len_pos := len(pos) / 2
					cho.CountDocument += len_pos
					start_doc_pos += len_pos
				}
			}
			idx++
			jdx++
		}
	}
	buff = append(buff, buffdoc...)
	buff = append(buff, buffpos...)
	cho.Allocate = len(buff)
	cho.CountDocument = len(buffdoc) / 10
	cho.CountPosition = len(buffpos) / 2
	cho.StartPosition = len(buffdoc)
	return
}

func invposition(buff []byte) (start, len int) {
	b3 := make([]byte, 0)
	b3 = append(b3, buff[5:8]...)
	b3 = append(b3, 0)
	start = int(binary.LittleEndian.Uint32(b3))
	len = int(binary.LittleEndian.Uint16(buff[8:10]))
	return
}

func comparepharse(bo1 []byte, bo2 []byte, diff int) (buff []byte) {
	idx := 0
	jdx := 0
	buff = make([]byte, 0)
	if len(bo1) == 0 || len(bo2) == 0 {
		return
	}
	idx = 0
	jdx = 0
	for idx < len(bo1) && jdx < len(bo2) {
		vali := binary.LittleEndian.Uint16(bo1[idx : idx+2])
		valj := binary.LittleEndian.Uint16(bo2[jdx : jdx+2])
		if vali+uint16(diff) > valj {
			jdx += 2
		} else if vali+uint16(diff) < valj {
			idx += 2
		} else {
			buff = append(buff, bo1[idx:idx+2]...)
			idx += 2
			jdx += 2
		}
	}
	return
}

func Load(query string) {
	x := "SELECT * FROM IDX WHERE WORD IN " + query
	rows, err := Db.Query(x)
	checkERROR(err)

	CkData = map[string]ChunkData{}

	for rows.Next() {
		var WORD string
		var INDEX int
		var POSITION int
		var ALLOCATE int
		var COUNTDOCMENT int
		var STARTPOSITION int
		var COUNTPOSITION int

		err := rows.Scan(&WORD, &INDEX, &POSITION, &ALLOCATE, &COUNTDOCMENT, &STARTPOSITION, &COUNTPOSITION)
		checkERROR(err)
		CkData[WORD] = ChunkData{INDEX, POSITION, ALLOCATE, COUNTDOCMENT, STARTPOSITION, COUNTPOSITION}
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
}

func MakeTree(query []qp.FlatNode) (root *Treenode) {
	Data := map[int]qp.FlatNode{}
	Found := map[int]int{}
	for _, i := range query {
		Data[i.Idx] = i
		Found[i.Idx]++
		Found[i.Idx]--
		Found[i.Lt]++
		Found[i.Rt]++
	}
	Head := -1
	for x, y := range Found {
		if y == 0 {
			Head = x
		}
	}
	root = maketree(Data, Head)
	return
}

func maketree(data map[int]qp.FlatNode, head int) (root *Treenode) {
	root = &Treenode{}
	root.Value = data[head].Val
	if data[head].Lt != -1 {
		root.Left = maketree(data, data[head].Lt)
	}
	if data[head].Rt != -1 {
		root.Right = maketree(data, data[head].Rt)
	}
	return
}

func checkERROR(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func ParseStr(str string) time.Time {
	Prased, err := time.Parse("2006-01-02 15:04:05 -0700 MST", str)
	if err != nil {
		fmt.Println("Error:", err)
		return time.Time{}
	}
	return Prased
}

func TimeToStr(time ftime.CTime) (str string) {
	str = time.Format("2006-01-02 15:04:05 -0700 MST")
	return
}

func docInvert(year, month, day, hour, running int) (buff []byte) {
	year -= 1950
	month -= 1
	day -= 1
	Value := year*12*31*24 + month*31*24 + day*24 + hour
	buff = make([]byte, 5)
	buff[0] = byte(Value >> 12)
	buff[1] = byte(Value >> 4)
	buff[2] = byte(Value<<4) | byte(running>>8)
	buff[3] = byte(running)
	buff[4] = byte(running >> 8)
	return
}
