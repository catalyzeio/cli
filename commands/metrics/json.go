package metrics

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
)

// JSONTransformer is a concrete implementation of Transformer transforming data
// into pretty printed JSON format.
type JSONTransformer struct{}

type cpu struct {
	ServiceName string  `json:"service_name,omitempty"`
	TS          int     `json:"ts"`
	Min         float64 `json:"min"`
	Max         float64 `json:"max"`
	AVG         float64 `json:"avg"`
	Total       float64 `json:"total"`
}

type mem struct {
	ServiceName string  `json:"service_name,omitempty"`
	TS          int     `json:"ts"`
	Min         float64 `json:"min"`
	Max         float64 `json:"max"`
	AVG         float64 `json:"avg"`
	Total       float64 `json:"total"`
}

type netin struct {
	ServiceName string  `json:"service_name,omitempty"`
	TS          int     `json:"ts"`
	RXKB        float64 `json:"rx_kb"`
	RXPackets   float64 `json:"rx_packets"`
}

type netout struct {
	ServiceName string  `json:"service_name,omitempty"`
	TS          int     `json:"ts"`
	TXKB        float64 `json:"tx_kb"`
	TXPackets   float64 `json:"tx_packets"`
}

// TransformGroupCPU transforms an entire environment's cpu data into json
// format. This outputs TransformSingleCPU for every service in the environment.
func (j *JSONTransformer) TransformGroupCPU(metrics *[]models.Metrics) {
	var data []cpu
	for _, m := range *metrics {
		if _, ok := blacklist[m.ServiceLabel]; !ok && m.Data != nil && m.Data.CPUUsage != nil {
			for _, d := range *m.Data.CPUUsage {
				data = append(data, cpu{m.ServiceLabel, d.TS, d.Min / 1000.0, d.Max / 1000.0, d.AVG / 1000.0, d.Total / 1000.0})
			}
		}
	}
	b, _ := json.MarshalIndent(data, "", "    ")
	logrus.Println(string(b))
}

// TransformGroupMemory transforms an entire environment's memory data into json
// format. This outputs TransformSingleMemory for every service in the
// environment.
func (j *JSONTransformer) TransformGroupMemory(metrics *[]models.Metrics) {
	var data []mem
	for _, m := range *metrics {
		if _, ok := blacklist[m.ServiceLabel]; !ok && m.Data != nil && m.Data.MemoryUsage != nil {
			for _, d := range *m.Data.MemoryUsage {
				data = append(data, mem{m.ServiceLabel, d.TS, d.Min / 1024.0, d.Max / 1024.0, d.AVG / 1024.0, float64(m.Size.RAM) * 1024.0})
			}
		}
	}
	b, _ := json.MarshalIndent(data, "", "    ")
	logrus.Println(string(b))
}

// TransformGroupNetworkIn transforms an entire environment's received network
// data into json format. This outputs TransformSingleNetworkIn for every
// service in the environment.
func (j *JSONTransformer) TransformGroupNetworkIn(metrics *[]models.Metrics) {
	var data []netin
	for _, m := range *metrics {
		if _, ok := blacklist[m.ServiceLabel]; !ok && m.Data != nil && m.Data.NetworkUsage != nil {
			for _, d := range *m.Data.NetworkUsage {
				data = append(data, netin{m.ServiceLabel, d.TS, d.RXKB, d.RXPackets})
			}
		}
	}
	b, _ := json.MarshalIndent(data, "", "    ")
	logrus.Println(string(b))
}

// TransformGroupNetworkOut transforms an entire environment's transmitted
// network data into json format. This outputs TransformSingleNetworkOut for
// every service in the environment.
func (j *JSONTransformer) TransformGroupNetworkOut(metrics *[]models.Metrics) {
	var data []netout
	for _, m := range *metrics {
		if _, ok := blacklist[m.ServiceLabel]; !ok && m.Data != nil && m.Data.NetworkUsage != nil {
			for _, d := range *m.Data.NetworkUsage {
				data = append(data, netout{m.ServiceLabel, d.TS, d.TXKB, d.TXPackets})
			}
		}
	}
	b, _ := json.MarshalIndent(data, "", "    ")
	logrus.Println(string(b))
}

// TransformSingleCPU transforms a single service's cpu data into json format.
func (j *JSONTransformer) TransformSingleCPU(metric *models.Metrics) {
	var data []cpu
	if metric.Data != nil && metric.Data.CPUUsage != nil {
		for _, d := range *metric.Data.CPUUsage {
			data = append(data, cpu{TS: d.TS, Min: d.Min / 1000.0, Max: d.Max / 1000.0, AVG: d.AVG / 1000.0, Total: d.Total / 1000.0})
		}
	}
	b, _ := json.MarshalIndent(data, "", "    ")
	logrus.Println(string(b))
}

// TransformSingleMemory transforms a single service's memory data into json
// format.
func (j *JSONTransformer) TransformSingleMemory(metric *models.Metrics) {
	var data []mem
	if metric.Data != nil && metric.Data.MemoryUsage != nil {
		for _, d := range *metric.Data.MemoryUsage {
			data = append(data, mem{TS: d.TS, Min: d.Min / 1024.0, Max: d.Max / 1024.0, AVG: d.AVG / 1024.0, Total: float64(metric.Size.RAM) * 1024.0})
		}
	}
	b, _ := json.MarshalIndent(data, "", "    ")
	logrus.Println(string(b))
}

// TransformSingleNetworkIn transforms a single service's received network data
// into json format.
func (j *JSONTransformer) TransformSingleNetworkIn(metric *models.Metrics) {
	var data []netin
	if metric.Data != nil && metric.Data.NetworkUsage != nil {
		for _, d := range *metric.Data.NetworkUsage {
			data = append(data, netin{TS: d.TS, RXKB: d.RXKB, RXPackets: d.RXPackets})
		}
	}
	b, _ := json.MarshalIndent(data, "", "    ")
	logrus.Println(string(b))
}

// TransformSingleNetworkOut transforms a single service's transmitted network
// data into json format.
func (j *JSONTransformer) TransformSingleNetworkOut(metric *models.Metrics) {
	var data []netout
	if metric.Data != nil && metric.Data.NetworkUsage != nil {
		for _, d := range *metric.Data.NetworkUsage {
			data = append(data, netout{TS: d.TS, TXKB: d.TXKB, TXPackets: d.TXPackets})
		}
	}
	b, _ := json.MarshalIndent(data, "", "    ")
	logrus.Println(string(b))
}
