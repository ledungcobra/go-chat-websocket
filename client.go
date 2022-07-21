package main

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type client struct {
	socket   *websocket.Conn
	send     chan *message
	room     *room
	userData map[string]interface{}
}

func (c *client) read() {
	defer func(socket *websocket.Conn) {
		err := socket.Close()
		if err != nil {
			log.Println(err)
		}
		c.room.leave <- c
		close(c.send)
	}(c.socket)
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			if avatarUrl, ok := c.userData["avatar_url"]; ok {
				msg.AvatarURL = avatarUrl.(string)
			}
			c.room.forward <- msg
		} else {
			break
		}
	}
}

func (c *client) write() {
	defer func(socket *websocket.Conn) {
		err := socket.Close()
		if err != nil {
			log.Println(err)
		}
		c.room.tracer.Trace("Client write disconnected")
	}(c.socket)
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
}
