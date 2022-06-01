package rabbitmqlib

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/imaimang/mgo/idlib"
	"github.com/streadway/amqp"
)

//Connection Connection
type Connection struct {
	channels     map[string]*Channel
	lockChannels *sync.Mutex
	conn         *amqp.Connection
	mqAddress    string
	closeChan    chan *amqp.Error
	connected    bool
}

//Connect Connect
func (c *Connection) Connect(mqAddress string) {
	c.channels = make(map[string]*Channel, 0)
	c.lockChannels = new(sync.Mutex)
	c.mqAddress = mqAddress
	c.connect()
}

//GetChannelCount 获取Channel数量
func (c *Connection) GetChannelCount() int {
	return len(c.channels)
}

func (c *Connection) connect() {
	var err error
	for {
		c.conn, err = amqp.Dial(c.mqAddress)
		if err == nil {
			break
		}
		fmt.Println("connect rabbitmq server faild,reconnect...", err)
		time.Sleep(2 * time.Second)
	}
	fmt.Println("connect rabbitmq server success", c.mqAddress)
	c.connected = true
	c.closeChan = make(chan *amqp.Error, 0)
	c.conn.NotifyClose(c.closeChan)
	go c.monitorClose()
}

func (c *Connection) monitorClose() {
	for {
		err, ok := <-c.closeChan
		if !ok {
			break
		}
		fmt.Println("rabbitmq already closed", err, ok)
		c.lockChannels.Lock()
		c.connected = false
		c.lockChannels.Unlock()
		c.conn.Close()
		c.connect()
		c.lockChannels.Lock()
		c.connected = true
		c.lockChannels.Unlock()
	}
}

//GetChannel GetChannel
func (c *Connection) GetChannel() (*Channel, error) {
	c.lockChannels.Lock()
	var channel *Channel
	var err error
	if c.connected {
		index := 0
		for {
			for _, item := range c.channels {
				if !item.IsUsed {
					item.IsUsed = true
					channel = item
					break
				}
			}
			if channel != nil {
				break
			} else if channel == nil && len(c.channels) < 100 {
				var ch *amqp.Channel
				ch, err = c.conn.Channel()
				if err == nil {
					channel = new(Channel)
					channel.CH = ch
					channel.IsUsed = true
					channel.Identify = idlib.CreateMD5()
					c.channels[channel.Identify] = channel
				} else {
					fmt.Println("创建Channel失败", err, len(c.channels))
				}
				break
			} else if channel == nil && index == 100 {
				err = errors.New("channel busy")
				break
			}
			index++
			time.Sleep(10 * time.Millisecond)
		}
	} else {
		err = errors.New("rabbitmq connect break")
	}
	c.lockChannels.Unlock()
	return channel, err
}

//CloseChannel CloseChannel
func (c *Connection) CloseChannel(channel *Channel) {
	if channel != nil {
		channel.CH.Close()
		c.lockChannels.Lock()
		delete(c.channels, channel.Identify)
		c.lockChannels.Unlock()
	}
}
