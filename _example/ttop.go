// Copyright 2016 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gizak/termui"
	"github.com/gizak/termui/extra"
)

const statFilePath = "/proc/stat"
const meminfoFilePath = "/proc/meminfo"

type CpuStat struct {
	user   float32
	nice   float32
	system float32
	idle   float32
}

type CpusStats struct {
	stat map[string]CpuStat
	proc map[string]CpuStat
}

func NewCpusStats(s map[string]CpuStat) *CpusStats {
	return &CpusStats{stat: s, proc: make(map[string]CpuStat)}
}

func (cs *CpusStats) String() (ret string) {
	for key, _ := range cs.proc {
		ret += fmt.Sprintf("%s: %.2f %.2f %.2f %.2f\n", key, cs.proc[key].user, cs.proc[key].nice, cs.proc[key].system, cs.proc[key].idle)
	}
	return
}

func subCpuStat(m CpuStat, s CpuStat) CpuStat {
	return CpuStat{user: m.user - s.user,
		nice:   m.nice - s.nice,
		system: m.system - s.system,
		idle:   m.idle - s.idle}
}

func procCpuStat(c CpuStat) CpuStat {
	sum := c.user + c.nice + c.system + c.idle
	return CpuStat{user: c.user / sum * 100,
		nice:   c.nice / sum * 100,
		system: c.system / sum * 100,
		idle:   c.idle / sum * 100}
}

func (cs *CpusStats) tick(ns map[string]CpuStat) {
	for key, _ := range cs.stat {
		proc := subCpuStat(ns[key], cs.stat[key])
		cs.proc[key] = procCpuStat(proc)
		cs.stat[key] = ns[key]
	}
}

type errIntParser struct {
	err error
}

func (eip *errIntParser) parse(s string) (ret int64) {
	if eip.err != nil {
		return 0
	}
	ret, eip.err = strconv.ParseInt(s, 10, 0)
	return
}

type LineProcessor interface {
	process(string) error
	finalize() interface{}
}

type CpuLineProcessor struct {
	m map[string]CpuStat
}

func (clp *CpuLineProcessor) process(line string) (err error) {
	r := regexp.MustCompile("^cpu([0-9]*)")

	if r.MatchString(line) {
		tab := strings.Fields(line)
		if len(tab) < 5 {
			err = errors.New("cpu info line has not enough fields")
			return
		}
		parser := errIntParser{}
		cs := CpuStat{user: float32(parser.parse(tab[1])),
			nice:   float32(parser.parse(tab[2])),
			system: float32(parser.parse(tab[3])),
			idle:   float32(parser.parse(tab[4]))}
		clp.m[tab[0]] = cs
		err = parser.err
		if err != nil {
			return
		}
	}
	return
}

func (clp *CpuLineProcessor) finalize() interface{} {
	return clp.m
}

type MemStat struct {
	total int64
	free  int64
}

func (ms MemStat) String() (ret string) {
	ret = fmt.Sprintf("TotalMem: %d, FreeMem: %d\n", ms.total, ms.free)
	return
}

func (ms *MemStat) process(line string) (err error) {
	rtotal := regexp.MustCompile("^MemTotal:")
	rfree := regexp.MustCompile("^MemFree:")
	var aux int64
	if rtotal.MatchString(line) || rfree.MatchString(line) {
		tab := strings.Fields(line)
		if len(tab) < 3 {
			err = errors.New("mem info line has not enough fields")
			return
		}
		aux, err = strconv.ParseInt(tab[1], 10, 0)
	}
	if err != nil {
		return
	}

	if rtotal.MatchString(line) {
		ms.total = aux
	}
	if rfree.MatchString(line) {
		ms.free = aux
	}
	return
}

func (ms *MemStat) finalize() interface{} {
	return *ms
}

func processFileLines(filePath string, lp LineProcessor) (ret interface{}, err error) {
	var statFile *os.File
	statFile, err = os.Open(filePath)
	if err != nil {
		fmt.Printf("open: %v\n", err)
	}
	defer statFile.Close()

	statFileReader := bufio.NewReader(statFile)

	for {
		var line string
		line, err = statFileReader.ReadString('\n')
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			fmt.Printf("open: %v\n", err)
			break
		}
		line = strings.TrimSpace(line)

		err = lp.process(line)
	}

	ret = lp.finalize()
	return
}

func getCpusStatsMap() (m map[string]CpuStat, err error) {
	var aux interface{}
	aux, err = processFileLines(statFilePath, &CpuLineProcessor{m: make(map[string]CpuStat)})
	return aux.(map[string]CpuStat), err
}

func getMemStats() (ms MemStat, err error) {
	var aux interface{}
	aux, err = processFileLines(meminfoFilePath, &MemStat{})
	return aux.(MemStat), err
}

type CpuTabElems struct {
	GMap   map[string]*termui.Gauge
	LChart *termui.LineChart
}

func NewCpuTabElems(width int) *CpuTabElems {
	lc := termui.NewLineChart()
	lc.Width = width
	lc.Height = 12
	lc.X = 0
	lc.Mode = "dot"
	lc.Border.Label = "CPU"
	return &CpuTabElems{GMap: make(map[string]*termui.Gauge),
		LChart: lc}
}

func (cte *CpuTabElems) AddGauge(key string, Y int, width int) *termui.Gauge {
	cte.GMap[key] = termui.NewGauge()
	cte.GMap[key].Width = width
	cte.GMap[key].Height = 3
	cte.GMap[key].Y = Y
	cte.GMap[key].Border.Label = key
	cte.GMap[key].Percent = 0 //int(val.user + val.nice + val.system)
	return cte.GMap[key]
}

func (cte *CpuTabElems) Update(cs CpusStats) {
	for key, val := range cs.proc {
		p := int(val.user + val.nice + val.system)
		cte.GMap[key].Percent = p
		if key == "cpu" {
			cte.LChart.Data = append(cte.LChart.Data, 0)
			copy(cte.LChart.Data[1:], cte.LChart.Data[0:])
			cte.LChart.Data[0] = float64(p)
		}
	}
}

type MemTabElems struct {
	Gauge  *termui.Gauge
	SLines *termui.Sparklines
}

func NewMemTabElems(width int) *MemTabElems {
	g := termui.NewGauge()
	g.Width = width
	g.Height = 3
	g.Y = 0

	sline := termui.NewSparkline()
	sline.Title = "MEM"
	sline.Height = 8

	sls := termui.NewSparklines(sline)
	sls.Width = width
	sls.Height = 12
	sls.Y = 3
	return &MemTabElems{Gauge: g, SLines: sls}
}

func (mte *MemTabElems) Update(ms MemStat) {
	used := int((ms.total - ms.free) * 100 / ms.total)
	mte.Gauge.Percent = used
	mte.SLines.Lines[0].Data = append(mte.SLines.Lines[0].Data, 0)
	copy(mte.SLines.Lines[0].Data[1:], mte.SLines.Lines[0].Data[0:])
	mte.SLines.Lines[0].Data[0] = used
	if len(mte.SLines.Lines[0].Data) > mte.SLines.Width-2 {
		mte.SLines.Lines[0].Data = mte.SLines.Lines[0].Data[0 : mte.SLines.Width-2]
	}
}

func main() {
	if runtime.GOOS != "linux" {
		panic("Currently works only on Linux")
	}
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	termWidth := 70

	termui.UseTheme("helloworld")

	header := termui.NewPar("Press q to quit, Press j or k to switch tabs")
	header.Height = 1
	header.Width = 50
	header.HasBorder = false
	header.TextBgColor = termui.ColorBlue

	tabCpu := extra.NewTab("CPU")
	tabMem := extra.NewTab("MEM")

	tabpane := extra.NewTabpane()
	tabpane.Y = 1
	tabpane.Width = 30
	tabpane.HasBorder = false

	cs, errcs := getCpusStatsMap()
	cpusStats := NewCpusStats(cs)

	if errcs != nil {
		panic("error")
	}

	cpuTabElems := NewCpuTabElems(termWidth)

	Y := 0
	cpuKeys := make([]string, 0, len(cs))
	for key := range cs {
		cpuKeys = append(cpuKeys, key)
	}
	sort.Strings(cpuKeys)
	for _, key := range cpuKeys {
		g := cpuTabElems.AddGauge(key, Y, termWidth)
		Y += 3
		tabCpu.AddBlocks(g)
	}
	cpuTabElems.LChart.Y = Y
	tabCpu.AddBlocks(cpuTabElems.LChart)

	memTabElems := NewMemTabElems(termWidth)
	ms, errm := getMemStats()
	if errm != nil {
		panic(errm)
	}
	memTabElems.Update(ms)
	tabMem.AddBlocks(memTabElems.Gauge)
	tabMem.AddBlocks(memTabElems.SLines)

	tabpane.SetTabs(*tabCpu, *tabMem)

	termui.Render(header, tabpane)

	evt := termui.EventCh()
	for {
		select {
		case e := <-evt:
			if e.Type == termui.EventKey {
				switch e.Ch {
				case 'q':
					return
				case 'j':
					tabpane.SetActiveLeft()
					termui.Render(header, tabpane)
				case 'k':
					tabpane.SetActiveRight()
					termui.Render(header, tabpane)
				}
			}
		case <-time.After(time.Second):
			cs, errcs := getCpusStatsMap()
			if errcs != nil {
				panic(errcs)
			}
			cpusStats.tick(cs)
			cpuTabElems.Update(*cpusStats)

			ms, errm := getMemStats()
			if errm != nil {
				panic(errm)
			}
			memTabElems.Update(ms)

			termui.Render(header, tabpane)
		}
	}
}
