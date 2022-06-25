package regression_test

import (
	"encoding/xml"
	"testing"

	"github.com/stillmatic/pummel/pkg/regression"
	"github.com/stretchr/testify/assert"
)

var regressionTableXML = []byte(`
  <RegressionTable intercept="132.37">
	<NumericPredictor name="age" exponent="1" coefficient="7.1"/>
	<NumericPredictor name="salary" exponent="1" coefficient="0.01"/>
	<CategoricalPredictor name="car_location" value="carpark" coefficient="41.1"/>
	<CategoricalPredictor name="car_location" value="street" coefficient="325.03"/>
  </RegressionTable>
`)

func TestRegressionTable(t *testing.T) {
	var rt regression.RegressionTable
	err := xml.Unmarshal(regressionTableXML, &rt)
	assert.NoError(t, err)
	assert.Equal(t, float64(132.37), rt.Intercept)
	assert.Equal(t, 4, len(rt.Predictors))
	inputs := map[string]interface{}{
		"age":              float64(30),
		"salary":           float64(1000),
		"car_location":     "carpark",
		"number_of_claims": float64(0),
	}
	out, err := rt.Evaluate(inputs)
	assert.NoError(t, err)
	assert.Equal(t, float64(396.47), out)
}

//nolint
func BenchmarkRegressionTable(b *testing.B) {
	var rt regression.RegressionTable
	xml.Unmarshal(regressionTableXML, &rt)
	inputs := map[string]interface{}{
		"age":              float64(30),
		"salary":           float64(1000),
		"car_location":     "carpark",
		"number_of_claims": float64(0),
	}
	for i := 0; i < b.N; i++ {
		rt.Evaluate(inputs)
	}
}
