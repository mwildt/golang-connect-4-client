package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

type TokenInsertedHandler func(playerName string, colNumber uint8)
type ResultHandler func(result string, playerName string, reason string)
type TokenHandler func(token string)

type ProgrammableClient struct {
	sender chan string
	receiver chan string
	tokenInsertedHandler TokenInsertedHandler
	newSeasonHandler TokenHandler
	newGameHandler TokenHandler
	yourTurnHandler TokenHandler
	resultHandler ResultHandler
}

func (client *ProgrammableClient) onTokenInserted(handler TokenInsertedHandler) {}
func (client *ProgrammableClient) onNewSeason(handler TokenHandler) {}
func (client *ProgrammableClient) onNewGame(handler TokenHandler) {}
func (client *ProgrammableClient) onYourTurn(handler TokenHandler) {}

func (client *ProgrammableClient) register(name string) {
	client.sender<- fmt.Sprintf("REGISTER:%s", name)
}
func (client *ProgrammableClient) unregister() {
	client.sender<- "UNREGISTER"
}
func (client *ProgrammableClient) joinSession(sessionId string) {
	client.sender<- fmt.Sprintf("JOIN:%s", sessionId)
}
func (client *ProgrammableClient) insert(colNumber uint8, token string) {
	client.sender<- fmt.Sprintf("INSERT;%d;%s",colNumber, token)
}

func (client *ProgrammableClient) run(){
	for {
		msg := <-client.receiver
		fmt.Print("gOT MESSAGE FROM CHANNEL", msg)
		parts := strings.Split(msg, ";")

		switch parts[0] {

		case "NEW SEASON" : client.newSeasonHandler(parts[1])

		case "NEW GAME" : client.newGameHandler(parts[1])

		case "TOKEN INSERTED" :
			colNr, err := strconv.ParseUint(parts[2],2, 8)
			checkError(err)
			client.tokenInsertedHandler(parts[1], uint8(colNr))

		case "YOURTURN" : client.yourTurnHandler(parts[1])

		case "RESULT" : client.resultHandler(parts[1], parts[2], parts[3])

		}
	}
}


func NewProgrammableClient(sender chan string, receiver chan string) *ProgrammableClient {
	bot := & ProgrammableClient{
		sender: sender,
		receiver: receiver,
		tokenInsertedHandler: func(playerName string, colNumber uint8) {},
		newSeasonHandler: func(token string) {},
		newGameHandler: func(token string) {},
		yourTurnHandler: func(token string) {},
		resultHandler: func(result string, playerName string, reason string) {},
	}
	go bot.run()
	return bot
}

func createBotHandler(name string) Handler {
	fmt.Println("create BotHandler")
	return func (sender chan string, receiver chan string) {
		fmt.Println("run BotHandler")
		bot := NewProgrammableClient(sender, receiver)

		bot.newSeasonHandler = func (seasonId string) {
			fmt.Printf("[%s] -> NEW SEASON | SEND JOIN\n", name)
			bot.joinSession(seasonId)
		}

		bot.newGameHandler = func (seasonId string) {
			fmt.Printf("[%s] -> NEW GAME\n", name)
		}

		bot.yourTurnHandler = func(token string) {
			bot.insert(uint8(rand.Intn(6)), token)
		}

		bot.register(name)
	}
}

func createHandler(name string) Handler {

	return func (sender chan string, receiver chan string) {

		fmt.Println("running handler fn")
		sender<- fmt.Sprintf("REGISTER;%s", name)

		seasonCount := 0

		send := func(msg string) {
			fmt.Printf("[%s] -> [%s]\n", name, msg)
			sender<- msg
		}

		for {
			msg := <-receiver
			fmt.Printf("[%s] <- %s\n", name, msg)
			parts := strings.Split(msg, ";")

			switch parts[0] {

			case "NEW SEASON" :
				if seasonCount < 1 {
					seasonCount += 1
					send(fmt.Sprintf("JOIN;%s", parts[1]))
				} else {
					send("UNREGISTER")
				}

			case "NEW GAME" :
				fmt.Printf("[%s] Attend to new Game %s\n", name, parts[1])

			case "YOURTURN" :
				send(fmt.Sprintf("INSERT;%d;%s", rand.Intn(6), parts[1]))
				fmt.Printf("[%s] Attend to new Game %s\n", name, parts[1])

			case "TOKEN INSERTED" :
				fmt.Printf("[%s] TOKEN INSERTED in column %s by %s\n", name, parts[2], parts[1])

			case "RESULT" :
				fmt.Printf("[%s] RESULT %s %s %s\n", name, parts[2], parts[1], parts[3])
			}
		}
	}
}

func main() {
	fmt.Println("running client")

	NewUdpConnection("127.0.0.1:4446", 13330).Run(createHandler("gobot-1"))
	NewUdpConnection("127.0.0.1:4446", 13331).Run(createBotHandler("gobot-2"))

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}