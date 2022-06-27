package testdata_test

import (
	"bufio"
	"encoding/xml"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stillmatic/pummel/pkg/model"
	"github.com/stretchr/testify/assert"
)

var MAX_GOROUTINES = 4

type AuditInput struct {
	Age        int
	Employment string
	Education  string
	Marital    string
	Occupation string
	Income     float64
	Gender     string
	Deductions bool
	Hours      int
	Adjusted   int
}

func ParseAuditInput(input string) (map[string]interface{}, error) {
	ai := make(map[string]interface{}, 9)
	// parse from csv
	splits := strings.Split(input, ",")
	age, err := strconv.Atoi(splits[0])
	if err != nil {
		return nil, err
	}
	ai["Age"] = age
	ai["Employment"] = splits[1]
	ai["Education"] = splits[2]
	ai["Marital"] = splits[3]
	ai["Occupation"] = splits[4]
	income, err := strconv.ParseFloat(splits[5], 64)
	if err != nil {
		return nil, err
	}
	ai["Income"] = income
	ai["Gender"] = splits[6]
	deductions, err := strconv.ParseBool(splits[7])
	if err != nil {
		return nil, err
	}
	hours, err := strconv.Atoi(splits[8])
	if err != nil {
		return nil, err
	}
	adj, err := strconv.Atoi(splits[9])
	if err != nil {
		return nil, err
	}
	ai["Deductions"] = deductions
	ai["Hours"] = hours
	ai["Adjusted"] = adj
	return ai, nil
}

//nolint
func BenchmarkAuditLR(b *testing.B) {
	// unmarshal model from file
	var model model.PMMLRegressionModel
	lrmlIO, _ := ioutil.ReadFile("LogisticRegressionAudit.pmml")
	xml.Unmarshal(lrmlIO, &model)
	// load test data
	f, err := os.Open("audit.csv")
	assert.NoError(b, err)
	defer f.Close()
	inputs := make([]map[string]interface{}, 0)

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	for scanner.Scan() {
		input := scanner.Text()
		ai, err := ParseAuditInput(input)
		if err == nil {
			inputs = append(inputs, ai)
		}
	}
	s := rand.NewSource(42)
	r := rand.New(s)

	for i := 0; i < b.N; i++ {
		// for j := 0; j < 10000; j++ {
		inp := inputs[r.Intn(len(inputs))]
		model.RegressionModel.Evaluate(inp)
		// }
	}
}

// nolint
func BenchmarkAuditLRConcurrently(b *testing.B) {
	// unmarshal model from file
	var model model.PMMLRegressionModel
	lrmlIO, _ := ioutil.ReadFile("LogisticRegressionAudit.pmml")
	xml.Unmarshal(lrmlIO, &model)
	// load test data
	f, err := os.Open("audit.csv")
	assert.NoError(b, err)
	defer f.Close()
	inputs := make([]map[string]interface{}, 0)

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	for scanner.Scan() {
		input := scanner.Text()
		ai, err := ParseAuditInput(input)
		if err == nil {
			inputs = append(inputs, ai)
		}
	}
	guard := make(chan struct{}, MAX_GOROUTINES)
	for i := 0; i < b.N; i++ {
		guard <- struct{}{}
		go func(i int) {
			// for j := 0; j < 10000; j++ {
			inp := inputs[i%len(inputs)]
			model.RegressionModel.Evaluate(inp)
			<-guard
			// }
		}(i)
	}
}

func TestRegressionModel(t *testing.T) {
	// unmarshal model from file
	var model model.PMMLRegressionModel
	lrmlIO, _ := ioutil.ReadFile("LogisticRegressionAudit.pmml")
	xml.Unmarshal(lrmlIO, &model)
	assert.Equal(t, 3, len(model.RegressionModel.LocalTransformations.DerivedFields))
	assert.Equal(t, 2, len(model.RegressionModel.RegressionTables))

	// load test data
	f, err := os.Open("audit.csv")
	assert.NoError(t, err)
	defer f.Close()
	inputs := make([]map[string]interface{}, 0)

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	for scanner.Scan() {
		input := scanner.Text()
		ai, err := ParseAuditInput(input)
		if err == nil {
			inputs = append(inputs, ai)
		}
	}
	s := rand.NewSource(42)
	r := rand.New(s)

	for j := 0; j < 10; j++ {
		inp := inputs[r.Intn(len(inputs))]
		res, err := model.RegressionModel.Evaluate(inp)
		t.Log(res)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	}
}

func TestRFModel(t *testing.T) {
	// unmarshal model from file
	var model model.PMMLMiningModel
	lrmlIO, _ := ioutil.ReadFile("RandomForestAudit.pmml")
	xml.Unmarshal(lrmlIO, &model)
	assert.Equal(t, 6, len(model.MiningModel.LocalTransformations.DerivedFields))

	// load test data
	f, err := os.Open("audit.csv")
	assert.NoError(t, err)
	defer f.Close()
	inputs := make([]map[string]interface{}, 0)

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	for scanner.Scan() {
		input := scanner.Text()
		ai, err := ParseAuditInput(input)
		if err == nil {
			inputs = append(inputs, ai)
		}
	}
	s := rand.NewSource(42)
	r := rand.New(s)

	for j := 0; j < 2; j++ {
		inp := inputs[r.Intn(len(inputs))]
		t.Log(inp)
		res, err := model.MiningModel.Evaluate(inp)
		t.Log(res)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 0, 1)
	}
}

//nolint
func BenchmarkAuditRF(b *testing.B) {
	// unmarshal model from file
	var model model.PMMLMiningModel
	lrmlIO, _ := ioutil.ReadFile("RandomForestAudit.pmml")
	xml.Unmarshal(lrmlIO, &model)
	// load test data
	f, err := os.Open("audit.csv")
	assert.NoError(b, err)
	defer f.Close()
	inputs := make([]map[string]interface{}, 0)

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	for scanner.Scan() {
		input := scanner.Text()
		ai, err := ParseAuditInput(input)
		if err == nil {
			inputs = append(inputs, ai)
		}
	}
	s := rand.NewSource(42)
	r := rand.New(s)

	for i := 0; i < b.N; i++ {
		inp := inputs[r.Intn(len(inputs))]
		model.MiningModel.Evaluate(inp)
	}
}

//nolint
func BenchmarkAuditRFConcurrently(b *testing.B) {
	// unmarshal model from file
	var model model.PMMLMiningModel
	lrmlIO, _ := ioutil.ReadFile("RandomForestAudit.pmml")
	xml.Unmarshal(lrmlIO, &model)
	// load test data
	f, err := os.Open("audit.csv")
	assert.NoError(b, err)
	defer f.Close()
	inputs := make([]map[string]interface{}, 0)

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	for scanner.Scan() {
		input := scanner.Text()
		ai, err := ParseAuditInput(input)
		if err == nil {
			inputs = append(inputs, ai)
		}
	}

	guard := make(chan struct{}, MAX_GOROUTINES)
	for i := 0; i < b.N; i++ {
		// for j := 0; j < 10000; j++ {
		guard <- struct{}{}
		go func(i int) {
			inp := inputs[i%len(inputs)]
			model.MiningModel.Evaluate(inp)
			<-guard
		}(i)
		// }
	}
}
