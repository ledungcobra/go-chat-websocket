package main

import (
	"complexhttp/tracer"
	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
	"log"
	"net/http"
	"runtime"
)

type room struct {
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  tracer.Tracer
	avatar  Avatar
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		case client := <-r.leave:
			r.tracer.Trace("Disconnected client")
			delete(r.clients, client)
		case msg := <-r.forward:
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send message to channel
					r.tracer.Trace("Sending message to client")
				default:
					r.tracer.Trace("Client is not ready")
					// Failed to send
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Panic("Cannot server because of error")
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Panic("Cannot server because of error")
	}

	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() {
		r.leave <- client
	}()
	go client.write()
	client.read()
	log.Printf("Number of groutines %d", runtime.NumGoroutine())
}

func newRoom(avatar Avatar) *room {
	r := &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		avatar:  avatar,
	}
	return r
}
