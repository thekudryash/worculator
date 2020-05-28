package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"
)

type WorkerInterface interface {
	Name() string
	Min() int
	Max() int
	Start()
	Stop()
	DeliveryRate() int
	AckRate() int

	CalculatorInterface
}

type workerInfo struct {
	instancesCount int
}

type worculator struct {
	workers map[string]*workerInfo
}

func (w worculator) calculate(wi WorkerInterface, workerHash string) int {
	recommendedWorkersInstances := wi.Calculate(
		wi.DeliveryRate(),
		wi.AckRate(),
		w.workers[workerHash].instancesCount,
	)

	if recommendedWorkersInstances > wi.Max() {
		recommendedWorkersInstances = wi.Max()
	} else if recommendedWorkersInstances < wi.Min() {
		recommendedWorkersInstances = wi.Min()
	}

	return recommendedWorkersInstances
}

func (w worculator) startWorker(wi WorkerInterface, workerHash string) {
	instancesCount := w.workers[workerHash].instancesCount
	go wi.Start()
	w.workers[workerHash].instancesCount = instancesCount + 1
}

func (w worculator) stopWorker(wi WorkerInterface, workerHash string) {
	instancesCount := w.workers[workerHash].instancesCount
	go wi.Stop()
	w.workers[workerHash].instancesCount = instancesCount - 1
}

func workerHash(wi WorkerInterface) string {
	return base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%v", wi)),
	)
}

var w = worculator{
	workers: make(map[string]*workerInfo),
}

func Manage(
	ctx context.Context,
	wi WorkerInterface,
) {
	workerHash := workerHash(wi)

	if w.workers[workerHash] == nil {
		w.workers[workerHash] = &workerInfo{
			instancesCount: 0,
		}
	}

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			recommendedWorkersCount := w.calculate(wi, workerHash)
			currentWorkersCount := w.workers[workerHash].instancesCount

			if currentWorkersCount <= recommendedWorkersCount {
				for i := currentWorkersCount; i < recommendedWorkersCount; i++ {
					w.startWorker(wi, workerHash)
				}
			} else if currentWorkersCount > recommendedWorkersCount {
				for i := currentWorkersCount; i > recommendedWorkersCount; i-- {
					w.stopWorker(wi, workerHash)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
