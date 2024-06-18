package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	wbGo "github.com/Kard34/search-engine/dataxet/iq-wordbreak"
	"github.com/Kard34/search-engine/qp"
	"github.com/patrickmn/go-cache"
)

var (
	Cache = cache.New(5*time.Minute, 5*time.Minute)
)

func main() {
	var (
		di *wbGo.TriMain
	)
	app := fiber.New()

	_ = app

	if dic, found := Cache.Get("dictIQ"); found {
		di = dic.(*wbGo.TriMain)
	} else {
		di = nil
	}

	fx, x, err := qp.Parse(di, "\"ค่าเงินบาท\" and \"วันนี้\" and \"ธนาคาร\"")
	fmt.Println(fx, x, err)
	// search_engine.Search(query) น่าจะเป็นอะไรประมาณนี้

	// app.Listen(":8080")
}
