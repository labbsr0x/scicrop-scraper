package scicrop

type LoginResponse struct {
	Auth     AuthEntity     `json:"authEntity"`
	Response ResponseEntity `json:"responseEntity"`
}
