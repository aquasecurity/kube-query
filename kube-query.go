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
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	socketPath := flag.String("socket", "", "absolute path to the osquery socket")
	flag.Parse()
	if *kubeconfig == "" || *socketPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// initializing client and extension
	kubeclient, err := utils.CreateKubeClient(*kubeconfig)
	if err != nil {
		log.Fatalf("Error on creating kube-client: %s", err)
		panic(err)
	}
	extension, err := utils.CreateOsQueryExtension("kube-query", *socketPath)
	if err != nil {
		log.Fatalf("Error on registering osquery extension: %s", err)
		panic(err)
	}

	// creating tables and appending to list
	tableList := make([]tables.Table, 3)
	tableList[0] = tables.NewPodsTable(kubeclient)
	tableList[1] = tables.NewContainersTable(kubeclient)
	tableList[2] = tables.NewVolumesTable(kubeclient)

	// Registering all tables
	for _, t := range tableList {
		// Create and register a new table plugin with the server.
		extension.RegisterPlugin(osqueryTable.NewPlugin(t.Name(), t.Columns(), t.Generate))
	}

	if err := extension.Run(); err != nil {
		log.Fatalf("Error in registering tables: %s", err)
	}
}
