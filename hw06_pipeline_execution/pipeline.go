package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	go wrapper(in, done, out)
	in = out
	for _, stage := range stages {
		out := make(Bi)
		go wrapper(in, done, out)

		in = stage(out)
	}
	return in
}

func wrapper(in In, done In, out Bi) {
	defer func() {
		//nolint:revive
		for range in {
		}
		close(out)
	}()
	for {
		select {
		case <-done:
			return
		case i, ok := <-in:
			if !ok {
				return
			}
			select {
			case <-done:
				return
			case out <- i:
			}
		}
	}
}
