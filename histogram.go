/*
 * Copyright (C) 2021 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gonum/stat/distuv"
	"github.com/polynetwork/bridge-common/base"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	SRC = "Src->Poly"
	DST = "Poly->Dst"
)

type Point struct {
	Chain     uint64
	Duration  uint64
	Direction string
}

func (p Point) GetChain() Chain {
	return Chain{Chain: p.Chain, Direction: p.Direction}
}

func parseInt(v string) uint64 {
	i, _ := strconv.Atoi(v)
	return uint64(i)
}

type Chain struct {
	Chain     uint64
	Direction string
}

func (c Chain) Name() string {
	switch c.Direction {
	case SRC:
		return fmt.Sprintf("%s->Poly", base.GetChainName(c.Chain))
	case DST:
		return fmt.Sprintf("Poly->%s", base.GetChainName(c.Chain))
	}
	return ""
}

func feedFile(path string) (source map[Chain][]Point, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	points := []Point{}
	lines := strings.Split(string(data), "\n")
	for i := 1; i < len(lines); i++ {
		strs := []string{}
		for _, v := range strings.Split(lines[i], "\t") {
			v = strings.TrimSuffix(strings.TrimSpace(v), "\t")
			if v == "" {
				continue
			}
			strs = append(strs, v)
		}
		if len(strs) == 0 {
			continue
		}
		p1 := Point{
			Chain:     parseInt(strs[5]),
			Duration:  parseInt(strs[0]),
			Direction: SRC,
		}
		p2 := Point{
			Chain:     parseInt(strs[6]),
			Duration:  parseInt(strs[1]),
			Direction: DST,
		}
		if p1.Duration < 4000 {
			points = append(points, p1)
		}
		if p2.Duration < 4000 {
			points = append(points, p2)
		}
	}
	all := map[Chain][]Point{}
	for _, p := range points {
		all[p.GetChain()] = append(all[p.GetChain()], p)
	}
	return all, nil
}

func Draw(ps []Point, path string) (err error) {
	if len(ps) == 0 {
		err = fmt.Errorf("Empty points")
		return
	}
	chain := ps[0].GetChain()
	vs := make(plotter.Values, len(ps))
	for i, p := range ps {
		vs[i] = float64(p.Duration)
	}
	chart, err := plot.New()
	if err != nil {
		return
	}
	chart.Title.Text = chain.Name() + DESC
	h, err := plotter.NewHist(vs, 50)
	if err != nil {
		return
	}
	h.FillColor = plotutil.Color(2)
	// h.Normalize(1)
	chart.Legend.Color = plotutil.Color(3)
	chart.Add(h)
	norm := plotter.NewFunction(distuv.UnitNormal.CDF)
	norm.Color = color.RGBA{R: 255, G: 255, A: 255}
	norm.Width = vg.Points(2)
	chart.Add(norm)
	chart.X.Label.Text = "Duration(s)"
	chart.X.Tick.Length = vg.Length(5)
	chart.Y.Label.Text = "Count"
	chart.Y.Tick.Length = vg.Length(10)
	chart.X.Tick.Marker = Marker{16}
	chart.Y.Tick.Marker = Marker{10}
	return chart.Save(8*vg.Inch, 8*vg.Inch, path)
}

type Marker struct {
	N int
}

func (m Marker) Ticks(min, max float64) []plot.Tick {
	return Tick(min, max, m.N)
}
func Tick(min, max float64, n int) []plot.Tick {
	ticks := []plot.Tick{}
	gap := (max - min) / float64(n)
	for i := 0; i < n; i++ {
		v := gap*float64(i) + min
		ticks = append(ticks, plot.Tick{
			Value: v,
			Label: fmt.Sprintf("%v", int(v)),
		})
	}
	return ticks
}
