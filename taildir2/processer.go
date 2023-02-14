package taildir

type LineProcesser func(*TailD, <-chan *Line)
