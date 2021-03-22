package rabbitmq

import (
	"fmt"
	"net"
	"time"

	"github.com/streadway/amqp"
	"k8s.io/klog/v2"
)

type Messaging struct {
	Url       string
	QueueName string
	Timeout   time.Duration

	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func NewMessaging(url string, queueName string, timeout time.Duration) *Messaging {
	return &Messaging{Url: url, QueueName: queueName, Timeout: timeout}
}

func (m *Messaging) Connect() (err error) {
	if err = m.retry(3, m.Timeout, m.dial); err != nil {
		return err
	}

	if m.ch, err = m.conn.Channel(); err != nil {
		return err
	}

	if m.q, err = m.ch.QueueDeclare(
		m.QueueName, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	); err != nil {
		m.Disconnect()
		return err
	}

	return nil
}

func (m *Messaging) dial() (err error) {
	if m.conn, err = amqp.DialConfig(m.Url, amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, m.Timeout*time.Second)
		},
	}); err != nil {
		return err
	}

	return nil
}

func (m *Messaging) retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(sleep * time.Second)

		klog.InfoS("retrying after error:[%s]", err.Error())
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

func (m *Messaging) Consume() (<-chan amqp.Delivery, error) {
	return m.ch.Consume(
		m.q.Name, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
}

func (m *Messaging) Publish(contentType string, data []byte) error {
	if err := m.ch.Publish(
		"",
		m.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        data,
		},
	); err != nil {
		return err
	}

	return nil
}

func (m *Messaging) Disconnect() error {
	m.ch.Close()
	m.conn.Close()

	return nil
}
