package main

import (
	"net"
	"bufio"
	"fmt"
	"strings"
	_ "log"
	"github.com/tkanos/gonfig"
	_ "os"
	"path/filepath"
)

type configuration struct {
	ListenAddress string `json:"listenAddress"`
}

var cnfg configuration

func main() {

	getConfig()

	fmt.Print("Launching server ...")

	ln,err := net.Listen("tcp", cnfg.ListenAddress)
	if err!=nil {
		panic("Interface ")
	}
	conn,_ := ln.Accept()

	for {
		mess,_ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Mesagge receive:",mess)

		newmess := strings.ToUpper(mess)
		conn.Write([]byte(newmess+ "\n"))

	}
}

func getConfig() {

	cnfg = configuration{}

	fconfig := "config.dev.json"

	//Defaults
	cnfg.ListenAddress="127.0.0.1:8888"

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


