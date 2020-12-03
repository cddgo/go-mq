package gomq

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestOneTopic(t *testing.T) {
	client := NewClient()
	client.SetConditions(10)
	defer client.Close()

	ch, err := client.Subscribe("topic")
	if err != nil {
		t.Fatal("subscribe failed: ", err)
	}

	// 定时发送
	go func(c *Client) {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := c.Publish("topic", "make a mq")
				if err != nil {
					fmt.Println("publish msg error")
				}
			default:
			}
		}
	}(client)

	// 循环接收
	go func(ch <-chan interface{}, c *Client) {
		for {
			val := c.GetPayload(ch)
			fmt.Println("get message: ", val)
		}
	}(ch, client)

	time.Sleep(10 * time.Second)
}

func TestManyTopic(t *testing.T) {
	client := NewClient()
	client.SetConditions(10)
	defer client.Close()

	for i := 0; i < 3; i++ {
		topic := fmt.Sprintf("devhui_%02d", i)
		go sub(client, topic)
	}

	go pub(client)
	time.Sleep(10 * time.Second)
}

func sub(c *Client, topic string) {
	ch, err := c.Subscribe(topic)
	if err != nil {
		fmt.Printf("subscribe :%s, err:%v\n", topic, err)
	}
	for {
		if payload := c.GetPayload(ch); payload != nil {
			fmt.Printf("get %v, from %v\n", payload, topic)
		}
	}
}

func pub(c *Client) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			for i := 0; i < 3; i++ {
				topic := fmt.Sprintf("devhui_%02d", i)
				payload := fmt.Sprintf("make a mq_%02d", i)
				if err := c.Publish(topic, payload); err != nil {
					fmt.Println("publish msg err：", err)
				}
			}
		}
	}
}

func TestNewClient(t *testing.T) {
	//我们向不同的topic发送不同的信息，当订阅者收到消息后，就行取消订阅。
	client := NewClient()
	client.SetConditions(2)

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		topic := fmt.Sprintf("id-%d", i)
		payload := fmt.Sprintf("devhui-%d", i)
		ch, err := client.Subscribe(topic)
		if err != nil {
			t.Fatal(err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			e := client.GetPayload(ch)
			fmt.Println(e)
			if e != payload {
				t.Fatalf("%s expected but recv %s", payload, e)
			}

			if err := client.Unsubscribe(topic, ch); err != nil {
				t.Fatal(err)
			}
		}()

		if err = client.Publish(topic, payload); err != nil {
			t.Fatal(err)
		}
	}

	wg.Wait()
}
