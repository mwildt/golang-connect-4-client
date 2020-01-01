package main

import (
	"fmt"
	"strconv"
	"strings"
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


