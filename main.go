package main

import (
	"github.com/rs/xhandler"
	"github.com/rs/xlog"
	"github.com/rs/xmux"
	"golang.org/x/net/context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var response = []byte("pong")

const port = ":7999"

func main() {
	c := xhandler.Chain{}

	host, _ := os.Hostname()
	conf := xlog.Config{
		// Log info level and higher
		Level: xlog.LevelInfo,
		// Set some global env fields
		Fields: xlog.F{
			"role": "my-shitty-service",
			"host": host,
		},
		// Output everything on console
		Output: xlog.NewOutputChannel(xlog.NewConsoleOutput()),
	}

	// Add close notifier handler so context is cancelled when the client closes
	// the connection
	c.UseC(xhandler.CloseHandler)

	// Install the logger handler
	c.UseC(xlog.NewHandler(conf))

	// Add timeout handler (HAHA)
	//c.UseC(xhandler.TimeoutHandler(2 * time.Second))

	// Install some provided extra handler to set some request's context fields.
	// Thanks to those handler, all our logs will come with some pre-populated fields.
	c.UseC(xlog.MethodHandler("method"))
	c.UseC(xlog.URLHandler("url"))
	c.UseC(xlog.RemoteAddrHandler("ip"))
	c.UseC(xlog.UserAgentHandler("user_agent"))
	c.UseC(xlog.RefererHandler("referer"))
	c.UseC(xlog.RequestIDHandler("req_id", "Request-Id"))

	mux := xmux.New()

	mux.GET("/gamble", xhandler.HandlerFuncC(shittyHandler))

	log.Printf("Listening on http://localhost%v/gamble", port)
	http.ListenAndServe(port, c.Handler(mux))

}

func shittyHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	i := rand.Intn(1000)

	if i < 5 {
		eternalRequest(ctx, rw, r)
		return
	}
	if i < 30 {
		longRequest(ctx, rw, r)
		return

	}
	if i < 50 {
		emptyBody(ctx, rw, r)
		return
	}
	normalResponse(ctx, rw, r)
}

func normalResponse(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	l := xlog.FromContext(ctx)
	l.Info("normal response")
	rw.Write(response)
}

func longRequest(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	l := xlog.FromContext(ctx)

	minutes := time.Duration(rand.Intn(60) + 30) * time.Minute
	l.Infof("longRequest %v", minutes)

	time.Sleep(minutes)
	rw.Write(response)
}

func eternalRequest(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	l := xlog.FromContext(ctx)

	hours := time.Duration(rand.Intn(24)) * time.Hour
	l.Infof("eternalRequest %v", hours)

	time.Sleep(hours)
	rw.Write(response)
}

func emptyBody(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	l := xlog.FromContext(ctx)

	l.Info("emptyBody")

	rw.Write(nil)
}
