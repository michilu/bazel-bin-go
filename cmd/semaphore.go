package cmd

var (
	semaphore chan struct{}
)

func initSem() {
	semaphore = make(chan struct{}, flag.parallel)
}
