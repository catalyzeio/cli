package metrics

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
)

// CSVTransformer is a concrete implementation of Transformer transforming data
// into CSV format.
type CSVTransformer struct {
	HeadersWritten bool
	GroupMode      bool
	Buffer         *bytes.Buffer
	Writer         *csv.Writer
}

// WriteHeadersCPU outputs the csv headers needed for cpu data. If GroupMode
// is enabled, the service name is the first header.
func (csv *CSVTransformer) WriteHeadersCPU() {
	if !csv.HeadersWritten {
		headers := []string{"timestamp", "cpu_min", "cpu_max", "cpu_avg", "cpu_total"}
		if csv.GroupMode {
			headers = append([]string{"service_name"}, headers...)
		}
		csv.Writer.Write(headers)
		csv.HeadersWritten = true
	}
}

// WriteHeadersMemory outputs the csv headers needed for memory data. If
// GroupMode is enabled, the service name is the first header.
func (csv *CSVTransformer) WriteHeadersMemory() {
	if !csv.HeadersWritten {
		headers := []string{"timestamp", "memory_min", "memory_max", "memory_avg", "memory_total"}
		if csv.GroupMode {
			headers = append([]string{"service_name"}, headers...)
		}
		csv.Writer.Write(headers)
		csv.HeadersWritten = true
	}
}

// WriteHeadersNetworkIn outputs the csv headers needed for received network
// data. If GroupMode is enabled, the service name is the first header.
func (csv *CSVTransformer) WriteHeadersNetworkIn() {
	if !csv.HeadersWritten {
		headers := []string{"timestamp", "rx_bytes", "rx_packets"}
		if csv.GroupMode {
			headers = append([]string{"service_name"}, headers...)
		}
		csv.Writer.Write(headers)
		csv.HeadersWritten = true
	}
}

// WriteHeadersNetworkOut outputs the csv headers needed for transmitted network
// data. If GroupMode is enabled, the service name is the first header.
func (csv *CSVTransformer) WriteHeadersNetworkOut() {
	if !csv.HeadersWritten {
		headers := []string{"timestamp", "tx_bytes", "tx_packets"}
		if csv.GroupMode {
			headers = append([]string{"service_name"}, headers...)
		}
		csv.Writer.Write(headers)
		csv.HeadersWritten = true
	}
}

// TransformGroupCPU transforms an entire environment's cpu data into csv
// format. This outputs TransformSingleCPU for each service in the environment.
func (csv *CSVTransformer) TransformGroupCPU(metrics *[]models.Metrics) {
	csv.GroupMode = true
	for _, metric := range *metrics {
		if _, ok := blacklist[metric.ServiceLabel]; !ok {
			csv.TransformSingleCPU(&metric)
		}
	}
	csv.Writer.Flush()
	logrus.Println(csv.Buffer.String())
}

// TransformGroupMemory transforms an entire environment's memory data into csv
// format. This outputs TransformSingleMemory for each service in the
// environment.
func (csv *CSVTransformer) TransformGroupMemory(metrics *[]models.Metrics) {
	csv.GroupMode = true
	for _, metric := range *metrics {
		if _, ok := blacklist[metric.ServiceLabel]; !ok {
			csv.TransformSingleMemory(&metric)
		}
	}
	csv.Writer.Flush()
	logrus.Println(csv.Buffer.String())
}

// TransformGroupNetworkIn transforms an entire environment's received network
// data into csv format. This outputs TransformSingleNetworkIn for each service
// in the environment.
func (csv *CSVTransformer) TransformGroupNetworkIn(metrics *[]models.Metrics) {
	csv.GroupMode = true
	for _, metric := range *metrics {
		if _, ok := blacklist[metric.ServiceLabel]; !ok {
			csv.TransformSingleNetworkIn(&metric)
		}
	}
	csv.Writer.Flush()
	logrus.Println(csv.Buffer.String())
}

// TransformGroupNetworkOut transforms an entire environment's transmitted
// network data into csv format. This outputs TransformSingleNetworkOut for
// each service in the environment.
func (csv *CSVTransformer) TransformGroupNetworkOut(metrics *[]models.Metrics) {
	csv.GroupMode = true
	for _, metric := range *metrics {
		if _, ok := blacklist[metric.ServiceLabel]; !ok {
			csv.TransformSingleNetworkOut(&metric)
		}
	}
	csv.Writer.Flush()
	logrus.Println(csv.Buffer.String())
}

// TransformSingleCPU transforms a single service's CPU data into csv format.
func (csv *CSVTransformer) TransformSingleCPU(metric *models.Metrics) {
	csv.WriteHeadersCPU()
	if metric.Data != nil && metric.Data.CPUUsage != nil {
		for _, data := range *metric.Data.CPUUsage {
			row := []string{
				fmt.Sprintf("%d", data.TS),
				fmt.Sprintf("%f", data.Min/1000.0),
				fmt.Sprintf("%f", data.Max/1000.0),
				fmt.Sprintf("%f", data.AVG/1000.0),
				fmt.Sprintf("%f", data.Total/1000.0),
			}
			if csv.GroupMode {
				row = append([]string{metric.ServiceLabel}, row...)
			}
			csv.Writer.Write(row)
		}
	}
	if !csv.GroupMode {
		csv.Writer.Flush()
		logrus.Println(csv.Buffer.String())
	}
}

// TransformSingleMemory transforms a single service's memory data into csv
// format.
func (csv *CSVTransformer) TransformSingleMemory(metric *models.Metrics) {
	csv.WriteHeadersMemory()
	if metric.Data != nil && metric.Data.MemoryUsage != nil {
		for _, data := range *metric.Data.MemoryUsage {
			row := []string{
				fmt.Sprintf("%d", data.TS),
				fmt.Sprintf("%f", data.Min/1024.0),
				fmt.Sprintf("%f", data.Max/1024.0),
				fmt.Sprintf("%f", data.AVG/1024.0),
				fmt.Sprintf("%f", float64(metric.Size.RAM*1024.0)),
			}
			if csv.GroupMode {
				row = append([]string{metric.ServiceLabel}, row...)
			}
			csv.Writer.Write(row)
		}
	}
	if !csv.GroupMode {
		csv.Writer.Flush()
		logrus.Println(csv.Buffer.String())
	}
}

// TransformSingleNetworkIn transforms a single service's received network data
// into csv format.
func (csv *CSVTransformer) TransformSingleNetworkIn(metric *models.Metrics) {
	csv.WriteHeadersNetworkIn()
	if metric.Data != nil && metric.Data.NetworkUsage != nil {
		for _, data := range *metric.Data.NetworkUsage {
			row := []string{
				fmt.Sprintf("%d", data.TS),
				fmt.Sprintf("%f", data.RXKB),
				fmt.Sprintf("%f", data.RXPackets),
			}
			if csv.GroupMode {
				row = append([]string{metric.ServiceLabel}, row...)
			}
			csv.Writer.Write(row)
		}
	}
	if !csv.GroupMode {
		csv.Writer.Flush()
		logrus.Println(csv.Buffer.String())
	}
}

// TransformSingleNetworkOut transforms a single service's transmitted network
// data into csv format.
func (csv *CSVTransformer) TransformSingleNetworkOut(metric *models.Metrics) {
	csv.WriteHeadersNetworkOut()
	if metric.Data != nil && metric.Data.NetworkUsage != nil {
		for _, data := range *metric.Data.NetworkUsage {
			row := []string{
				fmt.Sprintf("%d", data.TS),
				fmt.Sprintf("%f", data.TXKB),
				fmt.Sprintf("%f", data.TXPackets),
			}
			if csv.GroupMode {
				row = append([]string{metric.ServiceLabel}, row...)
			}
			csv.Writer.Write(row)
		}
	}
	if !csv.GroupMode {
		csv.Writer.Flush()
		logrus.Println(csv.Buffer.String())
	}
}
