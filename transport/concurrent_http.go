package transport

import (
	j "github.com/uber/jaeger-client-go/thrift-gen/jaeger"
)

type ConcurrentHttpTransport struct {
	*HTTPTransport
	pool chan int
}

func NewConcurrentHttpTransport(httpTransport *HTTPTransport) *ConcurrentHttpTransport {
	ctp := &ConcurrentHttpTransport{
		HTTPTransport: httpTransport,
	}
	ctp.pool = make(chan int, 100)
	for i := 0; i < 100; i++ {
		ctp.pool <- 0
	}
	return ctp
}

func (c *ConcurrentHttpTransport) send(spans []*j.Span) error {
	<-c.pool
	defer func() {
		c.pool <- 0
	}()
	return c.HTTPTransport.send(spans)
}
