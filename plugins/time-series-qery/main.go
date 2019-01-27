package main

/*
Run the time series query
*/

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

var (
	metricType string
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

	queryTimeSeries(ctx, client)
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
		Short: "The Sensu Go gcpmetrics handler for gcpmetrics management",
		RunE:  run,
	}

	cmd.Flags().StringVar(&metricType, "metric_type", "custom.googleapis.com/custom_measurement", "Metric ID as defined by Google Monitoring API.")
	cmd.MarkFlagRequired("metric_type")
	return cmd
}

func queryTimeSeries(ctx context.Context, client *monitoring.MetricClient) error {
	startTime := time.Now().UTC().Add(time.Minute * -5).Unix()
	endTime := time.Now().UTC().Unix()

	req := &monitoringpb.ListTimeSeriesRequest{
		Filter: fmt.Sprintf("metric.type=\"%s\"", metricType),
		Interval: &monitoringpb.TimeInterval{
			StartTime: &timestamp.Timestamp{Seconds: startTime},
			EndTime:   &timestamp.Timestamp{Seconds: endTime},
		},
	}

	iter := client.ListTimeSeries(ctx, req)

	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("could not read time series value, %v ", err)
		}
		fmt.Println(resp)
	}

	return nil
}
