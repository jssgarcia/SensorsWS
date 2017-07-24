package lgg

import (
	"github.com/Sirupsen/logrus"
	//"github.com/gogap/logrus_mate"
	_ "github.com/gogap/logrus_mate/hooks/file"
	"github.com/gogap/logrus_mate"
	"path/filepath"
	"os"
)

var Lgdef *logrus.Logger

func InitLogger() {

	//Comprobar que existe la carpeta log.
	folderLog,_ := filepath.Abs("./log")
	if _,err := os.Stat(folderLog); os.IsNotExist(err) {
		//folder no existe
		os.Mkdir(folderLog,os.ModePerm)
	}

	fconfig := "logger.conf"
	f,_ := filepath.Abs("./" + fconfig)

	mate, _ := logrus_mate.NewLogrusMate(
		logrus_mate.ConfigFile(f),
	)

	Lgdef = mate.Logger("default")
}