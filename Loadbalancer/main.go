package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server interface {
	Address() string
	IsALive() bool
	Serve(w http.ResponseWriter, r *http.Request)
}

type simpleServer struct {
	address string
	proxy   *httputil.ReverseProxy
}

type Loadbalancer struct {
	port          string
	roundRobinCnt int
	server        []Server
}

func NewLoadBalancer(port string, servers []Server) *Loadbalancer {
	return &Loadbalancer{
		port:          port,
		roundRobinCnt: 0,
		server:        servers,
	}
}

func newServer(address string) *simpleServer {
	serverUrl, err := url.Parse(address)
	if err != nil {
		log.Fatal(err)
	}

	return &simpleServer{
		address: address,
		proxy:   httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func (s *simpleServer) Address() string {
	return s.address
}

func (s *simpleServer) IsAlive() bool {
	return true
}

func (s *simpleServer) Serve(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
}

func (lb *Loadbalancer) getNextAvailableServer() Server {
	ser := lb.server[lb.roundRobinCnt%len(lb.server)]
	for !ser.IsALive() {
		lb.roundRobinCnt++
		ser = lb.server[lb.roundRobinCnt%len(lb.server)]
	}
	lb.roundRobinCnt++
	return ser
}

func (lb *Loadbalancer) serveProxy(w http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Println("Forwarding request to address: ", targetServer.Address())
	targetServer.Serve(w, r)
}

func main() {
	servers := []Server{
		newServer("https://www.google.com"),
		newServer("https://www.github.com"),
		newServer("https://www.instagram.com"),
	}

	lb := NewLoadBalancer("4000", servers)

	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		lb.serveProxy(w, r)
	}
	http.HandleFunc("/", handleRedirect)

	fmt.Println("listening at port: ", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
