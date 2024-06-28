package main

import (
	_ "embed"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"

	wbGo "github.com/Kard34/search-engine/dataxet/iq-wordbreak"
	"github.com/Kard34/search-engine/ftime"
	"github.com/Kard34/search-engine/qp"
	search_engine "github.com/Kard34/search-engine/search-engine"
	"github.com/gofiber/fiber/v2/middleware/cors"

	_ "github.com/mattn/go-sqlite3"
)

type inputData struct {
	Query    string `json: "query"`
	Limit    int    `json: "limit"`
	Offset   int    `json: "offset"`
	Filename string `json: filename`
}

var (
	Cache = cache.New(5*time.Minute, 5*time.Minute)

	DictDataIQ string = strings.TrimSpace(dictdataiq)
	//go:embed godict.txt
	dictdataiq string
)

func main() {
	log.Print("Starting the service.")
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST",
		AllowHeaders: "Content-Type",
	}))
	app.Post("/search", searchData)
	wordbreak("")
	log.Print("The service is ready to listen and serve.")
	app.Listen(":8080")
}

func searchData(c *fiber.Ctx) error {
	var (
		di        *wbGo.TriMain
		wordList  []string
		startTime ftime.CTime
		endTime   ftime.CTime
	)
	inputdata := new(inputData)
	if err := c.BodyParser(inputdata); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if dic, found := Cache.Get("dictIQ"); found {
		di = dic.(*wbGo.TriMain)
	} else {
		di = nil
	}

	if t := c.Query("start_time"); t != "" {
		startTime.Parse(t)
	}
	if t := c.Query("end_time"); t != "" {
		endTime.Parse(t)
	}

	fx, x, err := qp.Parse(di, (*inputdata).Query)
	fmt.Println(startTime, endTime)
	fmt.Println(fx, x, err)

	for _, item := range fx {
		if item.Lt == -1 && item.Rt == -1 {
			wordList = append(wordList, "'"+item.Val+"'")
		}
	}
	var wordQuery string
	wordQuery += "(" + strings.Join(wordList, ",") + ")"

	Result := search_engine.SearchFile(wordQuery, fx, (*inputdata).Limit, (*inputdata).Offset, startTime, endTime, (*inputdata).Filename)
	fmt.Println("RESULT : ", len(Result))
	return c.JSON(Result)
}

func wordbreak(str string) (result string) {
	var di *wbGo.TriMain
	if dic, found := Cache.Get("dictIQ"); found {
		di = dic.(*wbGo.TriMain)
	} else {
		_, di = wbGo.LoadDictString(DictDataIQ, "")
		Cache.Add("dictIQ", di, cache.NoExpiration)
	}
	breakonly := false
	cutbetweencodepage := false
	result = wbGo.BreakLine(di, str, true, breakonly, cutbetweencodepage)
	return
}
