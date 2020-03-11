package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.sandmanbb.com/perfil-digital-agro/agro-ws/build/scicrop-scraper/models/schema"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type SchemaClient struct {
	uri       *url.URL
	client    *http.Client
}

func (c *SchemaClient) setJsonContent(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
}

func (c *SchemaClient) makeRequest(req *http.Request) (*http.Response, error) {
	logrus.Debugf("Making %s request to %s", req.Method, req.URL.String())
	logrus.Debugf("Payload %v", req)
	return c.client.Do(req)
}

func NewSchemaClient(v *viper.Viper) *SchemaClient {
	instance := new(SchemaClient)

	uri, err := url.Parse(v.GetString("schema-uri"))
	if err != nil {
		logrus.Errorf("Invalid schema-uri parameter. Panicking... ")
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

func (c *SchemaClient) CreateSchema(schema *schema.Schema) error  {
	if requestBody, marshErr := json.Marshal(schema); marshErr == nil {
		req, _ := http.NewRequest("POST", fmt.Sprintf("%s/schema", c.uri.String()), bytes.NewBuffer(requestBody))
		c.setJsonContent(req)
		r, err := c.makeRequest(req)

		if err != nil {
			return fmt.Errorf("Couldnt post read to schema-api reason: %s", err);
		} else {
			defer r.Body.Close()
			if r.StatusCode >= 300 {
				return fmt.Errorf("Post to schema-api was not accepted, status code %d", r.StatusCode);
			}
		}

	} else {
		return fmt.Errorf("Marshall error on posting schema to schema-api")
	}

	logrus.Debugf("Successfully posted schema to schema-api.")
	return nil
}

func (c *SchemaClient) GetSchema(schema_uri string)(*schema.Schema, error) {
	logrus.Traceln(fmt.Sprintf("Calling GetSchema to schema-api on host %s", c.uri.String()))
	requestBody, marshErr := json.Marshal(struct{}{})
	if marshErr != nil {
		return nil, fmt.Errorf("Marshall error on getting schema from schema-api: %s", marshErr)
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/schema/%s", c.uri.String(), schema_uri), bytes.NewBuffer(requestBody))
	c.setJsonContent(req)
	r, err := c.makeRequest(req)

	if err != nil {
		return nil, fmt.Errorf("Couldnt get schema on schema-api reason: %s", err);
	}
	defer r.Body.Close()
	if r.StatusCode >= 300 {
		if r.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("Get schema was not accepted on schema-api, status code %d", r.StatusCode);
	}
	b, _ := ioutil.ReadAll(r.Body)
	response := &schema.Schema{}
	umErr := json.Unmarshal(b, response)
	if umErr != nil {
		return nil, fmt.Errorf("Unmarshall error on getting schema from schema-api: %s", umErr)
	}
	logrus.Debugf("Get Schema From Schema-API - Response body %v", *response)
	return response, nil
}

func (c *SchemaClient) AssignSchema(schema_uri string, domain string) error  {
	assignPayload := struct { Domain string `json:"domain"`}{ Domain: domain}
	if requestBody, marshErr := json.Marshal(assignPayload); marshErr == nil {
		req, _ := http.NewRequest("POST", fmt.Sprintf("%s/schema/%s/assign", c.uri.String(), schema_uri), bytes.NewBuffer(requestBody))
		c.setJsonContent(req)
		r, err := c.makeRequest(req)

		if err != nil {
			return fmt.Errorf("Couldnt assign schema on schema-api to domain %s. Reason: %s", domain, err);
		} else {
			defer r.Body.Close()
			if r.StatusCode >= 300 {
				return fmt.Errorf("Assign post to schema-api was not accepted, status code %d", r.StatusCode);
			}
		}

	} else {
		return fmt.Errorf("Marshall error on assigning schema on schema-api.")
	}

	logrus.Debugf("Successfully assigned schema on schema-api.")
	return nil
}

func (c *SchemaClient) AssertSchema(targetSchema *schema.Schema) (error) {
	if existingSchema, err := c.GetSchema(targetSchema.Uri); err != nil {
		return fmt.Errorf("Error asserting schema %s: %s", targetSchema.Uri, err)
	} else if existingSchema == nil {
		return c.CreateSchema(targetSchema)
	}
	return nil
}

func (c *SchemaClient) AssertSchemas(schemas []*schema.Schema) (retErr error) {
	for _, schema := range schemas {
		if retErr != nil {
			break
		}
		retErr = c.AssertSchema(schema)
	}
	if retErr != nil {
		logrus.Error(retErr)
	}
	return retErr
}

func (c *SchemaClient) AssignSchemaToDomains(owner string, thing string, schemas []*schema.Schema) (retErr error) {
	for _, schema := range schemas {
		if retErr != nil {
			break
		}
		retErr = c.AssignSchema(schema.Uri, fmt.Sprintf("owner/%s/thing/%s/node/%s", owner, thing, schema.Node))
	}
	if retErr != nil {
		logrus.Error(retErr)
	}
	return retErr
}