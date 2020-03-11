package client

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.sandmanbb.com/perfil-digital-agro/agro-ws/build/scicrop-scraper/models/scicrop"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const SESSION_COOKIE_NAME = "JSESSIONID"

type ScicropClient struct {
	uri       *url.URL
	jwtToken  string
	sessionId string
	client    *http.Client
}

func NewScicropClient(v *viper.Viper) *ScicropClient {
	instance := new(ScicropClient)

	uri, err := url.Parse(v.GetString("scicrop-api"))
	if err != nil {
		logrus.Errorf("Invalid conductor-uri parameter. Panicking... ")
		panic(err)
	}
	instance.uri = uri

	httpTimeout := v.GetUint("http-timeout")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	instance.client = &http.Client{Timeout: time.Duration(httpTimeout) * time.Second, Transport: tr}

	return instance
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *ScicropClient) Login(username string, password string) bool {

	logrus.Debugln("Logging on agrodataAPI with user ", username)

	var jsonStr = []byte(fmt.Sprintf(`{"authEntity":{"userEntity":{"email":"%s"}}}`, username))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.uri.String(), "userauth/token"), bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization","Basic " + basicAuth(username,password))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")

	logrus.Debugln("Making login request")
	logrus.Tracef("Login request payload %v", req)
	r, err := c.client.Do(req)

	if err != nil {
		logrus.Errorf("erro client.Do(): ", err)
	}

	defer r.Body.Close()

	/*
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorln("Could not deserialize login response")
	}
	logrus.Tracef("Login http response: %s", string(b))

	 */

	response := new(scicrop.LoginResponse)
	json.NewDecoder(r.Body).Decode(response)
	logrus.Debugf("Response body %v", response)

	if loginOk := response.Response.ReturnId == 0; loginOk {
		logrus.Infoln("Logged in with username", username)
		c.jwtToken = response.Auth.JwtToken
		c.sessionId = response.Auth.SessionId
		logrus.Debugln("Obtained JwtToken=", c.jwtToken, " AND SESSIONID ", c.sessionId)
		return true
	}

	logrus.Errorln("Login attempt was refused for username", username)
	return false

}

func (c *ScicropClient) GetStationList() ([]scicrop.Station, error) {

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.uri.String(), "station/myAccess"), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	cookie := http.Cookie{Name: SESSION_COOKIE_NAME, Value: c.sessionId}
	req.AddCookie(&cookie)
	req.Header.Set("Authorization","Bearer " + c.jwtToken)

	logrus.Debugln("Making request to get station list")
	logrus.Tracef("Station list request payload %v", req)
	r, err := c.client.Do(req)

	if err != nil {
		logrus.Errorf("erro client.Do(): ", err)
		return nil, fmt.Errorf("Error making station request to scicrop: %s", err.Error())
	}

	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	logrus.Tracef("Station List http response: " + string(b))

	response := new(scicrop.StationResponse)
	if err := json.Unmarshal(b, response); err != nil {
		return nil, fmt.Errorf("Error unmarshaling stations: %s", err.Error())
	}
	logrus.Debugf("Station List Response body %v", response)

	if reqOk := response.Response.ReturnId == 0; reqOk {
		return response.Payload.List, nil
	}

	logrus.Errorln("Get Station list request failed")
	return nil, fmt.Errorf("Scicrop response message %s: ", response.Response.ReturnMsg)

}

func (c *ScicropClient) GetDailyData(stationId int, date string) map[string]interface{}  {

	var jsonStr = []byte(fmt.Sprintf(`{ "payloadEntity":{"stationLst":[{"stationId":%d,"stationDataLst":[{"date":"%s"}]}]}}`, stationId, date))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.uri.String(), "station/getScicropStationDataByDayOnStationDataT"), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	cookie := http.Cookie{Name: SESSION_COOKIE_NAME, Value: c.sessionId}
	req.AddCookie(&cookie)
	req.Header.Set("Authorization","Bearer " + c.jwtToken)


	logrus.Debugln("Making request to get Daily Data")
	logrus.Tracef("Daily Data request payload %v", req)
	r, err := c.client.Do(req)

	if err != nil {
		logrus.Errorf("erro client.Do(): ", err)
	}

	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	logrus.Tracef("Daily Data http response: " + string(b))

	response := new(scicrop.StationResponse)
	//response.Payload.List = make([]models.Station, 1)
	//response.Payload.List[0].Data = make([]map[string]float64, 1)
	//response.Payload.List[0].Data[0] = make(map[string]float64)

	if err := json.Unmarshal(b, response); err != nil {
		logrus.Warnf("Unmarshall error on daily data request")
	} else {
		logrus.Debugf("Station List Response body %v", response)
	}

	if reqOk := response.Response.ReturnId == 0; reqOk {
		return response.Payload.List[0].Data[0]
	}

	logrus.Errorln("Daily Data request failed")
	return nil

}
