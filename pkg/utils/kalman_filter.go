// File: kalman_filter.go
// Package: soc
// For statistics and control theory, Kalman filtering, also known as linear quadratic estimation (LQE),
// is an algorithm that uses a series of measurements observed over time, including statistical noise and
// other inaccuracies, and produces estimates of unknown variables that tend to be more accurate than 
// those based on a single measurement alone, by estimating a joint probability distribution over the 
// variables for each timeframe. The filter is named after Rudolf E. Kálmán, who was one of the primary 
// developers of its theory.
// Ref to https://en.wikipedia.org/wiki/Kalman_filter

package utils

import (
	"fmt"
	"math"
)

// KalmanFilter represents the Kalman filter state.
type KalmanFilter struct {
	xHat      float64 // State estimate
	p         float64 // Estimate error covariance
	q         float64 // Process noise covariance
	r         float64 // Measurement noise covariance
	k         float64 // Kalman gain
}

// NewKalmanFilter initializes a new KalmanFilter with provided initial values.
func NewKalmanFilter(initialValue, initialEstimateError, processNoise, measurementNoise float64) *KalmanFilter {
	return &KalmanFilter{
		xHat: initialValue,
		p:    initialEstimateError,
		q:    processNoise,
		r:    measurementNoise,
	}
}

// Update performs a single update step of the Kalman filter.
func (kf *KalmanFilter) Update(measurement float64) float64 {
	// Prediction step
	// State estimate prediction
	xHatMinus := kf.xHat
	// Estimate error covariance prediction
	pMinus := kf.p + kf.q

	// Update step
	// Kalman gain calculation
	kf.k = pMinus / (pMinus + kf.r)
	// State estimate update
	kf.xHat = xHatMinus + kf.k*(measurement-xHatMinus)
	// Estimate error covariance update
	kf.p = (1 - kf.k) * pMinus

	return kf.xHat
}
