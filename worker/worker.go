package worker

import (
	"fmt"
	"github.com/ggrcha/conductor-go-client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.sandmanbb.com/perfil-digital-agro/agro-ws/build/scicrop-scraper/worker/tasks"
	"net/url"
)

type ScicropWorker struct {
	conductorUri *url.URL
	taskList []string
	threadCount int
	pollingInterval int

	conductorWorker *conductor.ConductorWorker
}

func NewScicropWorker(v *viper.Viper) *ScicropWorker {
	instance := new(ScicropWorker)

	uri, err := url.Parse(fmt.Sprintf("%s://%s:%s/api", v.GetString("conductor-protocol"), v.GetString("conductor-host"), v.GetString("conductor-port")))
	if err != nil {
		logrus.Errorf("Invalid conductor-uri parameter. Panicking... ")
		panic(err)
	}

	instance.taskList = v.GetStringSlice("task-list")
	instance.threadCount = v.GetInt("threads-per-task")
	instance.pollingInterval = v.GetInt("pooling-interval")
	instance.conductorWorker = conductor.NewConductorWorker(uri.String(), instance.threadCount, 1000 * instance.pollingInterval)

	return instance
}

func (w *ScicropWorker) Run() {

	for idx, taskName := range w.taskList {
		wait := idx == len(w.taskList)-1
		switch taskName {
			case "last_read":
				w.conductorWorker.Start("scicrop_last_read", tasks.ScrapLast, wait)
			case "date":
				w.conductorWorker.Start("scicrop_date", tasks.ScrapRange, wait)
			default:
				logrus.Errorf("Invalid task name received -> %s\n The valid options are: \"last_read\" | \"date\"", taskName)
		}
	}
}

func (w *ScicropWorker) Stop() {
	// To be developed
}
