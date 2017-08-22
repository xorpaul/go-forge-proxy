package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var (
	debug          bool
	verbose        bool
	info           bool
	quiet          bool
	force          bool
	usemove        bool
	pfMode         bool
	pfLocation     string
	dryRun         bool
	check4update   bool
	checkSum       bool
	moduleDirParam string
	cacheDirParam  string
	branchParam    string
	configFile     string
	config         ConfigSettings
	wg             sync.WaitGroup
	mutex          sync.Mutex
	empty          struct{}
	buildtime      string
	maxworker      int
	syncForgeTime  float64

	forgeJsonParseTime float64
)

// ConfigSettings contains the key value pairs from the config file
type ConfigSettings struct {
	Port      int
	CacheDir  string `yaml:"cachedir"`
	Timeout   int    `yaml:"timeout"`
	Maxworker int
	ForgeUrl  string `yaml:forgeurl"`
}

func main() {

	var (
		configFileFlag = flag.String("config", "", "which config file to use")
		versionFlag    = flag.Bool("version", false, "show build time and version number")
	)
	flag.BoolVar(&debug, "debug", false, "log debug output, defaults to false")
	flag.BoolVar(&verbose, "verbose", false, "log verbose output, defaults to false")
	flag.BoolVar(&info, "info", false, "log info output, defaults to false")
	flag.BoolVar(&quiet, "quiet", false, "no output, defaults to false")

	flag.Parse()

	configFile = *configFileFlag
	version := *versionFlag

	if version {
		fmt.Println("go-forge-proxy Build time:", buildtime, "UTC")
		os.Exit(0)
	}

	if len(configFile) > 0 {
		config = readConfigfile(configFile)
		http.HandleFunc("/", handleRequest)
		Verbosef("Starting Forge mirror on :" + strconv.Itoa(config.Port))
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))

	} else {
		Fatalf("Error: you need to specify a config file!\nExample call: " + os.Args[0] + " -config test.yaml\n")
	}

}
