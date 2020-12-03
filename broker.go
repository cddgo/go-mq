package gomq

import (
	"errors"
	"sync"
	"time"
)

type BrokerImpl struct {
	exit chan bool
	cap  int

	//一个topic可以有多个订阅者，一个订阅者对应着一个通道
	topics map[string][]chan interface{}

	sync.RWMutex
}

func NewBroker() Broker {
	return &BrokerImpl{
		exit: make(chan bool),
		topics: make(map[string][]chan interface{}),
	}
}

func (b *BrokerImpl) publish(topic string, msg interface{}) error {
	select {
	case <-b.exit:
		return errors.New("broker closed")
	default:
	}
	b.RLock()
	subscribers, ok := b.topics[topic]
	b.RUnlock()
	if !ok {
		return errors.New(topic + " not found subscribers")
	}
	b.broadcast(msg, subscribers)
	return nil
}

func (b *BrokerImpl) subscribe(topic string) (<-chan interface{}, error) {
	select {
	case <-b.exit:
		return nil, errors.New("broker closed")
	default:
	}

	ch := make(chan interface{}, b.cap)
	b.Lock() // 对map的写必须加互斥锁，防止其他协程读
	b.topics[topic] = append(b.topics[topic], ch)
	b.Unlock()
	return ch, nil
}

func (b *BrokerImpl) unsubscribe(topic string, sub <-chan interface{}) error {
	select {
	case <-b.exit:
		return errors.New("broker closed")
	default:
	}

	b.RLock()
	subscribers, ok := b.topics[topic]
	b.RUnlock()
	if !ok {
		return errors.New(topic + " not found")
	}

	//删除订阅者
	for i, subscriber := range subscribers {
		if subscriber == sub {
			copy(subscribers[i:], subscribers[i+1:])
			subscribers = subscribers[:len(subscribers)-1]
			break
		}
	}

	// 对map的写必须加互斥锁，防止其他协程读
	b.Lock()
	b.topics[topic] = subscribers
	b.Unlock()
	return nil
}

func (b *BrokerImpl) close() {
	select {
	case <-b.exit:
		return
	default:
		close(b.exit)
		b.Lock()
		//这里主要是为了保证下一次使用该消息队列不发生冲突。
		b.topics = make(map[string][]chan interface{})
		b.Unlock()
		return
	}
}

func (b *BrokerImpl) setConditions(cap int) {
	b.cap = cap
}

func (b *BrokerImpl) broadcast(msg interface{}, subcribers []chan interface{}) {
	count := len(subcribers)

	concurrent := 0
	switch {
	case count > 1000:
		concurrent = 3
	case count > 100:
		concurrent = 2
	default:
		concurrent = 1
	}

	publish := func(start int) {
		// 注意这里向后移位goroutine的数量
		for j := start; j < count; j += concurrent {
			select {
			case subcribers[j] <- msg:
			case <-time.After(5 * time.Millisecond):
				//超时5ms后关闭向此订阅者发送，不保证百分百送达
			case <-b.exit:
				return
			}
		}
	}

	for i := 0; i < concurrent; i++ {
		go publish(i)
	}
}

func (b *BrokerImpl) GetPayload(sub <-chan interface{}) interface{} {
	for val := range sub {
		if val != nil {
			return val
		}
	}
	return nil
}
