package environment

type BaseEnvironment struct {
	LogLevel         string `env:"LOG_LEVEL" envDefault:"INFO"`
	ServiceID        string `env:"SERVICE_ID" envDefault:"UNWINDIA_GENERIC_SERVICE"`
	WorkerCount      int    `env:"WORKER_COUNT" envDefault:"-1" envDescription:"Number of workers to use. If 0, uses the number of CPUs, if -1 it uses the number of CPUs - 1"`
	HTTPPort         int    `env:"HTTP_PORT" envDefault:"8080"`
	WorkItemLockType string `env:"WORKITEM_LOCK" envDefault:"memory" envDescription:"Type of workitem lock to use. Valid values are 'memory' and 'mongodb'"`

	MongoDbURI string `env:"MONGODB_URI"`

	ConfigFileName     string `env:"CONFIG_FILENAME"`
	ConfigTemplatesDir string `env:"CONFIG_TEMPLATE_DIR"`

	PulsarURL        string `env:"PULSAR_URL" envDefault:"pulsar://localhost:6650"`
	PulsarBaseTopic  string `env:"PULSAR_TOPIC" envDefault:"persistent://unwindia/unwindia"`
	PulsarAuth       string `env:"PULSAR_AUTH" envDefault:"simple"`
	PulsarAuthParams string `env:"PULSAR_AUTH_PARAMS" required:"true"`
}
