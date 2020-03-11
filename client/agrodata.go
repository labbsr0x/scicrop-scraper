package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.sandmanbb.com/perfil-digital-agro/agro-ws/build/scicrop-scraper/models/schema"
	"net/http"
	"net/url"
	"time"
)

type AgrodataClient struct {
	uri       *url.URL
	jwtToken  string
	sessionId string
	client    *http.Client
}

func NewAgrodataClient(v *viper.Viper) *AgrodataClient {
	instance := new(AgrodataClient)

	uri, err := url.Parse(v.GetString("agro-uri"))
	if err != nil {
		logrus.Errorf("Invalid agro-data-uri parameter. Panicking... ")
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

func (c *AgrodataClient) PostSingleRead(owner string, thing string, node string, attributes map[string]string ) error {
	logrus.Tracef("Calling post read to agro-data on host %s for owner %s thing %s node %s ", c.uri.String(), owner, thing, node)
	if requestBody, marshErr := json.Marshal(attributes); marshErr == nil {

		req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v1/owner/%s/thing/%s/node/%s", c.uri.String(), owner, thing, node), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "*/*")
		logrus.Debugln("Making post request to agrodata")
		logrus.Tracef("Payload %v", req)
		r, err := c.client.Do(req)

		if err != nil {
			return fmt.Errorf("Couldnt post read to agrodata reason: %s", err);
		} else {
			defer r.Body.Close()
			if r.StatusCode >= 300 {
				return fmt.Errorf("status code %d", );
			}
		}

	} else {
		return fmt.Errorf("Marshall error on posting read to agro-data")
	}

	logrus.Debugf("Successfully posted read to agro-data for owner %s thing %s node %s with attributes %v", owner, thing, node, attributes)
	return nil
}

func (c *AgrodataClient) UploadGeojson(json string) error {
	logrus.Tracef("Calling post upload method to agro-data %s ", json)
	requestBody := []byte(json)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v0/files/upload/geojson", c.uri.String()), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	logrus.Debugln("Making post upload request to agrodata")
	logrus.Tracef("Payload %v", req)
	r, err := c.client.Do(req)

	if err != nil {
		return fmt.Errorf("Couldnt upload geoson to agrodata reason: %s", err);
	} else {
		defer r.Body.Close()
		if r.StatusCode >= 300 {
			return fmt.Errorf("status code %d", r.StatusCode);
		}
	}

	logrus.Debugf("Successfully posted geojson to agro-data")
	return nil
}

/*
	Post data to multiple schemas, one by one.
	Each schema has a different node assigned on Node property.
 */
func (c *AgrodataClient) PostToMultipleNodes(owner string, thing string, schemas []*schema.Schema, data map[string]string) (retErr error) {
	for _, schema := range schemas {
		if retErr != nil {
			break
		}
		retErr = c.PostSingleRead(owner, thing, schema.Node, MapUsingSchema(data, *schema))
	}
	if retErr != nil {
		logrus.Error(retErr)
	}
	return retErr
}