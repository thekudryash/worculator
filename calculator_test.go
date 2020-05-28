package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultCalculator_Calculate(t *testing.T) {
	type fields struct {
		ackRate                     int
		deliveryRate                int
		currentWorkerInstancesCount int
	}
	type want struct {
		calculateResult int
	}
	type test struct {
		name string
		fields
		want
	}

	tests := []test{
		{
			name: "positive difference between total ack with zero worker instances",
			fields: fields{
				ackRate:                     10,
				deliveryRate:                5,
				currentWorkerInstancesCount: 0,
			},
			want: want{
				calculateResult: 1,
			},
		},
		{
			name: "negative difference ack - delivery with non zero worker instances",
			fields: fields{
				ackRate:                     10,
				deliveryRate:                5,
				currentWorkerInstancesCount: 2,
			},
			want: want{
				calculateResult: 1,
			},
		},
		{
			name: "negative difference ack - delivery with zero worker instances",
			fields: fields{
				ackRate:                     10,
				deliveryRate:                50,
				currentWorkerInstancesCount: 0,
			},
			want: want{
				calculateResult: 5,
			},
		},
		{
			name: "negative difference ack - delivery with non zero worker instances",
			fields: fields{
				ackRate:                     10,
				deliveryRate:                50,
				currentWorkerInstancesCount: 10,
			},
			want: want{
				calculateResult: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultCalculatorResult := DefaultCalculator{}.Calculate(
				tt.deliveryRate,
				tt.ackRate,
				tt.currentWorkerInstancesCount,
			)

			assert.Equal(t, tt.calculateResult, defaultCalculatorResult)
		})
	}
}
