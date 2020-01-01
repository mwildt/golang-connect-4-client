package conn4api

import (
	"testing"
	"time"
)

func initTest() (chan string, chan string, *ProgrammableClient) {
	sender := make(chan string, 1)
	receiver := make(chan string, 1)
	handler := NewProgrammableClient(sender, receiver)
	go handler.run()
	return sender, receiver, handler
}

func checkMessage(channel chan string, expected string, t *testing.T) {
	select {
	case <-time.After(1 * time.Second):
		t.Error("timeout")
	case msg := <-channel:
		{
			if msg != expected {
				t.Errorf("expected [%s] but was [%s]", expected, msg)
			}
		}
	}
}

func TestNEW_SEASON(t *testing.T) {
	_, receiver, handler := initTest()
	tokenChan := make(chan string)
	handler.newSeasonHandler = func(token string) {
		tokenChan <- token
		defer handler.Close()
	}
	receiver<- "NEW SEASON;token-123"
	checkMessage(tokenChan, "token-123", t)
}

func TestNEW_GAME(t *testing.T) {
	_, receiver, handler := initTest()
	tokenChan := make(chan bool)
	handler.newSeasonHandler = func(token string) {
		tokenChan <- true
		defer handler.Close()
	}
	receiver<- "NEW SEASON;token-123"
	select {
	case <-time.After(1 * time.Second):
		t.Error("timeout")
	case <-tokenChan:
	}
}

func TestYOURTURN(t *testing.T) {
	_, receiver, handler := initTest()
	tokenChan := make(chan string)
	handler.yourTurnHandler = func(token string) {
		tokenChan <- token
		defer handler.Close()
	}
	receiver<- "YOURTURN;token-456"
	checkMessage(tokenChan, "token-456", t)
}

func TestTOKEN_INSERTED(t *testing.T) {
	_, receiver, handler := initTest()
	type InsertTuple struct {
		playerName string
		col uint8
	}

	tokenChan := make(chan InsertTuple)
	handler.tokenInsertedHandler = func(playerName string, col uint8) {
		tokenChan <- InsertTuple{playerName, col}
		defer handler.Close()
	}
	receiver<- "TOKEN INSERTED;Paul;4"

	select {
	case <-time.After(1 * time.Second):
		t.Error("timeout")
	case insertToken :=  <- tokenChan: {
		if insertToken.playerName != "Paul" {
			t.Errorf("expected [%s] but was [%s]", "Paul", insertToken.playerName)
		}
		if insertToken.col != 4 {
			t.Errorf("expected [%d] but was [%d]", 4, insertToken.col)
		}
	}
	}

}

func TestINSERT(t *testing.T) {
	sender, _, handler := initTest()
	handler.insert(5, "ABC-123")
	checkMessage(sender, "INSERT;5;ABC-123", t)
	handler.Close()
}

func TestREGISTER(t *testing.T) {
	sender, _, handler := initTest()
	handler.register("peter")
	checkMessage(sender, "REGISTER;peter", t)
	handler.Close()
}

func TestUNREGISTER(t *testing.T) {
	sender, _, handler := initTest()
	handler.unregister()
	checkMessage(sender, "UNREGISTER", t)
	handler.Close()
}

func TestJOIN_SESSION(t *testing.T) {
	sender, _, handler := initTest()
	handler.joinSession("ABC-DEF")
	checkMessage(sender, "JOIN;ABC-DEF", t)
	handler.Close()
}