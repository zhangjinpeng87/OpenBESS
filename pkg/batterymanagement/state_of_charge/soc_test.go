package soc

import (
	"math"
	"testing"
)

func TestDischargeCalculaterSOC(t *testing.T) {
	calculator := &DischargeCalculater{
		DisMaxVoltage:      4.25,
		DisMidHighVoltage:  3.7,
		DisMidLowVoltage:   3.0,
		DisMinVoltage:      2.5,
	}

	tests := []struct {
		name        string
		voltage     float64
		expectedSOC float64
	}{
		{"Voltage above max", 4.5, 100},
		{"Voltage in mid-high range", 3.8, 92.5},
		{"Voltage in mid-low range", 3.2, 70},
		{"Voltage in min range", 2.7, 30},
		{"Voltage below min", 2.0, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			soc := calculator.SOC(test.voltage)
			if math.Abs(soc-test.expectedSOC) > 0.1 {
				t.Errorf("Expected SOC %.2f for voltage %.2f, got %.2f", test.expectedSOC, test.voltage, soc)
			}
		})
	}
}

func TestChargeCalculaterSOC(t *testing.T) {
	calculator := &ChargeCalculater{
		ChMaxVoltage:         4.2,
		ChMidHighVoltage:     3.8,
		ChMidLowVoltage:      3.5,
		ChMinVoltage:         2.5,
		ChMaxChargingCurrent: 1.0,
	}

	tests := []struct {
		name         string
		voltage      float64
		current      float64
		expectedSOC  float64
	}{
		{"Constant Current Charge Stage, voltage below min", 2.0, 0.5, 0},
		{"Constant Current Charge Stage, voltage in mid-low range", 3.0, 0.5, 2.5},
		{"Constant Current Charge Stage, voltage in mid-high range", 3.7, 0.5, 15},
		{"Constant Current Charge Stage, voltage in max range", 4.0, 0.5, 80},
		{"Saturation Charge Stage", 4.0, 0.8, 60},
		{"Saturation Charge Stage, high current", 4.0, 1.0, 80},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			soc := calculator.SOC(test.voltage, test.current)
			if math.Abs(soc-test.expectedSOC) > 0.1 {
				t.Errorf("Expected SOC %.2f for voltage %.2f and current %.2f, got %.2f", test.expectedSOC, test.voltage, test.current, soc)
			}
		})
	}
}
