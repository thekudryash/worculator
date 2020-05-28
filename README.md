# Worculator

## TODO
- [ ] Auto-grow workers storage
- [ ] Tests for pool
- [ ] RabbitMQ adapter

## Example
```go
type worker struct {
	Ctx       context.Context
	CancelCtx context.CancelFunc
	InputChannel chan string

	DefaultCalculator
}

func (w *worker) Context() (context.Context, context.CancelFunc) { return w.Ctx, w.CancelCtx }
func (w *worker) Name() string                                   { return strconv.Itoa(rand.Int()) }
func (w *worker) Min() int                                       { return 1 }
func (w *worker) Max() int                                       { return 10 }
func (w *worker) Start() {
	w.Ctx, w.CancelCtx = context.WithCancel(context.Background())

	for {
		select {
		case message := <-w.InputChannel:
			fmt.Println(message)
		case <-w.Ctx.Done():
			return
		}
	}
}
func (w *worker) Stop()             { w.CancelCtx() }
func (w *worker) DeliveryRate() int { return 2 }
func (w *worker) AckRate() int      { return 1 }

var _ WorkerInterface = (*worker)(nil)

func main() {
	workerInputChannel := make(chan string, 1000)
	worker1 := &worker{
		InputChannel: workerInputChannel,
	}

	poolCtx, cancelPoolCtx := context.WithCancel(context.Background())
	go pool.Manage(poolCtx, worker1)

	c, _ := context.WithTimeout(context.Background(), 10*time.Second)
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)

		for {
			select {
			case <-ticker.C:
				workerInputChannel <- "string"
			}
		}
	}()

	<-c.Done()
	cancelPoolCtx()
}

```