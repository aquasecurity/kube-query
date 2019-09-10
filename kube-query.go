package main

import (
	"flag"
	"log"
	"os"

	tables "github.com/aquasecurity/kube-query/tables"
	utils "github.com/aquasecurity/kube-query/utils"
	osqueryTable "github.com/kolide/osquery-go/plugin/table"
)

func main() {
	// Parsing flags
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file (can be set by KUBECONFIG environment variable)")
	socketPath := flag.String("socket", "", "absolute path to the osquery socket")
	
	// currently we do not care for these flags, but they must be set for the auto loader of osquery
	flag.String("timeout", "", "flag for specifying wait time before registering on autoload") 
	flag.String("interval", "", "flag for specifying wait time before registering on autoload")
	
	flag.Parse()
	if len(*kubeconfig) == 0 {
		// if not specified from flag, try getting from env variable
		if *kubeconfig = os.Getenv("KUBECONFIG"); len(*kubeconfig) == 0 {
			log.Fatal("Kubeconfig was not specified. set KUBECONFIG environment variable or pass the --kubeconfig flag")
			os.Exit(1)
		}
	}
	if len(*socketPath) == 0 {
		log.Fatal("Socket was not specified, set the --socket flag")		
		os.Exit(1)
	}

	// initializing clients and extension
	kubeclient, err := utils.CreateKubeClient(*kubeconfig)
	if err != nil {
		log.Fatalf("Error on creating kube-client: %s", err)
		panic(err)
	}
	metricsclient, err := utils.CreateMetricsClient(*kubeconfig)
	if err != nil {
		log.Fatalf("Error on creating the metrics client: %s", err)
		panic(err)
	}
	extension, err := utils.CreateOsQueryExtension("kube-query", *socketPath)
	if err != nil {
		log.Fatalf("Error on registering osquery extension: %s", err)
		panic(err)
	}

	// creating tables and appending to list
	tableList := []tables.Table{
		tables.NewPodsTable(kubeclient),
		tables.NewContainersTable(kubeclient),
		tables.NewVolumesTable(kubeclient),
		tables.NewNodesTable(kubeclient, metricsclient), // specific columns use the metrics client
		tables.NewDeploymentsTable(kubeclient),
	}

	// Registering all tables
	for _, t := range tableList {
		// Create and register a new table plugin with the server.
		extension.RegisterPlugin(osqueryTable.NewPlugin(t.Name(), t.Columns(), t.Generate))
	}

	if err := extension.Run(); err != nil {
		log.Fatalf("Error in registering tables: %s", err)
	}
}
