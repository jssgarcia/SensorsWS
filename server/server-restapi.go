package server

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"ty/csi/ws/SensorsWS/Global"
	"ty/csi/ws/SensorsWS/lgg"
)

//region Variables

type HttpServerInfo struct {
	EndpointName string
	EndpointAddress string
}

type SensorData struct {
	SensorID string   `json:"SensorID,omitempty"`
	Data *Global.ItemInfo   //`json:"Value,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

//endregion

func InitServer(info HttpServerInfo) {

	if info.EndpointAddress =="" {
		panic("ERROR Initiation: Server adddres is not defined. " + info.EndpointAddress)
	}

	router := mux.NewRouter()
	router.HandleFunc("/{id}", GetSensorPoint).Methods("GET")

	http.ListenAndServe(info.EndpointAddress, router)
}



func GetSensorPoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	if val, ok := Global.Resources.Store.Data[params["id"]]; ok {
		//do something here
		json.NewEncoder(w).Encode(SensorData{SensorID:params["id"],Data: val})
	}else {
		json.NewEncoder(w).Encode(SensorData{SensorID:params["id"],Data:&Global.ItemInfo{}})
	}

	lgg.Lgdef.Infof("STORE Global: %s", Global.Resources.Store)
}


