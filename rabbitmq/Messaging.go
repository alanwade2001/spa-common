package rabbitmq

import (
	"github.com/streadway/amqp"
)

type Messaging struct {
	Url       string
	QueueName string

	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func NewMessaging(url string, queueName string) *Messaging {
	return &Messaging{Url: url, QueueName: queueName}
}

func (m *Messaging) Connect() (err error) {
	if m.conn, err = amqp.Dial(m.Url); err != nil {
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
