package messagebroker

import (
	"context"
	"fmt"
	"sync"

	"github.com/GSH-LAN/Unwindia_common/src/go/logger"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/gammazero/workerpool"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func init() {
	log = logger.GetSugaredLogger()
}

type MessageBroker interface {
	StartConsumer(topic string) error
	StartProducer(topic string) (pulsar.Producer, error)
	SendMessage(topic string, message []byte) (*pulsar.ProducerMessage, error)
}

type PulsarWorkerImpl struct {
	client        *PulsarClient
	serviceId     string
	baseTopic     string
	lastState     map[string]*Response
	lock          sync.Mutex
	workerpool    *workerpool.WorkerPool
	consumers     map[string]*pulsar.Consumer
	consumersLock sync.Mutex
	producers     map[string]pulsar.Producer
	producersLock sync.Mutex
}

func NewWorker(serviceId, baseTopic string, client *PulsarClient, workerpool *workerpool.WorkerPool) MessageBroker {
	return &PulsarWorkerImpl{
		serviceId:  serviceId,
		baseTopic:  baseTopic,
		client:     client,
		workerpool: workerpool,
		lastState:  make(map[string]*Response),
		consumers:  make(map[string]*pulsar.Consumer),
		producers:  make(map[string]pulsar.Producer),
	}
}

func (p *PulsarWorkerImpl) StartConsumer(topic string) error {
	// forever
	func() {
		consumer, err := p.client.client.Subscribe(pulsar.ConsumerOptions{
			Topic:            fmt.Sprintf("%s/%s", p.baseTopic, topic),
			SubscriptionName: p.serviceId,
		})

		if err != nil {
			log.Fatal(err)
		}

		defer consumer.Close()

		ctx := context.Background()

		for {
			msg, err := consumer.Receive(ctx)
			if err != nil {
				log.Error(err)
			} else {
				log.Infof("Received message : %v\n", string(msg.Payload()))
			}

			consumer.Ack(msg)
		}
	}()
	log.Infof("Started Pulsar consumer-worker [%s]", topic)
	return nil
}

func (p *PulsarWorkerImpl) updateState(id string) {
	defer p.lock.Unlock()
	p.lock.Lock()
}

func (p *PulsarWorkerImpl) StartProducer(topic string) (pulsar.Producer, error) {
	if producer, exists := p.producers[topic]; exists {
		return producer, fmt.Errorf("producer %s already exists", topic)
	}

	producer, err := p.client.client.CreateProducer(pulsar.ProducerOptions{
		Topic: fmt.Sprintf("%s/%s", p.baseTopic, topic),
	})

	p.addProducer(topic, producer)

	log.Infof("Started Pulsar producer-worker [%s]", topic)
	return producer, err
}

func (p *PulsarWorkerImpl) addProducer(id string, producer pulsar.Producer) error {
	defer p.producersLock.Unlock()
	p.producersLock.Lock()
	if _, exists := p.producers[id]; exists {
		return fmt.Errorf("producer %s already exists", id)
	}
	p.producers[id] = producer
	return nil
}

func (p *PulsarWorkerImpl) SendMessage(topic string, message []byte) (*pulsar.ProducerMessage, error) {
	producer, producerExists := p.producers[topic]
	if !producerExists {
		var err error
		producer, err = p.StartProducer(topic)
		if err != nil {
			return nil, err
		}
	}

	msg := pulsar.ProducerMessage{
		Payload: message,
	}

	messageID, err := producer.Send(context.Background(), &msg)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Infof("Successfully published to topic %s message ID %v", topic, messageID)
	return &msg, nil
}
