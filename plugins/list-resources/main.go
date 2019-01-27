package main

/*
List monitored resource descriptors
*/

import (
	"context"
	"fmt"
	"log"
	"os"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

func run(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		_ = cmd.Help()
		return fmt.Errorf("invalid argument(s) received")
	}
	ctx := context.Background()
	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	listResourceDescriptors(ctx, client)

	return nil
}

func main() {
	rootCmd := configureRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func configureRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gcpmetrics",
		Short: "The Sensu Go gcpmetrics handler for resource managemnet",
		RunE:  run,
	}
	return cmd
}

func listResourceDescriptors(ctx context.Context, client *monitoring.MetricClient) {
	request := &monitoringpb.ListMonitoredResourceDescriptorsRequest{}
	response := client.ListMonitoredResourceDescriptors(ctx, request)
	index := 0
	for {
		index++
		resource, err := response.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Errorf("Error while listing resources, %v ", err)
		}
		fmt.Print("Resource Descriptor #", index)
		fmt.Print("\tName:", resource.Name)
		fmt.Print("\tType:", resource.Type)
		fmt.Print("\tDisplay Name:", resource.DisplayName)
		fmt.Print("\tDescription:", resource.Description)
		fmt.Print("\tLabels:")
		subindex := 0
		for _, label := range resource.Labels {
			subindex++
			fmt.Print("\t\tLabel Descritor #", subindex)
			fmt.Print("\t\t\tKey:", label.Key)
			fmt.Print("\t\t\tvalue_type:", label.ValueType)
			fmt.Print("\t\t\tdescription:", label.Description)
		}
		fmt.Println("")
	}
}
