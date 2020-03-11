package tasks

import (
	"fmt"
	"github.com/ggrcha/conductor-go-client/task"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.sandmanbb.com/perfil-digital-agro/agro-ws/build/scicrop-scraper/client"
)

func isValidLastTask(t *task.Task) bool {
	_, usernameOk := t.InputData["username"];
	_, passwordOk := t.InputData["password"];
	return usernameOk && passwordOk
}

func ScrapLast(t *task.Task) (taskResult *task.TaskResult, err error) {
	if isValidLastTask(t) {
		client := client.NewScicropClient(viper.GetViper())
		client.Login(t.InputData["username"].(string), t.InputData["password"].(string))
		stationList, _ := client.GetStationList()
		for _, station := range stationList {
			fmt.Println("Station ", station.Id, station.Name)
		}

	} else {
		logrus.Warnf("Received task with invalid input: %v", t.InputData)
	}

	taskResult = task.NewTaskResult(t)
	taskResult.Status = "COMPLETED"
	err = nil
	return taskResult, err
}