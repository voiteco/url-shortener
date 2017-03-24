package main

import (
	"fmt"
	"flag"
	"log"
	"runtime"
	"math/rand"
	"time"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"./handler"
	"./storage"
)

var (
	port   *uint
	redis  *string
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	port = flag.Uint("port", 8080, "Bind port")
	redis = flag.String("redis", "127.0.0.1:6379", "Redis address")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	s := storage.InitStorage(*redis)
	c := handler.Configuration{}
	c.ParseConfig()

	router := fasthttprouter.New()
	router.GET("/create", handler.CreateUrlHandler(handler.Handler, s, c))
	router.GET("/get/:uid", handler.GetUrlHandler(handler.Handler, s, c))
	router.GET("/u/:uid", handler.RedirectHandler(handler.Handler, s))

	log.Fatal(fasthttp.ListenAndServe(fmt.Sprintf(":%d", *port), router.Handler))
}
