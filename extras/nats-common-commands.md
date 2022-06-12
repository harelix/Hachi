# Streams
start nats stream using docker 
```bash 
docker run -d --name nats-hachi -p 4222:4222 -p 6222:6222 -p 8222:8222 nats -js
```
create stream from configuration file
```bash
nats str add controller --config orders.json
```
[Streams commands](https://docs.nats.io/running-a-nats-service/nats_admin/jetstream_admin/streams)

create stream + cli args
```bash
nats stream add neurostream --subjects "cns.brain.>" --ack --max-msgs=-1 --max-bytes=-1 --max-age=1y --storage file --retention limits --max-msg-size=-1 --discard=old
```

# Consumers

### Creating Pull-Based Consumers
add named consumer 'brain' to 'neurostream' stream:
```bash
nats consumer add neurostream brain_consumer --filter 'cns.brain.*' --ack explicit --pull --deliver all --max-deliver=-1 --sample 100
```

listing consumers on streams:
```bash
nats con ls {stream name}
```

One can store the configuration in a JSON file, the format of this is the same as:
```bash
nats con info ORDERS HACHI -j | jq .config > hachi-con.config
```

### listener on consumer

```bash 
nats consumer next neurostream brain_consumer --count 1000
```

```bash
nats pub neurostream --subject cns.brain.neurolink --count=10 --sleep 1s "publication #{{Count}} @ {{TimeStamp}}"
```



