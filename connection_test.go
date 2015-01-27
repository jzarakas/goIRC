package main

import (
	"fmt"
	"testing"
)

func TestGetHead(t *testing.T) {
	client := BuildUser()

	client.Host = "127.0.0.1"
	client.Ident = "username"
	client.Nick = "someone"

	if client.GetHead() != fmt.Sprintf(":%s!%s@%s", client.Nick, client.Ident, client.Host) {
		t.Errorf("GetHead() did not parse nick, ident, and host correctly. Returned %q", client.GetHead())
	}
}

func TestHandleNick(t *testing.T) {
	futureNick := "afuturenickname"
	client := BuildUser()
	client.Status = UserPassSent

	handleNick(nil, &client, futureNick, "")

	if client.Nick != futureNick {
		t.Errorf("client.Nick should be %q but is %q", futureNick, client.Nick)
	}

	//make sure that if we aren't registered, we have at least supplied a nickname
	if client.Status != UserNickSent {
		t.Errorf("client.Status is not UserNickSent, it should be!")
	}

}

func TestIsChannel(t *testing.T) {
	aChannel := "#gophers"
	notAChannel := "gophers"

	if isChannel(aChannel) != true {
		t.Errorf("%q is not a valid channel according to isChannel(), it is.", aChannel)
	}

	if isChannel(notAChannel) == true {
		t.Errorf("%q is a valid channel according to isChannel(), it isn't.", notAChannel)
	}
}

func TestHandlePart(t *testing.T) {
	buses := make(map[string]*EventBus)
	client := BuildUser()

	chanName := "#gophers"
	newChannel := Channel{name: chanName, topic: "gogo new channel!", mode: make(map[string]Mode)}
	buses[newChannel.name] = &EventBus{subscribers: make(map[EventType][]Subscriber), channel: &newChannel}

	buses[newChannel.name].Subscribe(UserJoin, &client)
	buses[newChannel.name].Subscribe(UserPart, &client)
	buses[newChannel.name].Subscribe(PrivMsg, &client)

	anotherClient := BuildUser()
	buses[newChannel.name].Subscribe(UserJoin, &anotherClient)
	buses[newChannel.name].Subscribe(UserPart, &anotherClient)
	buses[newChannel.name].Subscribe(PrivMsg, &anotherClient)

	handlePart(buses, &client, newChannel.name, "")

	if len(buses[newChannel.name].GetSubscribers(UserJoin)) != 1 {
		t.Errorf("handlePart did not unsubscribe user")
	}

	handlePart(buses, &anotherClient, newChannel.name, "")

	if buses[newChannel.name] != nil {
		t.Errorf("handlePart did not delete channel")
	}

}

func TestCheckEventBus(t *testing.T) {
	client := BuildUser()
	buses := make(map[string]*EventBus)
	key := "#gophers"
	randomKey := "#gjfkldfjglfdgfd"

	buses[key] = &EventBus{}

	if checkEventBus(buses, &client, key) != true {
		t.Errorf("checkEventBus states %q does not exist. It does.", key)
	}

	if checkEventBus(buses, &client, randomKey) == true {
		t.Errorf("checkEventBus states %q does exist, it doesn't.", randomKey)
	}
}

func TestCheckSubscribe(t *testing.T) {
	client := BuildUser()
	chanName := "#gophers"
	newChannel := Channel{name: chanName, topic: "gogo new channel!", mode: make(map[string]Mode)}
	buses := make(map[string]*EventBus)
	buses[newChannel.name] = &EventBus{subscribers: make(map[EventType][]Subscriber), channel: &newChannel}
	buses[newChannel.name].Subscribe(UserJoin, &client)
	checkSubscribed(buses[newChannel.name], &client, UserJoin)

	if checkSubscribed(buses[newChannel.name], &client, UserJoin) != true {
		t.Errorf("checkSubscribed says %q isn't subscribed to %q, it is.", client.Nick, newChannel.name)
	}

	if checkSubscribed(buses[newChannel.name], &client, UserPart) == true {
		t.Errorf("checkSubscribed says %q is subscribed, but it is not.", client.Nick)
	}
}