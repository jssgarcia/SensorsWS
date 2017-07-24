package main

import (
	"path/filepath"
	"github.com/tkanos/gonfig"
	"ty/csi/ws/SensorsWS/Global"
	"context"
	"time"
	"ty/csi/ws/SensorsWS/TcpClient"
	"github.com/kardianos/service"
	"log"
	"github.com/Sirupsen/logrus"
	"github.com/gogap/logrus_mate"

)

//region Service-Functions
type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}
//endregion


func main() {

	//Obtenemos la configuracion
	getConfig()

    ctx,cancel := context.WithCancel(context.Background())
	defer cancel()

	TcpClient.InitClient(ctx)

	for n:=0;n<10;n++ {
		time.Sleep(time.Second)
	}
}

func initService(){

	svcConfig := &service.Config{
		Name:        "GoServiceExampleSimple",
		DisplayName: "Go Service Example",
		Description: "This is an example Go service.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

func getConfig() {

	cnfg := Global.Configuration{}


	fconfig := "config.json"

	//Defaults
	cnfg.ServerAddress="127.0.0.1:8888"

	f,_ := filepath.Abs("./" + fconfig)

	//Obtenemos la configuracion
	err := gonfig.GetConf(f,&cnfg)
	if err!=nil {
		panic("ERROR obtener configuracion: " + err.Error())
	}

	Global.Resources.Config = cnfg

}