package main

import (
	"bufio"
	"fmt"
	"os"
	"log"
	"strconv"
	"strings"
        "math/rand"
	"math"
	"github.com/unixpickle/wav"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("usage: zre_proj1_linux cb_lpc.txt cb_gain.txt in.cod out.wav")
		os.Exit(1)
	}

	codedFile := CodedFile{}
	codedFile.read(os.Args[3])

	cbLpc := Matrix{}
	cbLpc.load(os.Args[1])

	cbGain := Matrix{}
	cbGain.load(os.Args[2])

	codedFile.decode(cbLpc, cbGain)
	codedFile.synthetize()

	data := codedFile.result

	sampleRate := 8000
	sound := wav.NewPCM16Sound(1, sampleRate)
	for _, val := range data{
		value := wav.Sample(val)
		sound.SetSamples(append(sound.Samples(), value))
	}
	wav.WriteFile(sound, os.Args[4])
	fmt.Println("Done decoding.")
}

type Matrix struct {
	data [][]float64
}

func (m *Matrix) load(fileName string)  {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var cols, rows int
	var data []float64
	for scanner.Scan() {
		lineParts := strings.Fields(scanner.Text())
		cols = 0
		for _, word := range lineParts {
			num, err := strconv.ParseFloat(word, 64)
			data = append(data, num)

			if err != nil {
				log.Fatal(err)
			}
			cols++
		}
		rows++
	}
	file.Close()
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}

	for i := 0; i < cols; i++  {
		vector := make([]float64, rows)
		for j := 0; j < rows; j++ {
			vector[j] = data[j*cols+i]
		}
		m.data = append(m.data, vector)
	}
}

type CodedFile struct {
	aCoded, gCoded, ls []int
	signal []float64
	gDecoded, aDecoded [][]float64
	result []float64
}

func (c *CodedFile) decode(codebookLpc, codebookGain Matrix) {
	for _, value := range c.aCoded {
		c.aDecoded = append(c.aDecoded, codebookLpc.data[value -1])
	}

	for _, value := range c.gCoded {
		c.gDecoded = append(c.gDecoded, codebookGain.data[value -1])
	}
}

func (c *CodedFile) synthetize() {
	var result []float64
	frameLength := 160
	frameCount := len(c.gCoded)
	initAmem := make([]float64, 11)
	nextvoiced := 0;

	for i := 0; i < frameCount; i++ {
		aDecoded := c.aDecoded[i]
		aDecoded = append([]float64{0.0}, aDecoded...)
		gainDecoded := c.gDecoded[i][0]
		lag := c.ls[i]
		var excit []float64
		if lag == 0 {
			excit = randomFrames(frameLength)
		} else {
			var where []int
			for nextvoiced < frameLength {
				where = append(where, nextvoiced)
				nextvoiced += lag
			}
			nextvoiced = where[len(where)-1] + lag - frameLength

			excit = make([]float64, frameLength)
			for _, value := range where {
				excit[value] = 1
			}
		}
		total := 0.0
		for _, value := range excit {
			total += math.Pow(value, 2)
		}
		total /= float64(frameLength)
		total = math.Sqrt(total)

		for key, val := range excit {
			excit[key] = val / total
		}

		res, amem := filter(excit, gainDecoded, initAmem, aDecoded)
		initAmem = amem
		result = append(result, res...)
	}

	c.result = result
}

func filter(excit []float64, gain float64, amem []float64, asym []float64) ([] float64, []float64) {
	coefficinetCount := 10
	var result []float64

	for _, value := range excit {
		sum := value * gain

		for i:= coefficinetCount; i >= 1; i--  {
			sum -= amem[i] * asym[i]
			amem[i] = amem[i-1]
		}
		amem[1] = sum

		result = append(result, sum)
	}
	return result, amem
}

func randomFrames(size int) []float64  {
	frames := make([]float64, size)
	for i := range frames {
		frames[i] = rand.Float64() * (1.0 - 0.0) + 0.0
	}
	return frames;
}

func (c *CodedFile) read(fileName string){
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineParts := strings.Fields(scanner.Text())
		num, err := strconv.Atoi(lineParts[0])
		c.aCoded = append(c.aCoded, num)
		num, err = strconv.Atoi(lineParts[1])
		c.gCoded = append(c.gCoded, num)
		num, err = strconv.Atoi(lineParts[2])
		c.ls = append(c.ls, num)
		if err != nil {
			log.Fatal(err)
		}
	}

	file.Close()
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
}