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
			forum.Debug("开始执行", client.Name(), "的签到操作")
			c, _ := client.Sign()
			forum.Debug(client.Name(), "签到完成，打印结果")
			for m := range c {
				log.Println(m)
			}
			wg.Done()
		}(client)
	}
	wg.Wait()
}
