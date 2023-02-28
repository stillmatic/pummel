package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
	"github.com/stillmatic/pummel/pkg/model"
	"github.com/stillmatic/pummel/pkg/utils"
)

type TestModelCmd struct {
	ModelFile          string `arg:"" help:"Model file to test"`
	SampleInputFile    string `arg:"" help:"Sample input csv to test"`
	ExpectedOutputFile string `arg:"" help:"Expected output file"`
}

func loadOutputFile(filename string) ([]float64, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	outputs := make([]float64, 0)
	for scanner.Scan() {
		txt := scanner.Text()
		if strings.HasSuffix(txt, "]") {
			txt = txt[:len(txt)-1]
		}
		if strings.HasPrefix(txt, "[") {
			txt = txt[1:]
		}
		val, err := strconv.ParseFloat(txt, 64)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, val)
	}
	return outputs, nil
}

const float64EqualityThreshold = 1e-8

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func (t *TestModelCmd) Run(ctx *kong.Context) error {
	var model model.PMMLMiningModel
	lrmlIO, _ := ioutil.ReadFile(t.ModelFile)
	err := xml.Unmarshal(lrmlIO, &model)
	if err != nil {
		return errors.Wrap(err, "couldnt unmarshal model")
	}
	data, err := utils.LoadData(t.SampleInputFile)
	if err != nil {
		return errors.Wrap(err, "couldnt load sample input")
	}
	expectedOutputs, err := loadOutputFile(t.ExpectedOutputFile)
	if err != nil {
		return errors.Wrap(err, "couldnt load expected output")
	}
	for i, row := range data {
		res, err := model.MiningModel.Evaluate(row)
		if err != nil {
			return errors.Wrap(err, "couldnt evaluate model")
		}

		val, ok := res["probability(1)"]
		if !ok {
			val, ok = res["probability(1.0)"]
			if !ok {
				return errors.Errorf("couldnt find probability in result")
			}
		}

		floatVal, ok := val.(float64)
		if !ok {
			return errors.Errorf("expected float64, got %T", val)
		}
		if !almostEqual(floatVal, expectedOutputs[i]) {
			fmt.Println(fmt.Errorf("failed on %v, expected %f, got %f", i, expectedOutputs[i], floatVal))
		}
	}

	return nil

}
