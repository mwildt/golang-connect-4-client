```go
package main

import (
	conn4api "github.com/mwildt/golang-connect-4-client"
	"math/rand"
	"sync"
)

func createBotHandler(name string) conn4api.Handler {

	return func (sender chan string, receiver chan string) {

		bot := conn4api.NewProgrammableClient(sender, receiver)

		bot.NewSeasonHandler = func (seasonId string) {
			bot.JoinSession(seasonId)
		}

		bot.YourTurnHandler = func(token string) {
			bot.Insert(uint8(rand.Intn(6)), token)
		}

		bot.Register(name)
	}
}

func main() {
	var wg sync.WaitGroup
	conn4api.NewUdpConnection("127.0.0.1:4446", 13330).Run(createBotHandler("gobot-1"))
	wg.Add(1)
	wg.Wait()
}
```