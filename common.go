package conn4api

import "log"

func checkError(err error){
	if nil != err {
		log.Fatal("FATAL ERROR ", err)
	}
}