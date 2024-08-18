package consumer

type Config struct {
	Brokers      string
	Topic        string
	Partition    int
	ConsumerName string
}
