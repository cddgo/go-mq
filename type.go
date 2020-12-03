package gomq

type Broker interface {
	//消息推送，传入订阅的主题、要传递的消息
	publish(topic string, msg interface{}) error
	//消息订阅，传入订阅的主题和对应的通道
	subscribe(topic string) (<-chan interface{}, error)
	//取消订阅，传入订阅的主题和对应的通道
	unsubscribe(topic string, sub <-chan interface{}) error

	//关闭消息队列
	close()

	//这个属于内部方法，作用是进行广播，对推送的消息进行广播，保证每一个订阅者都可以收到
	broadcast(msg interface{}, subscribers []chan interface{})
	//这里是用来设置条件，条件就是消息队列的容量，这样我们就可以控制消息队列的大小了
	setConditions(cap int)
}
