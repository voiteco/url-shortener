package handler

import (
	"fmt"
	"encoding/json"
	"log"
	"github.com/valyala/fasthttp"
	"../storage"
)

func Handler(ctx *fasthttp.RequestCtx) {}


func CreateUrlHandler(h fasthttp.RequestHandler, s *storage.Storage, c Configuration) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		if !CheckAuthentication(c, ctx) {
			ctx.SetStatusCode(403)
		} else {
			query := ctx.URI().QueryArgs()
			urlAddress := string(query.Peek("url"))
			url := storage.CreateUrl(s, urlAddress)
			fmt.Fprintf(ctx, "%s", ConvertUrlToJson(url))
			ctx.SetContentType("application/json")
		}
		h(ctx)
	})
}

func GetUrlHandler(h fasthttp.RequestHandler, s *storage.Storage, c Configuration) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		if !CheckAuthentication(c, ctx) {
			ctx.SetStatusCode(403)
		} else {
			uid := ctx.UserValue("uid").(string)
			url := storage.GetUrl(s, uid)
			fmt.Fprintf(ctx, "%s", ConvertUrlToJson(*url))
			ctx.SetContentType("application/json")
		}
		h(ctx)
	})
}

func RedirectHandler(h fasthttp.RequestHandler, s *storage.Storage) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		uid := ctx.UserValue("uid").(string)
		url := storage.GetUrl(s, uid)
		storage.UpdateStatistics(s, uid)
		ctx.Redirect(url.Url, 302)
		h(ctx)
	})
}

func CheckAuthentication(c Configuration, ctx *fasthttp.RequestCtx) bool {
	if c.Authentication {
		token := ctx.Request.Header.Peek(c.AuthenticationParameter)
		return string(token) == c.AuthenticationToken
	}
	return true
}

func ConvertUrlToJson(url storage.Url) string {
	j, err := json.Marshal(url)
	if err != nil {
		log.Println(err)
	}
	return string(j);
}