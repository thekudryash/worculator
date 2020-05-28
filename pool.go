package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"
)

type WorkerInterface interface {
	Context() (context.Context, context.CancelFunc)
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

type pool struct {
	workers map[string]*workerInfo
}

func (p pool) calculate(wi WorkerInterface, workerHash string) int {
	recommendedWorkersInstances := wi.Calculate(
		wi.DeliveryRate(),
		wi.AckRate(),
		p.workers[workerHash].instancesCount,
	)

	if recommendedWorkersInstances > wi.Max() {
		recommendedWorkersInstances = wi.Max()
	} else if recommendedWorkersInstances < wi.Min() {
		recommendedWorkersInstances = wi.Min()
	}

	return recommendedWorkersInstances
}

var p = pool{
	workers: make(map[string]*workerInfo, 10),
}

func (p pool) startWorker(wi WorkerInterface, workerHash string) {
	instancesCount := p.workers[workerHash].instancesCount
	go wi.Start()
	p.workers[workerHash].instancesCount = instancesCount + 1
}

func (p pool) stopWorker(wi WorkerInterface, workerHash string) {
	instancesCount := p.workers[workerHash].instancesCount
	go wi.Stop()
	p.workers[workerHash].instancesCount = instancesCount - 1
}

func workerHash(wi WorkerInterface) string {
	return base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%v", wi)),
	)
}

func Manage(
	ctx context.Context,
	wi WorkerInterface,
) {
	workerHash := workerHash(wi)

	if p.workers[workerHash] == nil {
		p.workers[workerHash] = &workerInfo{
			instancesCount: 0,
		}
	}

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			recommendedWorkersCount := p.calculate(wi, workerHash)
			currentWorkersCount := p.workers[workerHash].instancesCount

			if currentWorkersCount <= recommendedWorkersCount {
				for i := currentWorkersCount; i < recommendedWorkersCount; i++ {
					p.startWorker(wi, workerHash)
				}
			} else if currentWorkersCount > recommendedWorkersCount {
				for i := currentWorkersCount; i > recommendedWorkersCount; i-- {
					p.stopWorker(wi, workerHash)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
