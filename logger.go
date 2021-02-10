package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/jcxplorer/cwlogger"
)

var (
	logger *cwlogger.Logger
	cw     *cloudwatch.CloudWatch
)

// SendLog sends the report as a log entry to AWS
func (report Report) SendLog() {
	if len(os.Getenv("RE_LOG_GROUP")) == 0 {
		log.Fatal("no log group defined (did you specify env RE_LOG_GROUP?)")
	}

	var err error

	if logger == nil {
		logger, err = cwlogger.New(&cwlogger.Config{
			LogGroupName: os.Getenv("RE_LOG_GROUP"),
			Client: cloudwatchlogs.New(session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))),
		})

		if err != nil {
			panic(fmt.Errorf("error connecting to CloudWatch: %v", err))
		}
	}

	logger.Log(time.Now(), report.ToString())
}

// PublishMetric publishes the collected statistics as a metric to AWS
func (report Report) PublishMetric() {
	if len(os.Getenv("RE_NAMESPACE")) == 0 {
		log.Fatal("no namespace defined (did you specify env RE_NAMESPACE?)")
	}

	if cw == nil {
		cw = cloudwatch.New(session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		})))
	}

	_, err := cw.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace: aws.String(os.Getenv("RE_NAMESPACE") + "/" + report.Hostname),
		MetricData: []*cloudwatch.MetricDatum{
			{
				MetricName: aws.String("CPU utilization"),
				Timestamp:  &report.Timestamp,
				Value:      aws.Float64(float64(report.CPULoad * 100.0)),
				Unit:       aws.String("Percent"),
			},
			{
				MetricName: aws.String("Uptime"),
				Timestamp:  &report.Timestamp,
				Value:      aws.Float64(float64(report.Uptime)),
				Unit:       aws.String("Seconds"),
			},
			{
				MetricName: aws.String("Memory utilization (MB)"),
				Timestamp:  &report.Timestamp,
				Value:      aws.Float64(float64(report.MemoryUsed)),
				Unit:       aws.String("Megabytes"),
			},
			{
				MetricName: aws.String("Memory utilization (%)"),
				Timestamp:  &report.Timestamp,
				Value:      aws.Float64((float64(report.MemoryUsed) / float64(report.MemoryTotal)) * 100.0),
				Unit:       aws.String("Percent"),
			},
			{
				MetricName: aws.String("Swap utilization (MB)"),
				Timestamp:  &report.Timestamp,
				Value:      aws.Float64(float64(report.SwapUsed)),
				Unit:       aws.String("Megabytes"),
			},
			{
				MetricName: aws.String("Swap utilization (%)"),
				Timestamp:  &report.Timestamp,
				Value:      aws.Float64((float64(report.SwapUsed) / float64(report.SwapTotal)) * 100.0),
				Unit:       aws.String("Percent"),
			},
			{
				MetricName: aws.String("Load average"),
				Timestamp:  &report.Timestamp,
				Value:      &report.LoadAverage[0],
				Unit:       aws.String("Count"),
			},
		},
	})

	if err != nil {
		log.Fatalf("error connecting to Cloudwatch: %v", err)
	}
}
