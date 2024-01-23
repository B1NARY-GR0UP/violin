package violin

// TODO: refer to ants

var defaultViolin *Violin

func init() {
	defaultViolin = New(WithMaxWorkers(DefaultViolinMaxWorkerNum), WithWorkerIdleTimeout(DefaultViolinWorkerIdleTimeout))
}

func Submit(task func()) error {
	//return defaultViolin.Submit(task)
	return nil
}

func SubmitWait(task func()) error {
	//return defaultViolin.SubmitWait(task)
	return nil
}

func WorkerNum() uint32 {
	return defaultViolin.WorkerNum()
}

func MaxWorkerNum() int {
	return defaultViolin.MaxWorkerNum()
}
