package scicrop

type ResponseEntity struct {
	EpochTime   int 	`json:"epochTime"`
	ReturnId 	int 	`json:"returnId"`
	ReturnMsg   string  `json:"returnMsg"`
}