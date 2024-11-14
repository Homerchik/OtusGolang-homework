package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		out := make(Bi)
		go wrapper(in, done, out)

		in = stage(out)
	}
	return in
}

func wrapper(in In, done In, out Bi) {
	defer close(out)
	for {
		select {
		case <-done:
			// It appears that without this one TestAllStageStop stucks on first stages,
			//  trying to send additional data to the channel no one reads
			<-in
			return
		case i, ok := <-in:
			if !ok {
				return
			}
			select {
			case <-done:
				<-in
				return
			case out <- i:
			}
		}
	}
}
