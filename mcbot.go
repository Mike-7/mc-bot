package main

import (
	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/chat"
	"os"
	"fmt"
	"log"
	"time"
	"net"
	"strings"
)

const timeout = 10

var (
	c *bot.Client
	watch chan time.Time
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mcbot <hostname>")
		return
	}

	c = bot.NewClient()
	c.Name = "Steve [BOT]"

	//Login
	addr, port := GetAddrPort(os.Args[1])
	err := c.JoinServer(addr, port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Login success")

	//Register event handlers
	c.Events.GameStart = onGameStart
	c.Events.ChatMsg = onChatMsg
	c.Events.Disconnect = onDisconnect
	c.Events.Die = onDie

	//JoinGame
	err = c.HandleGame()
	if err != nil {
		log.Fatal(err)
	}
}

func GetAddrPort(hostname string) (string, int) {
	_, srvs, err := net.LookupSRV("minecraft", "tcp", hostname)
	if err != nil {
		panic(err)
	}

	srv := srvs[0]
	addr := strings.TrimSuffix(srv.Target, ".")
	port := srv.Port

	return addr, int(port)
}

func onDie() error {
	time.AfterFunc(3 * time.Second, func() {
		c.Respawn()
		c.Chat("/tp -118.5 58 -139.5")
	})

	return nil;
}

func onGameStart() error {
	log.Println("Game start")
	c.Chat("/tp -118.5 58 -139.5")

	watch = make(chan time.Time)
	go watchDog()

	return c.UseItem(0)
}

func onChatMsg(c chat.Message, pos byte) error {
	log.Println("Chat:", c)
	return nil
}

func onDisconnect(c chat.Message) error {
	log.Println("Disconnect:", c)
	return nil
}

func watchDog() {
	to := time.NewTimer(time.Second * timeout)
	for {
		select {
		case <-watch:
		case <-to.C:
			if err := c.UseBlock(0, -119, 58, -141, 0, 0, 0, 0, false); err != nil {
				
			}
		}
		to.Reset(time.Second * timeout)
	}
}
