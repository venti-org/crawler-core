package base

type QueueFlag struct {
	Push bool
	Pop  bool
}

type Queue interface {
	Init(QueueFlag) error
	Pop() interface{}
	Push(interface{})
	Close()
}

type Channel struct {
	C chan interface{}
}

func NewChannel(ch chan interface{}) *Channel {
	return &Channel{
		C: ch,
	}
}

func (c *Channel) Init(flag QueueFlag) error {
	return nil
}

func (c *Channel) Pop() interface{} {
	return <-c.C
}

func (c *Channel) Push(msg interface{}) {
	c.C <- msg
}

func (c *Channel) Close() {
	close(c.C)
}
