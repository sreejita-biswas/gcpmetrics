package main

/*
List available metric descriptors
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

	listMetricsDescriptors(ctx, client)

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
		Short: "The Sensu Go gcpmetrics handler for metrics descriptor management",
		RunE:  run,
	}
	return cmd
}

func listMetricsDescriptors(ctx context.Context, client *monitoring.MetricClient) {
	fmt.Print("Defined metric descriptors:")
	index := 0
	request := &monitoringpb.ListMetricDescriptorsRequest{}
	response := client.ListMetricDescriptors(ctx, request)
	for {
		index++
		metrics, err := response.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Errorf("Error while listing metrics, %v ", err)
		}
		fmt.Print("Metric Descriptor #", index)
		fmt.Print("\tName:", metrics.Name)
		fmt.Print("\tType:", metrics.Type)
		fmt.Print("\tMetric_Kind:", metrics.MetricKind)
		fmt.Print("\tValue_type:", metrics.ValueType)
		fmt.Print("\tunit:", metrics.Unit)
		fmt.Print("\tDisplay Name:", metrics.DisplayName)
		fmt.Print("\tDescription:", metrics.Description)
		fmt.Println("")
	}
}
