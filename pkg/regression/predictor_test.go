package regression_test

import (
	"encoding/xml"
	"testing"

	"github.com/stillmatic/pummel/pkg/regression"
	"github.com/stretchr/testify/assert"
)

func TestNumericPredictor(t *testing.T) {
	numPredictorXML := []byte(`
		<NumericPredictor name="age" exponent="1" coefficient="7.1"/>
	`)
	var numPredictor regression.NumericPredictor
	err := xml.Unmarshal(numPredictorXML, &numPredictor)
	assert.NoError(t, err)
	assert.Equal(t, "age", numPredictor.Name)
	assert.Equal(t, float64(1), numPredictor.Exponent)
	assert.Equal(t, float64(7.1), numPredictor.Coefficient)
	inputs := map[string]interface{}{
		"age":    float64(30),
		"decade": int(4),
	}
	out, err := numPredictor.Evaluate(inputs)
	assert.NoError(t, err)
	assert.Equal(t, float64(7.1*30), out)
	// testing missing value
	inputs = map[string]interface{}{
		"decade": int(4),
	}
	out, err = numPredictor.Evaluate(inputs)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), out)
}

func TestCategoricalPredictor(t *testing.T) {
	catPredictorXML := []byte(`
		<CategoricalPredictor name="car_location" value="carpark" coefficient="41.1"/>
	`)
	var catPredictor regression.CategoricalPredictor
	err := xml.Unmarshal(catPredictorXML, &catPredictor)
	assert.NoError(t, err)
	assert.Equal(t, "car_location", catPredictor.Name)
	assert.Equal(t, "carpark", catPredictor.Value)
	assert.Equal(t, float64(41.1), catPredictor.Coefficient)
	inputs := map[string]interface{}{
		"car_location": "carpark",
	}
	out, err := catPredictor.Evaluate(inputs)
	assert.NoError(t, err)
	assert.Equal(t, float64(41.1), out)
	// testing incorrect value
	inputs = map[string]interface{}{
		"car_location": "street",
	}
	out, err = catPredictor.Evaluate(inputs)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), out)
	// testing missing value
	inputs = map[string]interface{}{
		"sesame": "street",
	}
	out, err = catPredictor.Evaluate(inputs)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), out)

}
