package tests

import (
	"bytes"
	"github.com/wcharczuk/go-chart" //exposes "chart"
	"testing"
)

func Test_BasicPngChart(t *testing.T) {

	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0},
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)

	if err != nil {
		panic(err)
	}
}
