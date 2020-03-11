package tasks

import (
	"fmt"
	"github.com/ggrcha/conductor-go-client/task"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.sandmanbb.com/perfil-digital-agro/agro-ws/build/scicrop-scraper/client"
	"gitlab.sandmanbb.com/perfil-digital-agro/agro-ws/build/scicrop-scraper/models/schema"
	"gitlab.sandmanbb.com/perfil-digital-agro/agro-ws/build/scicrop-scraper/schemas"
	"strconv"
)

func isValidRangeTask(t *task.Task) bool {
	_, usernameOk := t.InputData["username"];
	_, passwordOk := t.InputData["password"];
	_, startOk := t.InputData["date"];

	return usernameOk && passwordOk && startOk
}

func ScrapRange(t *task.Task) (taskResult *task.TaskResult, err error) {

	taskResult = task.NewTaskResult(t)
	taskResult.OutputData = make(map[string]interface{})

	if isValidRangeTask(t) {
		username, password, date := t.InputData["username"].(string), t.InputData["password"].(string), t.InputData["date"].(string)

		if taskResult.OutputData["date"], err = GetNextDay(date); err != nil {
			logrus.Errorf("Error on task %s output processing reason: %s", t.TaskId, err.Error())
			return taskResult, err
		}

		schemaIntegrationOn := viper.GetBool("check-schemas")
		schemaClient := client.NewSchemaClient(viper.GetViper())
		var schemaList = []*schema.Schema{schemas.Temp24hs, schemas.Humidity24hs, schemas.Radiation24hs, schemas.Wind24hs}

		if schemaIntegrationOn {
			if err := schemaClient.AssertSchemas(schemaList); err != nil {
				logrus.Error(err)
				return taskResult, err
			}
		}

		scicropClient := client.NewScicropClient(viper.GetViper())
		scicropClient.Login(username, password)
		if stationList, statErr := scicropClient.GetStationList(); statErr != nil {
			logrus.Errorf("Error on task %s output processing reason: %s", t.TaskId, statErr.Error())
			return taskResult, statErr
		} else {
			agroClient := client.NewAgrodataClient(viper.GetViper())
			owner := username
			for _, station := range stationList {
				logrus.Infoln("Processing Geolocation for Station", station.Id)
				thing := strconv.Itoa(station.Id)

				logrus.Infoln("Fetching Data for Station", station.Id)
				stationData := scicropClient.GetDailyData(station.Id, date)
				normalizedData := NormalizeDataMap(stationData)
				normalizedData["dateTime"] = fmt.Sprintf("%sT00:00:00Z", date)
				logrus.Infof("Data feched for Station %d:\n%v", station.Id, normalizedData)

				if schemaIntegrationOn {
					if assignErr := schemaClient.AssignSchemaToDomains(owner, thing, schemaList); assignErr != nil {
						logrus.Errorf("Error assign schemas to domains = %v ", assignErr.Error())
						return taskResult, assignErr
					}
				}

				if postErr := agroClient.PostToMultipleNodes(owner, thing, schemaList, normalizedData); postErr != nil {
					logrus.Errorf("Error posting read to agroData = %v ", postErr.Error())
					return taskResult, postErr
				} else {
					if viper.GetBool("enable-geo") {
						if uploadErr := agroClient.UploadGeojson(fmt.Sprintf(`
					{ "type": "FeatureCollection",
					  "features": [
						{ "type": "Feature",
						  "geometry": {"type": "Point", "coordinates": [%f, %f]},
						  "properties": {"owner": "%s", "thing": "%s"}
						}
					  ]
					}`, station.Location.Lon, station.Location.Lat, username, thing)); uploadErr != nil {
							logrus.Errorf("Error uploading geojson to agroData = %v ", uploadErr.Error())
							return taskResult, uploadErr
						}
					}
				}
			}
		}
	} else {
		logrus.Warnf("Received task with invalid input: %v", t.InputData)
	}

	taskResult.Status = "COMPLETED"
	err = nil
	return taskResult, err
}

