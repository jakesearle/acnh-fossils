package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	// "log"
	"math/rand"
	"sort"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

const totalFossils = 73
const nDailyFossils = 4

func main() {
	startTime := time.Now()
	test()
	elapsedTime := time.Since(startTime)
	fmt.Printf("Function took %s\n", elapsedTime)
}

func test() {
	results := runNSims(100_000)
	printStats(results)
	// printHistogram(results)
	plotHistogram(results)
}

func plotHistogram(intData []int) {
	// stdNorm returns the probability of drawing a
	// value from a standard normal distribution.
	// stdNorm := func(x float64) float64 {
	// 	const sigma = 1.0
	// 	const mu = 0.0
	// 	const root2π = 2.50662827459517818309
	// 	return 1.0 / (sigma * root2π) * math.Exp(-((x-mu)*(x-mu))/(2*sigma*sigma))
	// }

	n := len(intData)
	vals := make(plotter.Values, n)
	for i := 0; i < n; i++ {
		// fmt.Printf("%v vs %v", intData[i], float64(intData[i]))
		vals[i] = float64(intData[i])
	}

	p := plot.New()
	p.Title.Text = "Histogram"
	h, err := plotter.NewHist(vals, 16)
	if err != nil {
		log.Panic(err)
	}
	h.Normalize(1)
	p.Add(h)

	// The normal distribution function
	// norm := plotter.NewFunction(stdNorm)
	// norm.Color = color.RGBA{R: 255, A: 255}
	// norm.Width = vg.Points(2)
	// p.Add(norm)

	err = p.Save(500, 500, "./histogram.png")
	if err != nil {
		log.Panic(err)
	}
}

func printHistogram(data []int) {
	sort.Ints(data)
	maxValue := 0
	for _, value := range data {
		if value > maxValue {
			maxValue = value
		}
	}

	var histogram strings.Builder
	for i := maxValue; i > 0; i-- {
		for _, value := range data {
			if value >= i {
				histogram.WriteString("*")
			} else {
				histogram.WriteString(" ")
			}
		}
		histogram.WriteString("\n")
	}

	fmt.Println(histogram.String())
}

func printStats(intData []int) {
	data := intToFloat(intData)
	sort.Float64s(data)
	// Calculate the five-number summary
	min := data[0]
	q1 := stat.Quantile(0.25, stat.Empirical, data, nil)
	median := stat.Quantile(0.5, stat.Empirical, data, nil)
	q3 := stat.Quantile(0.75, stat.Empirical, data, nil)
	max := data[len(data)-1]
	mean := stat.Mean(data, nil)
	mode, count := stat.Mode(data, nil)
	freq := count / float64(len(intData))

	fmt.Printf("Minimum: %v\n", min)
	fmt.Printf("1st Quartile (Q1): %v\n", q1)
	fmt.Printf("Median: %v\n", median)
	fmt.Printf("3rd Quartile (Q3): %v\n", q3)
	fmt.Printf("Maximum: %v\n", max)
	fmt.Printf("Mean: %.2f\n", mean)
	fmt.Printf("Most common: %v, which happened %.2f%% of the time\n", mode, freq)
}

func intToFloat(input []int) []float64 {
	result := make([]float64, len(input))

	for i, v := range input {
		result[i] = float64(v)
	}

	return result
}

func runNSims(nSims int) (daysToSolve []int) {
	bar := pb.New(nSims)
	bar.Start()
	for i := 0; i < nSims; i++ {
		days, err := runSim()
		if err == nil {
			daysToSolve = append(daysToSolve, days)
		}
		bar.Increment()
	}
	bar.Finish()
	return
}

func runSim() (nDays int, e error) {
	fossilCounter := make(map[int]int)
	for day := 1; ; day++ {
		for df := 0; df < nDailyFossils; df++ {
			newFossil := rand.Intn(totalFossils)
			if val, ok := fossilCounter[newFossil]; ok {
				fossilCounter[newFossil] = val + 1
			} else {
				fossilCounter[newFossil] = 1
			}
		}
		if len(fossilCounter) == totalFossils {
			nDays = day
			break
		} else if len(fossilCounter) >= totalFossils {
			return 0, errors.New("Simulation broken")
		}
	}
	return
}
