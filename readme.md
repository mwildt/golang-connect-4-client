```go
package main

import (
	"math/rand"
	"sync"
	""
)

func createBotHandler(name string) Handler {

	return func (sender chan string, receiver chan string) {

		bot := NewProgrammableClient(sender, receiver)

		bot.newSeasonHandler = func (seasonId string) {
			bot.joinSession(seasonId)
		}

		bot.yourTurnHandler = func(token string) {
			bot.insert(uint8(rand.Intn(6)), token)
		}

		bot.register(name)
	}
}

func main() {
	var wg sync.WaitGroup
	NewUdpConnection("127.0.0.1:4446", 13330).Run(createBotHandler("gobot-1"))
	wg.Add(1)
	wg.Wait()
}
```