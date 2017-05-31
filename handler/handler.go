package handler

import (
	"fmt"
	"encoding/json"
	"log"
	"net/url"
	"github.com/valyala/fasthttp"
	"../storage"
)

func Handler(ctx *fasthttp.RequestCtx) {}


func CreateUrlHandler(h fasthttp.RequestHandler, s *storage.Storage, c Configuration) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		if !CheckAuthentication(c, ctx) {
			ctx.SetStatusCode(401)
			fmt.Fprintf(ctx, "Unauthorized")
		} else {
			var urlAddress string
			if ctx.Request.Header.IsPost() {
				params := ctx.Request.PostArgs()
				urlAddress = string(params.Peek("url"))
			} else {
				query := ctx.URI().QueryArgs()
				urlAddress = string(query.Peek("url"))
			}
			_, err := url.ParseRequestURI(urlAddress)
			if err != nil {
				ctx.SetStatusCode(500)
				fmt.Fprintf(ctx, "Invalid URL")
			} else {
				url := storage.CreateUrl(s, urlAddress)
				fmt.Fprintf(ctx, "%s", ConvertUrlToJson(url))
				ctx.SetContentType("application/json")
			}
		}
		h(ctx)
	})
}

func GetUrlHandler(h fasthttp.RequestHandler, s *storage.Storage, c Configuration) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		if !CheckAuthentication(c, ctx) {
			ctx.SetStatusCode(401)
		} else {
			uid := ctx.UserValue("uid").(string)
			url := storage.GetUrl(s, uid)
			if url.Url != "" {
				fmt.Fprintf(ctx, "%s", ConvertUrlToJson(*url))
				ctx.SetContentType("application/json")
			} else {
				ctx.SetStatusCode(404)
				fmt.Fprintf(ctx, "Not Found")
			}
		}
		h(ctx)
	})
}

func RedirectHandler(h fasthttp.RequestHandler, s *storage.Storage) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		uid := ctx.UserValue("uid").(string)
		url := storage.GetUrl(s, uid)
		if url.Url != "" {
			storage.UpdateStatistics(s, uid)
			ctx.Redirect(url.Url, 301)
		} else {
			ctx.SetStatusCode(404)
			fmt.Fprintf(ctx, "Not Found")
		}
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