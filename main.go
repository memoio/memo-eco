package main

import (
	"math/rand"
	"path"

	"github.com/memoio/memo-eco/model"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {
	rand.Seed(int64(0))

	// 有关代币相关的图
	p := plot.New()

	p.Title.Text = "Memo Token"
	p.X.Label.Text = "Time(Day)"
	p.Y.Label.Text = "Token(Memo)"

	// 有关代币相关的图
	pPay := plot.New()

	pPay.Title.Text = "Memo Token"
	pPay.X.Label.Text = "Time(Day)"
	pPay.Y.Label.Text = "Token(Memo)"

	// 有关空间占用的图
	pSize := plot.New()

	pSize.Title.Text = "Memo Size"
	pSize.X.Label.Text = "Time(Day)"
	pSize.Y.Label.Text = "Size(GB)"

	config := model.DefaultEconomicsConfig()
	points := model.EcoModelSimulate(config)

	err := plotutil.AddLinePoints(p,
		"Supply", points[model.SUPPLY_INDEX],
		"Liquid", points[model.LIQUID_INDEX],
		"Reward", points[model.REWARD_INDEX],
		"Pledge", points[model.PLEDGE_INDEX],
	)

	if err != nil {
		panic(err)
	}

	err = plotutil.AddLinePoints(pPay,
		"Paid", points[model.PAID_INDEX],
	)

	if err != nil {
		panic(err)
	}

	err = plotutil.AddLinePoints(pSize,
		"Size", points[model.SIZE_INDEX],
	)

	if err != nil {
		panic(err)
	}

	outputDir := "output"

	// Save the plot to a PNG file.
	if err := p.Save(96*vg.Inch, 48*vg.Inch, path.Join(outputDir, "token.png")); err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := pSize.Save(96*vg.Inch, 48*vg.Inch, path.Join(outputDir, "size.png")); err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := pPay.Save(96*vg.Inch, 48*vg.Inch, path.Join(outputDir, "pay.png")); err != nil {
		panic(err)
	}
}
