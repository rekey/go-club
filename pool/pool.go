package pool

import (
	"log"
)

type Pool struct {
	task       chan int
	result     chan int
	size       int
	consumerFn func(int)
	producerFn func(int)
	needLog    bool
}

func NewPool(size int, consumerFn func(i int), needLog bool) *Pool {
	return &Pool{
		task:       make(chan int, size),
		result:     make(chan int, size),
		size:       size,
		consumerFn: consumerFn,
		producerFn: func(i int) {},
		needLog:    needLog,
	}
}

func NewPoolWithResult(size int, consumerFn func(i int), producerFn func(i int), needLog bool) *Pool {
	return &Pool{
		task:       make(chan int, size),
		result:     make(chan int, size),
		size:       size,
		consumerFn: consumerFn,
		producerFn: producerFn,
		needLog:    needLog,
	}
}

func (p *Pool) consumer() {
	for w := 1; w <= p.size; w++ {
		go func() {
			for id := range p.task {
				p.consumerRun(id)
			}
		}()
	}
}

func (p *Pool) consumerRun(id int) {
	// p.log("consumer", id, "start")
	p.consumerFn(id)
	p.result <- id
}

func (p *Pool) producer() {
	for w := 1; w <= p.size; w++ {
		go func() {
			for id := range p.result {
				p.producerRun(id)
			}
		}()
	}
}

func (p *Pool) producerRun(id int) {
	// p.log("producer", id, "start")
	p.producerFn(id)
	p.task <- id
}

func (p *Pool) log(v ...any) {
	if !p.needLog {
		return
	}
	log.Println(v...)
}

func (p *Pool) Run() {
	p.consumer()
	p.log("for consumer", "init", p.size, len(p.task))
	p.producer()
	p.log("for producer", "init", p.size, len(p.result))
	// 初始化任务
	for j := 0; j < p.size; j++ {
		p.producerRun(j)
	}
}
