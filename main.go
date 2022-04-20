package main

import (
	"github.com/LovesAsuna/ForumSignin/forum"
	"log"
	"sync"
)

type client []forum.Sign

func main() {
	clients := client{
		forum.NewHuaHuoClient(),
		forum.NewZdfxClient(),
	}
	wg := sync.WaitGroup{}
	for _, client := range clients {
		wg.Add(1)
		go func(client forum.Sign) {
			c, _ := client.Sign()
			for m := range c {
				log.Println(m)
			}
			wg.Done()
		}(client)
	}
	wg.Wait()
}
