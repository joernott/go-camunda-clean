package cmd

import (
	"fmt"
	"os"

	"bufio"
	"path/filepath"

	_ "github.com/davecgh/go-spew/spew"
	"github.com/joernott/camunda-clean/camunda"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Connection *camunda.CamundaConnection

var rootCmd = &cobra.Command{
	Use:   "camunda-clean",
	Short: "Clean leftover camunda processes",
	Long: `Camunda-clean is a commandline tool to clean leftover camunda processes
using the REST API of the camunda engine.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := HandleConfigFile()
		if err != nil {
			panic(err)
		}
		err = InitLogging()
		if err != nil {
			fmt.Println("Error configuring logging")
			os.Exit(10)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		var Connection *camunda.CamundaConnection
		var list camunda.ProcessInstanceList
		var item camunda.ProcessInstance
		var err error
		var failed bool

		Connection = camunda.NewCamunda(viper.GetBool("ssl"),
			viper.GetString("host"),
			viper.GetInt("port"),
			viper.GetString("baseendpoint"),
			viper.GetString("user"),
			viper.GetString("password"),
			viper.GetBool("validatessl"),
			viper.GetString("proxy"),
			viper.GetBool("socks"))
		list, err = Connection.GetProcessInstanceList()
		if err != nil {
			os.Exit(20)
		}
		failed = false
		for _, item = range list {
			err = Connection.TerminateProcess(item.Id)
			if err != nil {
				failed = true
			}
		}
		if failed {
			os.Exit(21)
		}
	},
}

var ConfigFile string
var UseSSL bool
var ValidateSSL bool
var Host string
var Port int
var BaseEndpoint string
var User string
var Password string
var LogLevel int
var LogFile string
var Proxy string
var ProxyIsSocks bool

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	pwd, _ := os.Getwd()
	cfgpath := pwd + string(os.PathSeparator) + "camunda-clean.yml"
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "c", cfgpath, "Configuration file")
	rootCmd.PersistentFlags().BoolVarP(&UseSSL, "ssl", "s", false, "Use SSL")
	rootCmd.PersistentFlags().BoolVarP(&ValidateSSL, "validatessl", "v", true, "Validate SSL certificate")
	rootCmd.PersistentFlags().StringVarP(&Host, "host", "H", "localhost", "Hostname of the server")
	rootCmd.PersistentFlags().IntVarP(&Port, "port", "P", 8080, "Network port")
	rootCmd.PersistentFlags().StringVarP(&BaseEndpoint, "baseendpoint", "B", "/engine-rest", "Base endpoint for the camunda API")
	rootCmd.PersistentFlags().StringVarP(&User, "user", "u", "", "Username for Elasticsearch")
	rootCmd.PersistentFlags().StringVarP(&Password, "password", "p", "", "Password for the Elasticsearch user")
	rootCmd.PersistentFlags().IntVarP(&LogLevel, "loglevel", "l", 5, "Log level")
	rootCmd.PersistentFlags().StringVarP(&LogFile, "logfile", "L", "", "Log file (defaults to stdout)")
	rootCmd.PersistentFlags().StringVarP(&Proxy, "proxy", "y", "", "Proxy (defaults to none)")
	rootCmd.PersistentFlags().BoolVarP(&ProxyIsSocks, "socks", "Y", false, "This is a SOCKS proxy")

	viper.SetDefault("ssl", false)
	viper.SetDefault("validatessl", true)
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 8080)
	viper.SetDefault("user", "")
	viper.SetDefault("baseendpoint", "/engine-rest")
	viper.SetDefault("password", "")
	viper.SetDefault("loglevel", 5)
	viper.SetDefault("logfile", "")
	viper.SetDefault("proxy", "")
	viper.SetDefault("socks", false)

	viper.BindPFlag("ssl", rootCmd.PersistentFlags().Lookup("ssl"))
	viper.BindPFlag("validatessl", rootCmd.PersistentFlags().Lookup("validatessl"))
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("baseendpoint", rootCmd.PersistentFlags().Lookup("baseendpoint"))
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))
	viper.BindPFlag("proxy", rootCmd.PersistentFlags().Lookup("proxy"))
	viper.BindPFlag("socks", rootCmd.PersistentFlags().Lookup("socks"))
}

func HandleConfigFile() error {
	if ConfigFile != "" {
		log.Debug("Read config from " + ConfigFile)
		viper.SetConfigFile(ConfigFile)
	} else {
		log.Debug("Read config from home directory")
		home, err := homedir.Dir()
		if err != nil {
			log.Error(err)
			return err
		}
		viper.AddConfigPath(home)
		ex, err := os.Executable()
		if err != nil {
			log.Error(err)
			return err
		}
		pwd := filepath.Dir(ex)
		viper.AddConfigPath(pwd)
		viper.SetConfigName("gobana")
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Error("Can't read config: " + err.Error())
		return err
	}

	return nil
}

func InitLogging() error {
	LogFile = viper.GetString("logfile")
	LogLevel = viper.GetInt("loglevel")
	if LogFile == "" {
		log.SetOutput(os.Stdout)
	} else {
		f, err := os.Create(LogFile)
		if err != nil {
			fmt.Println("Could not create logfile '" + LogFile + "'")
			return err
		}
		w := bufio.NewWriter(f)
		log.SetOutput(w)
	}
	switch LogLevel {
	case 0:
		log.SetLevel(log.PanicLevel)
	case 1:
		log.SetLevel(log.FatalLevel)
	case 2:
		log.SetLevel(log.ErrorLevel)
	case 3:
		log.SetLevel(log.WarnLevel)
	case 4:
		log.SetLevel(log.InfoLevel)
	case 5:
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
	log.WithFields(log.Fields{
		"LogFile":  LogFile,
		"LogLevel": LogLevel,
	}).Debug("Logging configured")
	log.Debug("PersistentPreRun finished")
	return nil
}
