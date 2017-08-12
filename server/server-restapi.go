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
	SensorID  string   `json:"SensorID,omitempty"`
	LastValue *Global.ItemInfo //`json:"Value,omitempty"`
	ListValue []*Global.ItemInfo
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

	err := http.ListenAndServe(info.EndpointAddress, router)
	if (err!=nil) {
		lgg.Lgdef.Errorf("ERROR INICIAR HTTP-Endpoint %s. %s",info.EndpointAddress,err)
		panic(err)
	}
}

func GetSensorPoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	//if val, ok := Global.Resources.Store.DataLast[params["id"]]; ok {
	//	//do something here
	//	json.NewEncoder(w).Encode(SensorData{SensorID:params["id"],DataLast: val})
	//}else {
	//	if params["id"]!="favicon.ico" {
	//		lgg.Lgdef.Warnf("HTTPServer: SendorID '%s' no encontrado", params["id"])
	//		json.NewEncoder(w).Encode(SensorData{SensorID: params["id"], DataLast: &Global.ItemInfo{}})
	//	}
	//}

	if val, ok := Global.Resources.Store.DataQueue[params["id"]]; ok {
		//do something here
		//json.NewEncoder(w).Encode(SensorData{SensorID: params["id"], LastValue: val.GetLast(), ListValue: val.ToList()})
		json.NewEncoder(w).Encode(SensorData{SensorID: params["id"], LastValue: val.GetFirst(), ListValue: val.ToList()})
	}else {
		if params["id"]!="favicon.ico" {
			lgg.Lgdef.Warnf("HTTPServer: SendorID '%s' no encontrado", params["id"])
			json.NewEncoder(w).Encode(SensorData{SensorID: params["id"], LastValue: nil ,ListValue: []*Global.ItemInfo{}})
		}
	}
}



