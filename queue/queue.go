package queue

type Queue interface {
	Pop() interface{}
	Push(interface{})
}

type Channel struct {
	C chan interface{}
}

func NewChannel(ch chan interface{}) *Channel {
	return &Channel{
		C: ch,
	}
}

func (c *Channel) Pop() interface{} {
	for msg := range c.C {
		return msg
	}
	return nil
}

func (c *Channel) Push(msg interface{}) {
	c.C <- msg
}
