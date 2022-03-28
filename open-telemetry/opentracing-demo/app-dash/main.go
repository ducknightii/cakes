package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"sourcegraph.com/sourcegraph/appdash"
	appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"
	"sourcegraph.com/sourcegraph/appdash/traceapp"
)

var host = "http://localhost:8088"

func main() {
	store := appdash.NewMemoryStore()

	// Listen on any available TCP port locally.
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		panic(err)
	}
	collectorPort := l.Addr().(*net.TCPAddr).Port
	collectorAddr := fmt.Sprintf(":%d", collectorPort)

	// Start an Appdash collection server that will listen for spans and
	// annotations and add them to the local collector (stored in-memory).
	cs := appdash.NewServer(l, appdash.NewLocalCollector(store))
	go cs.Start()

	// Print the URL at which the web UI will be running.
	appdashPort := 8700
	appdashURLStr := fmt.Sprintf("http://localhost:%d", appdashPort)
	appdashURL, err := url.Parse(appdashURLStr)
	if err != nil {
		fmt.Printf("Error parsing %s: %s\n", appdashURLStr, err)
		return
	}
	fmt.Printf("To see your traces, go to %s/traces\n", appdashURL)

	// Start the web UI in a separate goroutine.
	tapp, err := traceapp.New(nil, appdashURL)
	if err != nil {
		panic(err)
	}
	tapp.Store = store
	tapp.Queryer = store
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", appdashPort), tapp)
		panic(err)
	}()

	tracer := appdashot.NewTracer(appdash.NewRemoteCollector(collectorAddr))
	opentracing.InitGlobalTracer(tracer)

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/home", homeHandler)
	mux.HandleFunc("/async", serviceHandler)
	mux.HandleFunc("/service", serviceHandler)
	mux.HandleFunc("/db", dbHandler)
	fmt.Printf("Go to %s/home to start a request!\n", host)
	fmt.Println(http.ListenAndServe(":8088", mux))
}

// Acts as our index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<a href="/home"> Click here to start a request </a>`))
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	span := opentracing.StartSpan("/home")
	defer span.Finish()

	w.Write([]byte("Request started\n"))
	// Since we have to inject our span into the HTTP headers, we create a request
	asyncReq, _ := http.NewRequest("GET", host+"/async", nil)
	// Inject the span context into the header
	carrier := opentracing.HTTPHeadersCarrier(asyncReq.Header)
	err := span.Tracer().Inject(span.Context(),
		opentracing.HTTPHeaders,
		carrier)
	if err != nil {
		fmt.Printf("Could not inject span context into header: %v\n", err)
		return
	}
	// ep: http.Header{"Ot-Tracer-Sampled":[]string{"true"}, "Ot-Tracer-Spanid":[]string{"330039d5a1cae2d3"}, "Ot-Tracer-Traceid":[]string{"78187e424b2180da"}}
	fmt.Printf("Header: %#v\n", asyncReq.Header)
	go func() {
		if _, err := http.DefaultClient.Do(asyncReq); err != nil {
			span.SetTag("error", true)
			span.LogFields(log.String("detail", fmt.Sprintf("GET /async error:%v", err)))
		}
	}()

	serviceReq, _ := http.NewRequest("GET", host+"/service", nil)
	// Inject the span context into the header
	sCarrier := opentracing.HTTPHeadersCarrier(serviceReq.Header)
	err = span.Tracer().Inject(span.Context(),
		opentracing.HTTPHeaders,
		sCarrier)

	_, err = http.DefaultClient.Do(serviceReq)
	if err != nil {
		ext.Error.Set(span, true)                                                      // Tag the span as errored
		span.LogFields(log.String("detail", fmt.Sprintf("GET service error:%v", err))) // Log the error
	}

	time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
	w.Write([]byte("Request done!"))
}

// Mocks a service endpoint that makes a DB call
func serviceHandler(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	opName := r.URL.Path
	// Attempt to join a trace by getting trace context from the headers.
	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		// If for whatever reason we can't join, go ahead an start a new root span.
		sp = opentracing.StartSpan(opName)
	} else {
		sp = opentracing.StartSpan(opName, opentracing.ChildOf(wireContext))
	}
	defer sp.Finish()
	// ...
	http.Get(host + "/db")
	time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
	// ...
}

// Mocks a DB call
func dbHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
	// here would be the actual call to a DB.
}
