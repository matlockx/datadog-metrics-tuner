package main

import (
	"flag"
	"fmt"
	"github.com/zorkian/go-datadog-api"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"time"
)

var gitHash = "No hash provided"
var buildTime = "No build time provided"

var apiKey = flag.String("api-key", "", "Datadog api key")
var appKey = flag.String("app-key", "", "Datadog app key")
var logIntervalInSeconds = flag.Int64("interval", 10, "Log interval in seconds.")
var configDirPath string

func init() {
	const (
		helpUsage = "Directory to metric files."
	)
	flag.StringVar(&configDirPath, "config-dir", "./metrics.d", helpUsage)
	flag.StringVar(&configDirPath, "c", "./metrics.d", helpUsage+" (shorthand)")
}

type metric struct {
	Name  string
	Value string
	Tags  []string
}

func main() {
	flag.Parse()

	log.Println("Git Commit Hash: ", gitHash)
	log.Println("UTC Build Time : ", buildTime)
	log.Println("Using API KEY ", *apiKey)
	log.Println("Using APP KEY ", *appKey)

	configuredMetrics := readMetricsFromFiles()
	datadogClient := datadog.NewClient(*apiKey, *appKey)

	ddmetrics := createDDMetrics(configuredMetrics)
	for {
		updateTimstampInDDMetrics(ddmetrics)
		postMetrics(datadogClient, ddmetrics)
		time.Sleep(time.Duration(*logIntervalInSeconds) * time.Second)
	}

}

func readMetricsFromFiles() []metric {
	log.Println("Parsing files from directory ", configDirPath)
	configDir, _ := filepath.Abs(configDirPath)
	files, err := ioutil.ReadDir(configDir)
	if err != nil {
		log.Fatalln("Dir [", configDirPath, "] not found: ", err)
	}
	configuredMetrics := []metric{}
	for _, file := range files {
		fileMetrics := []metric{}
		filePath := fmt.Sprintf("%s/%s", configDir, file.Name())
		fileContent, _ := ioutil.ReadFile(filePath)
		if err := yaml.Unmarshal(fileContent, &fileMetrics); err != nil {
			log.Fatalln("Unmarshalling error in ", filePath, ": ", err)
		}
		configuredMetrics = append(configuredMetrics, fileMetrics...)
	}
	return configuredMetrics
}

func updateTimstampInDDMetrics(ddmetrics []datadog.Metric) {
	timestamp := time.Now().Unix()
	for _, metric := range ddmetrics {
		for i := range metric.Points {
			metric.Points[i][0] = float64(timestamp)
		}
	}
}

func createDDMetrics(configuredMetrics []metric) []datadog.Metric {
	ddmetrics := []datadog.Metric{}
	timestamp := time.Now().Unix()
	for _, metric := range configuredMetrics {
		value, _ := strconv.ParseFloat(metric.Value, 64)
		ddmetric := datadog.Metric{
			Metric: metric.Name,
			Points: []datadog.DataPoint{
				{float64(timestamp), value},
			},
			Tags: metric.Tags,
		}
		ddmetrics = append(ddmetrics, ddmetric)
		createdMetric := fmt.Sprintf("%#v", ddmetric)
		log.Println("Will report the following metric: ", createdMetric)
	}
	return ddmetrics
}

func postMetrics(datadogClient *datadog.Client, metrics []datadog.Metric) {
	log.Println("Updating", len(metrics), "metrics.")
	if err := datadogClient.PostMetrics(metrics); err != nil {
		log.Println("ERROR: Cannot log to datadog: ", err)
	}
}
