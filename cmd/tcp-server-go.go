package main

import (
	"net"
	//"bufio"
	"fmt"
	_ "log"
	"github.com/tkanos/gonfig"
	"os"
	"path/filepath"
	"time"
)

type configuration struct {
	ListenAddress string `json:"listenAddress"`
}

var cnfg configuration
var _conn net.Conn

func main() {

	//getConfig()

	if len(os.Args) > 1 {
		cnfg.ListenAddress = os.Args[1]
	} else {
		panic("ServerAdress is not defined")
	}

	fmt.Printf("Launching server ... %v",cnfg.ListenAddress)

	ln,err := net.Listen("tcp", cnfg.ListenAddress)
	if err!=nil {
		panic("Interface ")
	}

	for {
		conn, _ := ln.Accept()
		_conn = conn;

		//
		tmr := time.NewTicker(2 * time.Second)

		bcontinue := true

		for bcontinue {
			select {
			case <-tmr.C:
				err := sendData(conn, "PESO1#PESO2")
				if err != nil {
					fmt.Printf("SERVER: error send data '%s'\n",err)
					bcontinue=false
				}

				break
			}
		}
	}

	fmt.Printf("SERVER == Finalizado main ==")
}

func sendData(conn net.Conn, data string) error {
	_,err := conn.Write([]byte(data+"\n"))
	fmt.Println("TCPServer: message send ..." + data)

	return err
}

func getConfig() {

	cnfg = configuration{}

	fconfig := "config.dev.json"

	//Defaults
	//cnfg.ListenAddress="127.0.0.1:8888"

	f,_ := filepath.Abs("./server/" + fconfig)

	//Obtenemos la configuracion
	err := gonfig.GetConf(f,&cnfg)
	if err!=nil {
		panic("ERROR obtener configuracion: " + err.Error())
	}

}

//func startTcpclient() {
//
//	ln,err := net.Listen("tcp","127.0.0.1:35349")
//	if err!=nil {
//		log.Fatal("FATAL")
//	}
//
//}


