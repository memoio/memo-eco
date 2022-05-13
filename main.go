package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/mitchellh/go-homedir"

	"github.com/memoio/memo-eco/model"
)

func main() {
	// load config
	var cp string
	flag.StringVar(&cp, "config", "", "config file path")
	flag.Parse()

	cfg := new(model.Config)

	outputDir := "~/.simu"
	outputDir, err := homedir.Expand(outputDir)
	if err != nil {
		return
	}
	err = os.MkdirAll(outputDir, 0700)
	if err != nil {
		return
	}
	if cp == "" {
		cp = path.Join(outputDir, "config.toml")
		fmt.Println("create and save config to:", cp)
		cfg = model.DefaultConfig()
		buf := new(bytes.Buffer)
		err := toml.NewEncoder(buf).Encode(cfg)
		if err != nil {
			return
		}

		err = ioutil.WriteFile(cp, buf.Bytes(), 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("load config from: ", cp)
		_, err := toml.DecodeFile(cp, cfg)
		if err != nil {
			return
		}
	}

	model.Simulate(cfg)

	plotToken(outputDir)
	plotSize(outputDir)

	end := "0.0.0.0:18081"
	fmt.Println("visit: ", end)
	fs := http.FileServer(http.Dir(outputDir))
	http.ListenAndServe(end, logRequest(fs))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func plotToken(outputDir string) {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title: "Memo Token and Reward",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Memo",
			SplitLine: &opts.SplitLine{
				Show: false,
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Day",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: true,
		}),
	)

	line.SetXAxis(model.PlotX).
		AddSeries("Supply", model.PlotData[model.SUPPLY_INDEX]).
		AddSeries("Liquid", model.PlotData[model.LIQUID_INDEX]).
		AddSeries("Reward", model.PlotData[model.REWARD_INDEX]).
		AddSeries("Pledge", model.PlotData[model.PLEDGE_INDEX]).
		AddSeries("Paid", model.PlotData[model.PAID_INDEX]).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				Smooth: true,
			}),
		)

	f, err := os.Create(filepath.Join(outputDir, "token.html"))
	if err != nil {
		panic(err)
	}

	line.Render(io.MultiWriter(f))
}

func plotSize(outputDir string) {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title: "Memo Size",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "TB",
			SplitLine: &opts.SplitLine{
				Show: false,
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Day",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: true,
		}),
	)

	line.SetXAxis(model.PlotX).
		AddSeries("Size", model.PlotData[model.SIZE_INDEX]).
		AddSeries("Total Size", model.PlotData[model.ASIZE_INDEX]).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				Smooth: true,
			}),
		)

	f, err := os.Create(filepath.Join(outputDir, "size.html"))
	if err != nil {
		panic(err)
	}

	line.Render(io.MultiWriter(f))
}
