package scicrop

type AuthEntity struct {
	User      UserEntity `json:"userEntity"`
	SessionId string          `json:"jSessionId"`
	JwtToken  string          `json:"jwtToken"`
}