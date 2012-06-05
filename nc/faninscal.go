package nc

import ()

// ScalarChan is like a chan float32, but with fan-out
// (replicate data over multiple output channels).
// It can only have one input side though.
type FanInScal struct {
	fanout []chan float32
}

func MakeFanInScal() (c FanInScal) {
	return
}

// Add a new fanout and return it.
// All fanouts should be created before using the channel.
func (v *FanInScal) Fanout(buf int) FanoutScal {
	v.fanout = append(v.fanout, make(chan float32, buf))
	return v.fanout[len(v.fanout)-1]
}

// Send operator.
func (v *FanInScal) Send(data float32) {
	if len(v.fanout) == 0 {
		panic("FanInScal.Send: no fanout")
	}
	for i := range v.fanout {
		v.fanout[i] <- data
	}
}