package scicrop

type Station struct {
	Id int 							`json:"stationId"`
	Name string 					`json:"name"`
	Data []map[string]interface{}   `json:"stationDataLst"`
	Location struct {
		Lat float64 				`json:"decimalLatLoc"`
		Lon float64 				`json:"decimalLongLoc"`
	}								`json:"location"`
}

type StationPayload struct {
	List []Station				`json:"stationLst"`
}

type StationResponse struct {
	Payload  StationPayload      `json:"payloadEntity"`
	Response ResponseEntity `json:"responseEntity"`
}