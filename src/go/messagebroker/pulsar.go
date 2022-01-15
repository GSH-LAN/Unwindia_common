package messagebroker

import (
	"io/ioutil"
	"time"

	"github.com/GSH-LAN/Unwindia_common/src/go/logger"
	zaphook "github.com/Sytten/logrus-zap-hook"
	"github.com/apache/pulsar-client-go/pulsar"
	pulsarlog "github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/gammazero/workerpool"
	"github.com/sirupsen/logrus"
)

var logrusZapBridge *logrus.Logger

func init() {
	if log == nil {
		log = logger.GetSugaredLogger()
	}

	logr := logrus.New()
	logr.ReportCaller = true       // So Zap reports the right caller
	logr.SetOutput(ioutil.Discard) // Prevent logrus from writing its logs

	hook, _ := zaphook.NewZapHook(log.Desugar())

	logr.Hooks.Add(hook)
	logrusZapBridge = logr
}

type PulsarClient struct {
	client     pulsar.Client
	workerpool *workerpool.WorkerPool
}

func NewClient(pulsarUrl string, workerpool *workerpool.WorkerPool) (*PulsarClient, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               pulsarUrl,
		OperationTimeout:  5 * time.Second,
		ConnectionTimeout: 5 * time.Second,
		// Logger:            log.NewLoggerWithLogrus(log.GetLogger()),
		Logger: pulsarlog.NewLoggerWithLogrus(logrusZapBridge),
	})

	if err != nil {
		return nil, err
	}

	return &PulsarClient{
		client:     client,
		workerpool: workerpool,
	}, nil
}
