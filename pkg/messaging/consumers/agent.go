package consumers

/*
func ConstructConsumerConfig(agent config.IAgent) []nats.ConsumerConfig {

	consumers := make([]nats.ConsumerConfig, len(agent.GetIdentifiers()))
	for idx, identifier := range agent.GetIdentifiers() {

		//todo: test: can we subscribe to multiple subjects on the same consumer
		consumers[idx] = nats.ConsumerConfig{
			Durable:       "true",
			FilterSubject: identifier,
			Description:   string(agent.GetType().String()) + " durable consumer",
			AckPolicy:     nats.AckExplicitPolicy,
			AckWait:       0,
			MaxDeliver:    3,
			MaxWaiting:    0,
			HeadersOnly:   false,
		}
	}
	return consumers
}
*/
