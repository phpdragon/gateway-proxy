package rabbit

/*
 * @Script: mq.go
 * Author: Pangxiaobo
 * Last Modified: Saturday December 8th 2018 1:46:07 pm
 * Modified By: the developer formerly known as Pangxiaobo at <10846295@qq.com>
 * @Last Modified By: pangxiaobo
 * @Last Modified At: 2018-12-12 14:24:22
 */

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"io"
	"log"
	"net/http"
)

var amqpUri string

type MqController struct{}

// MessageEntity Entity for HTTP Request Body: Message/Exchange/Queue/QueueBind JSON Input
type MessageEntity struct {
	Exchange     string `json:"exchange"`
	Key          string `json:"key"`
	DeliveryMode uint8  `json:"deliverymode"`
	Priority     uint8  `json:"priority"`
	Body         string `json:"body"`
}

type ExchangeEntity struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Durable    bool   `json:"durable"`
	AutoDelete bool   `json:"autodelete"`
	NoWait     bool   `json:"nowait"`
}

type QueueEntity struct {
	Name       string `json:"name"`
	Durable    bool   `json:"durable"`
	AutoDelete bool   `json:"autodelete"`
	Exclusive  bool   `json:"exclusive"`
	NoWait     bool   `json:"nowait"`
}

type QueueBindEntity struct {
	Queue    string   `json:"queue"`
	Exchange string   `json:"exchange"`
	NoWait   bool     `json:"nowait"`
	Keys     []string `json:"keys"` // bind/routing keys
}

// MqClient Operate Wrapper
type MqClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan error
}

type Options struct {
	Host     string
	User     string
	Password string
}

func NewClient(config *Options) *MqClient {
	amqpUri = fmt.Sprintf("amqp://%s:%s@%s", config.User, config.Password, config.Host)
	c := MqClient{}
	err := c.connect()
	if err != nil {
		return nil
	}
	return &c
}

func (r *MqClient) connect() (err error) {
	r.conn, err = amqp.Dial(amqpUri)
	if err != nil {
		log.Printf("[amqp] connect errorcode: %s\n", err)
		return err
	}
	r.channel, err = r.conn.Channel()
	if err != nil {
		log.Printf("[amqp] get channel errorcode: %s\n", err)
		return err
	}
	r.done = make(chan error)
	return nil
}

func (r *MqClient) Publish(exchange, key string, deliverymode, priority uint8, body string) (err error) {
	err = r.channel.Publish(exchange, key, false, false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			DeliveryMode:    deliverymode,
			Priority:        priority,
			Body:            []byte(body),
		},
	)
	if err != nil {
		log.Printf("[amqp] publish message errorcode: %s\n", err)
		return err
	}
	return nil
}

func (r *MqClient) DeclareExchange(name, typ string, durable, autodelete, nowait bool) (err error) {
	err = r.channel.ExchangeDeclare(name, typ, durable, autodelete, false, nowait, nil)
	if err != nil {
		log.Printf("[amqp] declare exchange errorcode: %s\n", err)
		return err
	}
	return nil
}

func (r *MqClient) DeleteExchange(name string) (err error) {
	err = r.channel.ExchangeDelete(name, false, false)
	if err != nil {
		log.Printf("[amqp] delete exchange errorcode: %s\n", err)
		return err
	}
	return nil
}

func (r *MqClient) DeclareQueue(name string, durable, autodelete, exclusive, nowait bool) (err error) {
	_, err = r.channel.QueueDeclare(name, durable, autodelete, exclusive, nowait, nil)
	if err != nil {
		log.Printf("[amqp] declare queue errorcode: %s\n", err)
		return err
	}
	return nil
}

func (r *MqClient) DeleteQueue(name string) (err error) {
	// TODO: other property wrapper
	_, err = r.channel.QueueDelete(name, false, false, false)
	if err != nil {
		log.Printf("[amqp] delete queue errorcode: %s\n", err)
		return err
	}
	return nil
}

func (r *MqClient) BindQueue(queue, exchange string, keys []string, nowait bool) (err error) {
	for _, key := range keys {
		if err = r.channel.QueueBind(queue, key, exchange, nowait, nil); err != nil {
			log.Printf("[amqp] bind queue errorcode: %s\n", err)
			return err
		}
	}
	return nil
}

func (r *MqClient) UnBindQueue(queue, exchange string, keys []string) (err error) {
	for _, key := range keys {
		if err = r.channel.QueueUnbind(queue, key, exchange, nil); err != nil {
			log.Printf("[amqp] unbind queue errorcode: %s\n", err)
			return err
		}
	}
	return nil
}

func (r *MqClient) ConsumeQueue(queue string, message chan []byte) (err error) {
	deliveries, err := r.channel.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("[amqp] consume queue errorcode: %s\n", err)
		return err
	}
	go func(deliveries <-chan amqp.Delivery, done chan error, message chan []byte) {
		for d := range deliveries {
			message <- d.Body
		}
		done <- nil
	}(deliveries, r.done, message)
	return nil
}

func (r *MqClient) Close() (err error) {
	err = r.conn.Close()
	if err != nil {
		log.Printf("[amqp] close errorcode: %s\n", err)
		return err
	}
	return nil
}

// QueueHandler HTTP Handlers
func (m *MqController) QueueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" || r.Method == "DELETE" {

		if r.Body == nil {
			fmt.Println("missing form body")
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		entity := new(QueueEntity)
		if err = json.Unmarshal(body, entity); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rabbit := new(MqClient)
		if err = rabbit.connect(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rabbit.Close()

		if r.Method == "POST" {
			if err = rabbit.DeclareQueue(entity.Name, entity.Durable, entity.AutoDelete, entity.Exclusive, entity.NoWait); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("declare queue ok"))
		} else if r.Method == "DELETE" {
			if err = rabbit.DeleteQueue(entity.Name); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("delete queue ok"))
		}
	} else if r.Method == "GET" {
		r.ParseForm()
		rabbit := new(MqClient)
		if err := rabbit.connect(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rabbit.Close()

		message := make(chan []byte)

		for _, name := range r.Form["name"] {
			if err := rabbit.ConsumeQueue(name, message); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Write([]byte(""))
		w.(http.Flusher).Flush()

		for {
			fmt.Printf("Received message %s\n", <-message)
			//fmt.Fprintf(w, "%s\n", <-message)
			w.(http.Flusher).Flush()
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (m *MqController) QueueBindHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" || r.Method == "DELETE" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		entity := new(QueueBindEntity)
		if err = json.Unmarshal(body, entity); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rabbit := new(MqClient)
		if err = rabbit.connect(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rabbit.Close()

		if r.Method == "POST" {
			if err = rabbit.BindQueue(entity.Queue, entity.Exchange, entity.Keys, entity.NoWait); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("bind queue ok"))
		} else if r.Method == "DELETE" {
			if err = rabbit.UnBindQueue(entity.Queue, entity.Exchange, entity.Keys); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("unbind queue ok"))
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (m *MqController) PublishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		entity := new(MessageEntity)
		if err = json.Unmarshal(body, entity); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rabbit := new(MqClient)
		if err = rabbit.connect(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rabbit.Close()

		if err = rabbit.Publish(entity.Exchange, entity.Key, entity.DeliveryMode, entity.Priority, entity.Body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("publish message ok => " + entity.Body))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (m *MqController) ExchangeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" || r.Method == "DELETE" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		entity := new(ExchangeEntity)
		if err = json.Unmarshal(body, entity); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rabbit := new(MqClient)
		if err = rabbit.connect(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rabbit.Close()

		if r.Method == "POST" {
			if err = rabbit.DeclareExchange(entity.Name, entity.Type, entity.Durable, entity.AutoDelete, entity.NoWait); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("declare exchange ok"))
		} else if r.Method == "DELETE" {
			if err = rabbit.DeleteExchange(entity.Name); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("delete exchange ok"))
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
