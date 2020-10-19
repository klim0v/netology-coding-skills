package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var port = flag.Int("port", 3000, "HTTP port to listen on")

func init() {
	flag.Parse()
}

func main() {
	sigs := make(chan os.Signal, 1)
	fail := make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	server := fasthttp.Server{Handler: router}
	go func() { fail <- server.ListenAndServe(":" + strconv.Itoa(*port)) }()

	select {
	case <-sigs:
		log.Println("Shutdown signal received.")
		err := server.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
	case err := <-fail:
		if err != nil {
			log.Fatal(err)
		}
	}
}

func router(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/api/slow":
		timeoutMiddleware(slowHandler, 5*time.Second)(ctx)
	default:
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
	}
}

func timeoutMiddleware(handler fasthttp.RequestHandler, duration time.Duration) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		done := make(chan struct{})
		go func() {
			handler(ctx)
			close(done)
		}()

		select {
		case <-time.After(duration):
			var resp fasthttp.Response
			resp.SetStatusCode(400)
			resp.Header.SetContentType("application/json")
			resp.SetBodyString("{\"error\":\"timeout too long\"}")
			ctx.TimeoutErrorWithResponse(&resp)
			return
		case <-done:
			return
		}
	}
}

type requestBody struct {
	Timeout int `json:"timeout"`
}

func slowHandler(ctx *fasthttp.RequestCtx) {
	var body requestBody
	if err := json.Unmarshal(ctx.Request.Body(), &body); err != nil {
		ctx.SetStatusCode(400)
		ctx.SetContentType("application/json")
		ctx.SetBodyString(fmt.Sprintf("{\"status\":\"%s\"}", err.Error()))
		return
	}

	timer := time.After(time.Duration(body.Timeout) * time.Millisecond)
	for {
		select {
		case <-timer:
			ctx.SuccessString("application/json", "{\"status\":\"ok\"}")
			return
		default:
			if ctx.LastTimeoutErrorResponse() == nil {
				continue
			}
			return
		}
	}
}
