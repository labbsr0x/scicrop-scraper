package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.sandmanbb.com/perfil-digital-agro/agro-ws/build/scicrop-scraper/worker"
)

// workerCmd represents the serve command
var workerCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts the worker to pool conductor scicrop tasks",
	Example: `  scicrop-scraper run --conductor-uri=http://localhost:8080 --log-level=debug
	
	  All command line options can be provided via environment variables by adding the prefix "SCICROP_SCRAPER_" 
	  and converting their names to upper case and replacing punctuation and hyphen with underscores. 
	  For example,
	
			command line option                 environment variable
			------------------------------------------------------------------------------------------
			--conductor-host                 							SCISCRAPER_CONDUCTOR_HOST
			--log-level              									SCISCRAPER_LOG_LEVEL
	`,
	RunE: runE,
}


func runE(_ *cobra.Command, _ []string) error {

	config := viper.GetViper()
	logLevel := config.GetString("log-level")
	switch logLevel {
		case "debug":
			logrus.SetLevel(logrus.DebugLevel)
			break
		case "warning":
			logrus.SetLevel(logrus.WarnLevel)
			break
		case "trace":
			logrus.SetLevel(logrus.TraceLevel)
			break
		case "error":
			logrus.SetLevel(logrus.ErrorLevel)
			break
		default:
			logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.Infof("Starting scicrop-scraper with config:")
	for _, key := range config.AllKeys() {
		logrus.Infof("%s : %v", key, config.Get(key))
	}

	worker := worker.NewScicropWorker(config)
	worker.Run()

	return nil
}

func init() {
	rootCmd.AddCommand(workerCmd)
	workerCmd.Flags().StringP("conductor-host", "H", "localhost", "Configure the host to access conductor web server")
	workerCmd.Flags().StringP("conductor-port", "p", "8080", "Configure the port to access conductor web server")
	workerCmd.Flags().StringP("conductor-protocol", "P", "http", "Configure the protocol to access conductor web server")
	workerCmd.Flags().StringP("scicrop-api", "s", "https://engine.scicrop.com/scicrop-engine-web/api/v1", "Sets the URI to connect on Scicrop AgroDataAPI")
	workerCmd.Flags().BoolP("enable-geo", "g", true, "Flag to enable geoindexing of scicrop stations")
	workerCmd.Flags().StringP("agro-uri", "d", "http://localhost:8000", "Sets the URI to connect on AgroDataAPI.")
	workerCmd.Flags().BoolP("check-schemas", "c", true, "Flag to enable schema assertion on schema-api. If true a valid url must be provided on schema-uri parameter.")
	workerCmd.Flags().StringP("schema-uri", "S", "http://localhost:8001/v1", "Sets the URI to connect on  SchemaAPI")
	workerCmd.Flags().StringP("http-timeout", "t", "30", "Sets the http timeout IN SECONDS for Scicrop rest integration")
	workerCmd.Flags().StringSliceP("task-list", "o", []string{"date"}, "Flag to set which tasks the worker will process")
	workerCmd.Flags().StringP("threads-per-task", "T", "2", "Sets the number of threads per task provided on task-list parameter")
	workerCmd.Flags().StringP("pooling-interval", "i", "40", "Sets the interval IN SECONDS between attempts to pool tasks from conductor")
	workerCmd.Flags().StringP("log-level", "l", "info", "Flag to set the log level of the application")

	err := viper.GetViper().BindPFlags(workerCmd.Flags())
	if err != nil {
		panic(err)
	}
}