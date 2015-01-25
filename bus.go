package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type EventType int

const (
	ChannelUserJoin EventType = iota
	ChannelUserPart
	ChannelPrivMsg
	ChannelMsg
)

type MasterBus struct {
	buses map[string]EventBus
}

type Subscriber interface {
	OnEvent(*Event)
}

type EventBus struct {
	subscribers map[EventType][]Subscriber
	channel     *Channel
}

type Event struct {
	event_type EventType
	event_data string
}

type User struct {
	Nick     string
	Ident    string
	RealName string
	Conn     net.Conn
	Status   ConnectionStatus
}

type Channel struct {
	name  string
	topic string
}

// something funky going on here
// type Subscriber interface {
// 	OnEvent(event *Event)
// }

func (b *MasterBus) addBus(target string) {
	var mutex
}

func (u *User) OnEvent(event *Event) {
	switch event.event_type {
	case ChannelUserJoin:
		//fmt.Printf("%q(%d)> %q\n", s.Nick, event.event_type, event.event_data)
		_, err := u.conn.Write([]byte(event.event_data))
		if err != nil {
			fmt.Println("Not looking too good")
		}
	case ChannelMsg:
		_, err := u.conn.Write([]byte(event.event_data))
		if err != nil {
			fmt.Println("Not looking too good")
		}
	}
}
func (bus *EventBus) Publish(event *Event) {
	fmt.Printf("\npublishing -%d- data: %q\n", event.event_type, event.event_data)
	for _, subscriber := range bus.subscribers[event.event_type] {
		go subscriber.OnEvent(event) //currently slower than without the goroutine
	}
	fmt.Println("done publishing")
}

func (bus *EventBus) Subscribe(event_type EventType, subscriber Subscriber) {
	bus.subscribers[event_type] = append(bus.subscribers[event_type], subscriber)
}

var buses map[string]*EventBus

func init() {
	// init event bus map
	buses = make(map[string]*EventBus)

	// make new channel #gophers
	gophers := Channel{name: "#gophers", topic: "gogo gophergala!"}

	buses[gophers.name] = &EventBus{make(map[EventType][]Subscriber), &gophers}
	fmt.Println("New Channel: " + buses[gophers.name].channel.name)
}

func main() {
	ln, err := net.Listen("tcp", ":3030")
	if err != nil {
		panic("Listen not WORKING")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic("nope not Accepting")
		}
		go handleConnection(conn)
	}

}
