package main

import "log"

func checkError(err error){
	if nil != err {
		log.Fatal(err)
	}
}