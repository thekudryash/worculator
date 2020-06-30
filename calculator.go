package worculator

import (
	"math"
)

type CalculatorInterface interface {
	Calculate(deliveryRate int, ackRate int, instancesCount int) int
}

type DefaultCalculator struct{}

func (c DefaultCalculator) Calculate(
	deliveryRate int,
	ackRate int,
	instancesCount int,
) int {
	totalAckRate := instancesCount * ackRate

	if instancesCount == 0 {
		totalAckRate = 0
	}

	recommendedDiff := deliveryRate - totalAckRate

	// worker_instances = 10
	// ack_rate = 2
	// delivery_rate = 15
	//
	// total_ack_rate (ack_rate * worker_instances) = 2 * 10 = 20
	//
	// recommended_rate (delivery_rate - total_ack_rate) = 15 - 20 = -5
	//
	// calucalate_result (worker_instances + (recommended_rate / ack_rate))
	// |-> 10 + (-5 / 2) = 10 - 2 = 8
	instancesDiff := int(math.Ceil(
		float64(recommendedDiff) / float64(ackRate),
	))

	return instancesCount + instancesDiff
}

var _ CalculatorInterface = (*DefaultCalculator)(nil)
