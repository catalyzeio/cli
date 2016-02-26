package metrics

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
)

// TextTransformer is a concrete implementation of Transformer transforming data
// into plain text.
type TextTransformer struct{}

// TransformGroupCPU transforms an entire environment's cpu data into text
// format. This outputs TransformSingleCPU for every service in the environment.
func (text *TextTransformer) TransformGroupCPU(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		if _, ok := blacklist[metric.ServiceLabel]; !ok {
			logrus.Printf("%s:", metric.ServiceLabel)
			text.TransformSingleCPU(&metric)
		}
	}
}

// TransformGroupMemory transforms an entire environment's memory data into
// text format. This outputs TransformSingleMemory for every service in the
// environment.
func (text *TextTransformer) TransformGroupMemory(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		if _, ok := blacklist[metric.ServiceLabel]; !ok {
			logrus.Printf("%s:", metric.ServiceLabel)
			text.TransformSingleMemory(&metric)
		}
	}
}

// TransformGroupNetworkIn transforms an entire environment's received network
// data into text format. This outputs TransformSingleNetworkIn for every
// service in the environment.
func (text *TextTransformer) TransformGroupNetworkIn(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		if _, ok := blacklist[metric.ServiceLabel]; !ok {
			logrus.Printf("%s:", metric.ServiceLabel)
			text.TransformSingleNetworkIn(&metric)
		}
	}
}

// TransformGroupNetworkOut transforms an entire environment's transmitted
// network data into text format. This outputs TransformSingleNetworkOut for
// every service in the environment.
func (text *TextTransformer) TransformGroupNetworkOut(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		if _, ok := blacklist[metric.ServiceLabel]; !ok {
			logrus.Printf("%s:", metric.ServiceLabel)
			text.TransformSingleNetworkOut(&metric)
		}
	}
}

// TransformSingleCPU transforms a single service's cpu data into text format.
func (text *TextTransformer) TransformSingleCPU(metric *models.Metrics) {
	prefix := "    "
	if metric.Data != nil && metric.Data.CPUUsage != nil {
		for _, data := range *metric.Data.CPUUsage {
			ts := time.Unix(int64(data.TS/1000.0), 0)
			logrus.Printf("%s%s | CPU Min: %6.2f | CPU Max: %6.2f | CPU AVG: %6.2f | CPU Total: %6.2f",
				prefix,
				fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second()),
				data.Min/1000.0,
				data.Max/1000.0,
				data.AVG/1000.0,
				data.Total/1000.0)
		}
	}
}

// TransformSingleMemory transforms a single service's memory data into text
// format.
func (text *TextTransformer) TransformSingleMemory(metric *models.Metrics) {
	prefix := "    "
	if metric.Data != nil && metric.Data.MemoryUsage != nil {
		for _, data := range *metric.Data.MemoryUsage {
			ts := time.Unix(int64(data.TS/1000.0), 0)
			logrus.Printf("%s%s | Memory Min: %.2f MB | Memory Max: %.2f MB | Memory AVG: %.2f MB | Memory Total: %.2f MB",
				prefix,
				fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second()),
				data.Min/1024.0,
				data.Max/1024.0,
				data.AVG/1024.0,
				float64(metric.Size.RAM)*1024.0)
		}
	}
}

// TransformSingleNetworkIn transforms a single service's received network data
// into text format.
func (text *TextTransformer) TransformSingleNetworkIn(metric *models.Metrics) {
	prefix := "    "
	if metric.Data != nil && metric.Data.NetworkUsage != nil {
		for _, data := range *metric.Data.NetworkUsage {
			ts := time.Unix(int64(data.TS/1000.0), 0)
			logrus.Printf("%s%s | Received Bytes: %.2f KB | Received Packets: %.2f",
				prefix,
				fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second()),
				data.RXKB,
				data.RXPackets)
		}
	}
}

// TransformSingleNetworkOut transforms a single service's transmitted network
// data into text format.
func (text *TextTransformer) TransformSingleNetworkOut(metric *models.Metrics) {
	prefix := "    "
	if metric.Data != nil && metric.Data.NetworkUsage != nil {
		for _, data := range *metric.Data.NetworkUsage {
			ts := time.Unix(int64(data.TS/1000.0), 0)
			logrus.Printf("%s%s | Transmitted Bytes: %.2f KB | Transmitted Packets: %.2f",
				prefix,
				fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second()),
				data.TXKB,
				data.TXPackets)
		}
	}
}
