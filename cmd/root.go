package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// rootCmd represents the base command when called without any sub-commands
var rootCmd = &cobra.Command{
	Use:   "scicrop-scraper",
	Short: "Worker to scrap station data from scicrop AgroAPI",
	Long: `Worker application to scrap data from scicrop AgroAPI within the conductor infrastructure on project perfil-digital.`,
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("SCISCRAPER")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv() // read in environment variables that match
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
