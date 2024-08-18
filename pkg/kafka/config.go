package kafka

type Config struct {
	Brokers      string
	Topic        string
	Partition    int
	WriteTimeout int
}

type Message struct {
	EventID   string
	EventName string
	Data      interface{}
}
