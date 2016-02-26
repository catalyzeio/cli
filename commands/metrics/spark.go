package metrics

import (
	"github.com/catalyzeio/cli/models"
	ui "github.com/gizak/termui"
)

const (
	titleColor      = ui.ColorWhite
	cpuColor        = ui.ColorBlue
	memoryColor     = ui.ColorGreen
	networkInColor  = ui.ColorRed
	networkOutColor = ui.ColorWhite
)

// SparkTransformer is a concrete implementation of Transformer transforming
// data into spark lines.
type SparkTransformer struct {
	SparkLines map[string]*ui.Sparklines
}

// TransformGroupCPU transforms an entire environment's cpu data into spark
// lines. This outputs TransformSingleCPU for every service in the environment.
func (spark *SparkTransformer) TransformGroupCPU(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		if _, ok := blacklist[metric.ServiceLabel]; !ok {
			spark.TransformSingleCPU(&metric)
		}
	}
}

// TransformGroupMemory transforms an entire environment's memory data into
// spark lines. This outputs TransformSingleMemory for every service in the
// environment.
func (spark *SparkTransformer) TransformGroupMemory(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		if _, ok := blacklist[metric.ServiceLabel]; !ok {
			spark.TransformSingleMemory(&metric)
		}
	}
}

// TransformGroupNetworkIn transforms an entire environment's received network
// data into spark lines. This outputs TransformSingleNetworkIn for every
// service in the environment.
func (spark *SparkTransformer) TransformGroupNetworkIn(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		if _, ok := blacklist[metric.ServiceLabel]; !ok {
			spark.TransformSingleNetworkIn(&metric)
		}
	}
}

// TransformGroupNetworkOut transforms an entire environment's transmitted
// network data into spark lines. This outputs TransformSingleNetworkOut for
// every service in the environment.
func (spark *SparkTransformer) TransformGroupNetworkOut(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		if _, ok := blacklist[metric.ServiceLabel]; !ok {
			spark.TransformSingleNetworkOut(&metric)
		}
	}
}

// TransformSingleCPU transforms a single service's cpu data into spark lines.
func (spark *SparkTransformer) TransformSingleCPU(metric *models.Metrics) {
	var cpuMin []int
	var cpuMax []int
	var cpuAvg []int
	var cpuTotal []int
	if metric.Data != nil && metric.Data.CPUUsage != nil {
		for _, data := range *metric.Data.CPUUsage {
			cpuMin = append(cpuMin, int(data.Min/1000.0))
			cpuMax = append(cpuMax, int(data.Max/1000.0))
			cpuAvg = append(cpuAvg, int(data.AVG/1000.0))
			cpuTotal = append(cpuTotal, int(data.Total/1000.0))
		}
	}
	var sparkLines = spark.SparkLines[metric.ServiceLabel]
	if sparkLines == nil {
		sparkLines = addSparkLine(metric.ServiceLabel, []string{"CPU Min", "CPU Max", "CPU AVG", "CPU Total"}, cpuColor)
		spark.SparkLines[metric.ServiceLabel] = sparkLines
	}
	for i := range sparkLines.Lines {
		if sparkLines.Lines[i].Title == "CPU Min" {
			sparkLines.Lines[i].Data = cpuMin
		} else if sparkLines.Lines[i].Title == "CPU Max" {
			sparkLines.Lines[i].Data = cpuMax
		} else if sparkLines.Lines[i].Title == "CPU AVG" {
			sparkLines.Lines[i].Data = cpuAvg
		} else if sparkLines.Lines[i].Title == "CPU Total" {
			sparkLines.Lines[i].Data = cpuTotal
		}
	}
	ui.Render(ui.Body)
}

// TransformSingleMemory transforms a single service's memory data into spark
// lines.
func (spark *SparkTransformer) TransformSingleMemory(metric *models.Metrics) {
	var memMin []int
	var memMax []int
	var memAvg []int
	var memTotal []int
	if metric.Data != nil && metric.Data.MemoryUsage != nil {
		for _, data := range *metric.Data.MemoryUsage {
			memMin = append(memMin, int(data.Min/1024.0))
			memMax = append(memMax, int(data.Max/1024.0))
			memAvg = append(memAvg, int(data.AVG/1024.0))
			memTotal = append(memTotal, int(metric.Size.RAM*1024.0))
		}
	}
	var sparkLines = spark.SparkLines[metric.ServiceLabel]
	if sparkLines == nil {
		sparkLines = addSparkLine(metric.ServiceLabel, []string{"Mem Min", "Mem Max", "Mem AVG", "Mem Total"}, memoryColor)
		spark.SparkLines[metric.ServiceLabel] = sparkLines
	}
	for i := range sparkLines.Lines {
		if sparkLines.Lines[i].Title == "Mem Min" {
			sparkLines.Lines[i].Data = memMin
		} else if sparkLines.Lines[i].Title == "Mem Max" {
			sparkLines.Lines[i].Data = memMax
		} else if sparkLines.Lines[i].Title == "Mem AVG" {
			sparkLines.Lines[i].Data = memAvg
		} else if sparkLines.Lines[i].Title == "Mem Total" {
			sparkLines.Lines[i].Data = memTotal
		}
	}
	ui.Render(ui.Body)
}

// TransformSingleNetworkIn transforms a single service's received network data
// into spark lines.
func (spark *SparkTransformer) TransformSingleNetworkIn(metric *models.Metrics) {
	var netinKB []int
	var netinPackets []int
	if metric.Data != nil && metric.Data.NetworkUsage != nil {
		for _, data := range *metric.Data.NetworkUsage {
			netinKB = append(netinKB, int(data.RXKB))
			netinPackets = append(netinPackets, int(data.RXPackets))
		}
	}
	var sparkLines = spark.SparkLines[metric.ServiceLabel]
	if sparkLines == nil {
		sparkLines = addSparkLine(metric.ServiceLabel, []string{"Received KB", "Received Packets"}, networkInColor)
		spark.SparkLines[metric.ServiceLabel] = sparkLines
	}
	for i := range sparkLines.Lines {
		if sparkLines.Lines[i].Title == "Received KB" {
			sparkLines.Lines[i].Data = netinKB
		} else if sparkLines.Lines[i].Title == "Received Packets" {
			sparkLines.Lines[i].Data = netinPackets
		}
	}
	ui.Render(ui.Body)
}

// TransformSingleNetworkOut transforms a single service's transmitted network
// data into spark lines.
func (spark *SparkTransformer) TransformSingleNetworkOut(metric *models.Metrics) {
	var netoutKB []int
	var netoutPackets []int
	if metric.Data != nil && metric.Data.NetworkUsage != nil {
		for _, data := range *metric.Data.NetworkUsage {
			netoutKB = append(netoutKB, int(data.TXKB))
			netoutPackets = append(netoutPackets, int(data.TXPackets))
		}
	}
	var sparkLines = spark.SparkLines[metric.ServiceLabel]
	if sparkLines == nil {
		sparkLines = addSparkLine(metric.ServiceLabel, []string{"Transmitted KB", "Transmitted Packets"}, networkOutColor)
		spark.SparkLines[metric.ServiceLabel] = sparkLines
	}
	for i := range sparkLines.Lines {
		if sparkLines.Lines[i].Title == "Transmitted KB" {
			sparkLines.Lines[i].Data = netoutKB
		} else if sparkLines.Lines[i].Title == "Transmitted Packets" {
			sparkLines.Lines[i].Data = netoutPackets
		}
	}
	ui.Render(ui.Body)
}

func addSparkLine(serviceName string, titles []string, color ui.Attribute) *ui.Sparklines {
	var sparkLines []ui.Sparkline
	for _, title := range titles {
		sparkLine := ui.NewSparkline()
		sparkLine.Height = 1
		sparkLine.Data = []int{}
		sparkLine.Title = title
		sparkLine.TitleColor = titleColor
		sparkLine.LineColor = color
		sparkLines = append(sparkLines, sparkLine)
	}
	sp := ui.NewSparklines(sparkLines...)
	sp.Height = 11
	sp.BorderLabel = serviceName

	ui.Body.AddRows(
		ui.NewRow(ui.NewCol(12, 0, sp)),
	)

	ui.Body.Align()
	ui.Render(sp)
	ui.Render(ui.Body)

	return sp
}

func sparkLinesEventLoop() {
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Render(ui.Body)
	})

	ui.Loop() // blocking call
}
