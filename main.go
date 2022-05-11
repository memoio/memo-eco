package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/BurntSushi/toml"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"github.com/memoio/memo-eco/model"
)

func main() {
	// load config
	var cp string
	flag.StringVar(&cp, "config", "", "config file path")
	flag.Parse()

	cfg := new(model.Config)

	outputDir := "output"
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
			return
		}
	} else {
		fmt.Println("load config from: ", cp)
		_, err := toml.DecodeFile(cp, cfg)
		if err != nil {
			return
		}
	}

	// 有关代币相关的图
	p := plot.New()

	p.Title.Text = "Memo Token"
	p.X.Label.Text = "Time(Day)"
	p.Y.Label.Text = "Token(Memo)"

	// 有关空间占用的图
	pSize := plot.New()

	pSize.Title.Text = "Memo Size"
	pSize.X.Label.Text = "Time(Day)"
	pSize.Y.Label.Text = "Size(GB)"

	points := model.Simulate(cfg)

	err := plotutil.AddLinePoints(p,
		"Supply", points[model.SUPPLY_INDEX],
		"Liquid", points[model.LIQUID_INDEX],
		"Reward", points[model.REWARD_INDEX],
		"Pledge", points[model.PLEDGE_INDEX],
		"Paid", points[model.PAID_INDEX],
	)

	if err != nil {
		panic(err)
	}

	err = plotutil.AddLinePoints(pSize,
		"Size", points[model.SIZE_INDEX],
		"ASize", points[model.ASIZE_INDEX],
	)

	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(96*vg.Inch, 48*vg.Inch, path.Join(outputDir, "token.png")); err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := pSize.Save(96*vg.Inch, 48*vg.Inch, path.Join(outputDir, "size.png")); err != nil {
		panic(err)
	}
}
