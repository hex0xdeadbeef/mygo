package chapter6

type Counter struct {
	currentValue int
}

func (c *Counter) GetValue() int {
	return c.currentValue
}

func (c *Counter) Increment() {
	c.currentValue++
}

func (c *Counter) Reset() {
	c.currentValue = 0
}
