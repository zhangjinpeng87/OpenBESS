package utils

import (
	"math"
	"testing"
)

func TestKalmanFilter(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name                 string
		initialValue         float64
		initialEstimateError float64
		processNoise         float64
		measurementNoise     float64
		sensorData           []float64
		expectedOutput       []float64
	}{
		{
			name:                 "Test1",
			initialValue:         0.0,
			initialEstimateError: 1.0,
			processNoise:         0.01,
			measurementNoise:     0.1,
			sensorData:           []float64{1.2, 1.5, 1.8, 2.0, 2.5},
			expectedOutput:       []float64{1.2, 1.35, 1.58, 1.83, 2.18},
		},
		// Add more test cases as needed
	}

	// Run tests
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			kf := NewKalmanFilter(testCase.initialValue, testCase.initialEstimateError, testCase.processNoise, testCase.measurementNoise)

			// Apply Kalman filter to smooth sensor data
			for i, measurement := range testCase.sensorData {
				smoothedValue := kf.Update(measurement)
				if math.Abs(smoothedValue-testCase.expectedOutput[i]) > 0.001 {
					t.Errorf("Test case %s failed. Expected: %.2f, Got: %.2f", testCase.name, testCase.expectedOutput[i], smoothedValue)
				}
			}
		})
	}
}
