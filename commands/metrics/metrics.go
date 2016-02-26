package metrics

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
	ui "github.com/gizak/termui"
)

var blacklist = map[string]struct{}{
	"logging":       struct{}{},
	"service_proxy": struct{}{},
	"monitoring":    struct{}{},
}

// Transformer specifies that all concrete implementations should be
// able to transform an entire environments metrics data (group) or a single
// service metrics data (single).
type Transformer interface {
	TransformGroupCPU(*[]models.Metrics)
	TransformGroupMemory(*[]models.Metrics)
	TransformGroupNetworkIn(*[]models.Metrics)
	TransformGroupNetworkOut(*[]models.Metrics)
	TransformSingleCPU(*models.Metrics)
	TransformSingleMemory(*models.Metrics)
	TransformSingleNetworkIn(*models.Metrics)
	TransformSingleNetworkOut(*models.Metrics)
}

// CmdMetrics prints out metrics for a given service or if the service is not
// specified, metrics for the entire environment are printed.
func CmdMetrics(svcName string, metricType MetricType, jsonFlag, csvFlag, sparkFlag, streamFlag bool, mins int, im IMetrics, is services.IServices) error {
	if streamFlag && (jsonFlag || csvFlag || mins != 1) {
		return fmt.Errorf("--stream cannot be used with a custom format and multiple records")
	}
	if mins > 1440 {
		return fmt.Errorf("--mins cannot be greater than 1440")
	}
	var mt Transformer
	if jsonFlag {
		mt = &JSONTransformer{}
	} else if csvFlag {
		buffer := &bytes.Buffer{}
		mt = &CSVTransformer{
			HeadersWritten: false,
			GroupMode:      false,
			Buffer:         buffer,
			Writer:         csv.NewWriter(buffer),
		}
	} else if sparkFlag {
		// the spark lines interface stays up until closed by the user, so
		// we might as well keep updating it as long as it is there
		streamFlag = true
		mins = 30
		err := ui.Init()
		if err != nil {
			return err
		}
		defer ui.Close()

		p := ui.NewPar("PRESS q TO QUIT")
		p.Border = false

		p2 := ui.NewPar(fmt.Sprintf("%s Usage Metrics", metricsTypeToString(metricType)))
		p2.Border = false

		ui.Body.AddRows(
			ui.NewRow(ui.NewCol(12, 0, p)),
			ui.NewRow(ui.NewCol(12, 0, p2)),
		)
		ui.Body.Align()
		ui.Render(ui.Body)

		mt = &SparkTransformer{
			SparkLines: map[string]*ui.Sparklines{},
		}
	} else {
		mt = &TextTransformer{}
	}
	if svcName != "" {
		service, err := is.RetrieveByLabel(svcName)
		if err != nil {
			return err
		}
		if service == nil {
			return fmt.Errorf("Could not find a service with the label \"%s\"", svcName)
		}
		return CmdServiceMetrics(metricType, streamFlag, sparkFlag, mins, service, mt, im)
	}
	return CmdEnvironmentMetrics(metricType, streamFlag, sparkFlag, mins, mt, im)
}

func CmdEnvironmentMetrics(metricType MetricType, stream, sparkLines bool, mins int, t Transformer, im IMetrics) error {
	done := make(chan struct{})
	go func() {
		for {
			metrics, err := im.RetrieveEnvironmentMetrics(mins)
			if err != nil {
				logrus.Fatal(err.Error())
			}
			switch metricType {
			case CPU:
				t.TransformGroupCPU(metrics)
			case Memory:
				t.TransformGroupMemory(metrics)
			case NetworkIn:
				t.TransformGroupNetworkIn(metrics)
			case NetworkOut:
				t.TransformGroupNetworkOut(metrics)
			}
			if !stream {
				break
			}
			time.Sleep(time.Minute)
		}
		done <- struct{}{}
	}()
	if sparkLines {
		sparkLinesEventLoop()
	} else {
		<-done
	}
	return nil
}

func CmdServiceMetrics(metricType MetricType, stream, sparkLines bool, mins int, service *models.Service, t Transformer, im IMetrics) error {
	done := make(chan struct{})
	go func() {
		for {
			metrics, err := im.RetrieveServiceMetrics(mins, service.ID)
			if err != nil {
				logrus.Fatal(err.Error())
			}
			switch metricType {
			case CPU:
				t.TransformSingleCPU(metrics)
			case Memory:
				t.TransformSingleMemory(metrics)
			case NetworkIn:
				t.TransformSingleNetworkIn(metrics)
			case NetworkOut:
				t.TransformSingleNetworkOut(metrics)
			}
			if !stream {
				break
			}
			time.Sleep(time.Minute)
		}
		done <- struct{}{}
	}()
	if sparkLines {
		sparkLinesEventLoop()
	} else {
		<-done
	}
	return nil
}

func metricsTypeToString(metricType MetricType) string {
	switch metricType {
	case CPU:
		return "CPU"
	case Memory:
		return "Memory"
	case NetworkIn:
		return "Network In"
	case NetworkOut:
		return "Network Out"
	default:
		return ""
	}
}

// RetrieveEnvironmentMetrics retrieves metrics data for all services in
// the associated environment.
func (m *SMetrics) RetrieveEnvironmentMetrics(mins int) (*[]models.Metrics, error) {
	headers := httpclient.GetHeaders(m.Settings.SessionToken, m.Settings.Version, m.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/metrics?time=%dm", m.Settings.PaasHost, m.Settings.PaasHostVersion, m.Settings.EnvironmentID, mins), headers)
	if err != nil {
		return nil, err
	}
	var metrics []models.Metrics
	err = httpclient.ConvertResp(resp, statusCode, &metrics)
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}

// RetrieveServiceMetrics retrieves metrics data for the given service.
func (m *SMetrics) RetrieveServiceMetrics(mins int, svcID string) (*models.Metrics, error) {
	headers := httpclient.GetHeaders(m.Settings.SessionToken, m.Settings.Version, m.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/metrics?time=%dm", m.Settings.PaasHost, m.Settings.PaasHostVersion, m.Settings.EnvironmentID, svcID, mins), headers)
	if err != nil {
		return nil, err
	}
	var metrics models.Metrics
	err = httpclient.ConvertResp(resp, statusCode, &metrics)
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}
