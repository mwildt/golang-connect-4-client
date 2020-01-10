package conn4api

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
	TokenInsertedHandler TokenInsertedHandler
	NewSeasonHandler TokenHandler
	NewGameHandler TokenHandler
	YourTurnHandler TokenHandler
	ResultHandler ResultHandler
}

func NewProgrammableClient(sender chan string, receiver chan string) *ProgrammableClient {
	bot := & ProgrammableClient{
		sender: sender,
		receiver: receiver,
		TokenInsertedHandler: func(playerName string, colNumber uint8) {},
		NewSeasonHandler: func(token string) {},
		NewGameHandler: func(token string) {},
		YourTurnHandler: func(token string) {},
		ResultHandler: func(result string, playerName string, reason string) {},
	}
	go bot.Run()
	return bot
}

func (client *ProgrammableClient) OnTokenInserted(handler TokenInsertedHandler) {}

func (client *ProgrammableClient) OnNewSeason(handler TokenHandler) {}

func (client *ProgrammableClient) OnNewGame(handler TokenHandler) {}

func (client *ProgrammableClient) OnYourTurn(handler TokenHandler) {}

func (client *ProgrammableClient) Register(name string) {
	client.sender<- fmt.Sprintf("REGISTER;%s", name)
}
func (client *ProgrammableClient) Unregister() {
	client.sender<- "UNREGISTER"
}
func (client *ProgrammableClient) JoinSession(sessionId string) {
	client.sender<- fmt.Sprintf("JOIN;%s", sessionId)
}
func (client *ProgrammableClient) Insert(colNumber uint8, token string) {
	client.sender<- fmt.Sprintf("INSERT;%d;%s",colNumber, token)
}

func (client *ProgrammableClient) Run() {
	var channelOpen = true
	for channelOpen {
		msg, ok := <-client.receiver
		if ok {
			parts := strings.Split(msg, ";")

			switch parts[0] {

			case "NEW SEASON" : client.NewSeasonHandler(parts[1])

			case "NEW GAME" : client.NewGameHandler(parts[1])

			case "TOKEN INSERTED" :
				colNr, err := strconv.ParseUint(parts[2],10, 8)
				checkError(err)
				client.TokenInsertedHandler(parts[1], uint8(colNr))

			case "YOURTURN" : client.YourTurnHandler(parts[1])

			case "RESULT" : client.ResultHandler(parts[1], parts[2], parts[3])

			}
		} else {
			channelOpen = false
		}
	}
}

func (client *ProgrammableClient) Close() {
	close(client.sender)
	close(client.receiver)
}


