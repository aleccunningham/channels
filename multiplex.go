package channels

// FAN OUT
numFinders := runtime.NumCPU()
finders := make([]<-chan int, numFinders)
for i := 0; i < numFinders; i++ {
	finders[i] = primeFinder(done, randIntStream)
}

// FAN IN
fanIn := func(
	done <-chan interface{},
	channels ...<-chan interface{},
) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})

	// when passed a channel, multiplex will read from
	// the channnel and pass the read value onto
	// the multiplexedStream channel
	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	// select from all the channels
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	// wait for all reads to complete
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

return multiplexedStream
