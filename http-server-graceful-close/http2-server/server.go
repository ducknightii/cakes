package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

// 证书生成: https://www.iszy.cc/posts/u62i8r/

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz)

	server := &http.Server{
		Addr:    ":8000",
		Handler: mux,
		ConnState: func(conn net.Conn, state http.ConnState) {
			fmt.Printf("conn: %s, state: %s\n", conn.RemoteAddr(), state.String())
		},
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		s := <-c
		fmt.Printf("System Single call: %+v\n", s)
		cancel()
	}()

	go func() {
		fmt.Println("Server http2 :8000 start")
		// When Shutdown is called, Serve, ListenAndServe, and
		// ListenAndServeTLS immediately return ErrServerClosed. Make sure the
		// program doesn't exit and waits instead for Shutdown to return.
		err := server.ListenAndServeTLS("server.crt", "server.key")
		fmt.Printf("Server close: %+v\n", err)
	}()

	// 优雅关闭
	// ListenAndServe 会直接返回 所以 我们主程序要改由 Shutdown 阻塞
	<-ctx.Done()
	fmt.Println("Received close single.")
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer func() {
		timeoutCancel()
	}()

	if err := server.Shutdown(timeoutCtx); err != nil {
		fmt.Printf("Server Shutdown err; %+v\n", err)
		return
	}
	fmt.Println("Server closed gracefully")

}

func healthz(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 15)

	httpState := http.StatusOK
	fmt.Fprintf(w, "ok")
	fmt.Println("resp ok...")

	defer func() {
		requestInfo(httpState, r)
	}()

}

func requestInfo(httpState int, r *http.Request) {
	cIP := clientIP(r)
	path := r.URL
	fmt.Printf("[%s] [%d] [%s] client ip: %s\n", time.Now(), httpState, path, cIP)
}

//  参照 gin.clientIP()
func clientIP(r *http.Request) string {
	cIP := r.Header.Get("X-Forwarded-For")
	cIP = strings.TrimSpace(strings.Split(cIP, ",")[0])
	if cIP == "" {
		cIP = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	}
	if cIP != "" {
		return cIP
	}

	if addr := r.Header.Get("X-Appengine-Remote-Addr"); addr != "" {
		return addr
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
