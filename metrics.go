package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2/google"

	"cloud.google.com/go/compute/metadata"
	"google.golang.org/api/monitoring/v3"
)

func projectResource(projectID string) string {
	return "projects/" + projectID
}

func customMetricType(metricType string) string {
	return "custom.googleapis.com/varnish." + metricType
}

func createService(ctx context.Context) (*monitoring.Service, error) {
	hc, err := google.DefaultClient(ctx, monitoring.MonitoringScope)
	if err != nil {
		return nil, err
	}
	s, err := monitoring.New(hc)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func createCustomMetric(s *monitoring.Service, projectID, metricType string) error {
	ld := monitoring.LabelDescriptor{Key: "application", ValueType: "STRING", Description: "Application that this varnish caches"}
	md := monitoring.MetricDescriptor{
		Type:        customMetricType(metricType),
		Labels:      []*monitoring.LabelDescriptor{&ld},
		MetricKind:  "GAUGE",
		ValueType:   "INT64",
		Unit:        "items",
		Description: "Varnishstat " + metricType + " per min",
		DisplayName: metricType,
	}
	_, err := s.Projects.MetricDescriptors.Create(projectResource(projectID), &md).Do()
	if err != nil {
		return fmt.Errorf("Could not create custom metric: %v", err)
	}
	return nil
}

// getCustomMetric reads the custom metric created.
// func getCustomMetric(s *monitoring.Service, projectID, metricType string) (*monitoring.ListMetricDescriptorsResponse, error) {
// 	resp, err := s.Projects.MetricDescriptors.List(projectResource(projectID)).
// 		Filter(fmt.Sprintf("metric.type=\"%s\"", metricType)).Do()
// 	if err != nil {
// 		return nil, fmt.Errorf("Could not get custom metric: %v", err)
// 	}
// 	return resp, nil
// }


func writeTimeSeriesValue(s *monitoring.Service, projectID, metricType string, application string, value int64) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	zone, _ := metadata.Zone()
	instanceID, _ := metadata.InstanceID()

	timeseries := monitoring.TimeSeries{
		Metric: &monitoring.Metric{
			Type: customMetricType(metricType),
			Labels: map[string]string{
				"application": application,
			},
		},
		Resource: &monitoring.MonitoredResource{
			Labels: map[string]string{
				"instance_id": instanceID,
				"zone":        zone,
			},
			Type: "gce_instance",
		},
		Points: []*monitoring.Point{
			{
				Interval: &monitoring.TimeInterval{
					StartTime: now,
					EndTime:   now,
				},
				Value: &monitoring.TypedValue{
					Int64Value: &value,
				},
			},
		},
	}

	createTimeseriesRequest := monitoring.CreateTimeSeriesRequest{
		TimeSeries: []*monitoring.TimeSeries{&timeseries},
	}
	_, err := s.Projects.TimeSeries.Create(projectResource(projectID), &createTimeseriesRequest).Do()
	if err != nil {
		return fmt.Errorf("Could not write time series value, %v ", err)
	}
	return nil
}
