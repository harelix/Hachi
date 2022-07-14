## Hachi 
### A Cloud Native polymorphic connectivity mesh

controller & The agents

> hachi 八 means "Eight" or Bee in Japanese


### <span style="background-color: #EFEFEF;color:#000;padding:0.2em">stream</span> Parameters 
 - <span style="background-color: #EFEFEF;color:green;padding:0.2em">async</span> (bool)
 - <span style="background-color: #EFEFEF;color:green;padding:0.2em">verb</span> (string)
   At this time, Hachi supports the POST and GET (static body payload) verbs
 - <span style="background-color: #EFEFEF;color:green;padding:0.2em">selectors</span> (pattern, all, dynamic)
   - pattern
   - all
   - dynamic
 - <span style="background-color: #EFEFEF;color:green;padding:0.2em">async</span> (bool)
 - <span style="background-color: #EFEFEF;color:green;padding:0.2em">async</span> (bool)
 
```hcl
 stream "trigger_webhook_const_event" {
      async = true
      verb = "POST"
      selectors {
        pattern {
          values = ["agents.selector.{{.route::selector}}"]
        }
      }

      local = "/selector/:selector"
      remote {
        webhook {
          event = "EVENT.DATA.SOME_DATA_CHANGE"
        }
      }
    }
```

Secure, Performant, Agile, Resilient

With reliable backbone services (NATS, Kafka*, Rabbit*), adaptive communication channels (personal and public), and leaf nodes customization, optimize communications for various invokation scenarios.
Hachi's Adaptive mesh Architecture allows for a perfect fit for unique needs to introduce and activate devices, edge, cloud, or hybrid deployments.

Hachi with its underlying backbone supports true multi-tenancy, securely isolate and share data. 
Security is bifurcated from topology, connect anywhere in a deployment and it will do the right thing - Based on Nats.

Hachi self-heals and can scale up, down with zero downtime. 
NATS topology allowing Hachi to future proof the system and meet the needs of today and tomorrow.

## Encryption at Rest
[JetStream encryption at rest.](https://docs.nats.io/running-a-nats-service/nats_admin/jetstream_admin/encryption_at_rest)
---
NATS single dev server and tools installation
 
#### just to get you up and running in now time
```bash
docker run -d --name nats-hachi -p 4222:4222 -p 6222:6222 -p 8222:8222 nats -js
```

### cli tools
[nsc CLI installation instructions](https://docs.nats.io/using-nats/nats-tools/nsc)

[nats CLI installation instructions](https://github.com/nats-io/natscli)

---

## HACHI Documentation

### What is HACHI?  
HACHI is a connective technology that powers a distributed mesh. A connective technology is responsible for addressing, discovery and exchanging of messages and commands that drive the common patterns in distributed systems; asking and answering questions, aka services/microservices, and making and processing statements, or stream processing.

### Goals
>  Hachi = (all agents - controller & agents)

- Hachi must be easy to configure and operate and be observable.

- Hachi must support multiple use cases.

- Hachi must self-heal and always be available.

- Hachi must expose an API for interaction and configuration.

- Hachi must adapt to its environment (interpolation, discovery, resolving, invocation and execution) 

- Hachi must allow messages to flow in a bidirectional manner as desired

- Hachi must allow sync and async behaviour

- Hachi must display payload agnostic behavior

- Encryption at rest of the messages being stored.

- Streams can capture more than one subject

- Replay policies (needed here?)

- Stream replication factor


Replicas=1 - Cannot operate during an outage of the server servicing the stream. Highly performant.

Replicas=3 - Can tolerate loss of one server servicing the stream. An ideal balance between risk and performance.

Replicas=5 - Can tolerate simultaneous loss of two servers servicing the stream. Mitigates risk at the expense of performance.

need to test pull

ordered push consumer or pull
any need for replay?
ReplayInstant ? ReplayOriginal?

message tracking 

## DeliverNew

DeliverLastPerSubject
When first consuming messages, start with the latest one for each filtered subject currently in the stream.


AckExplicit or AckALL

MaxDeliver - Some messages may cause your applications to crash and cause a never ending loop forever poisoning your system. The MaxDeliver setting allow you to set an upper bound to how many times a message may be delivered. 



1. values/envars - deployment time interpolation
{{.interpolated_key}}

2. api invokation value - RPC execution time
{{.api.interpolated_key}}

   
for instance; 
    /api/${api_version}/method
    on system start apoi_version value will be interpolated into the base url

{{.value}} = these values will always be interpolated on API invokation from dynamic value sent to the url

Helper methods on Server

1. server-url/ListAllRoutes
   
2. controller = root agent - controller
3. agent = leaf agent - edge controller

### controller Configuration 

```hcl
version = 1

agent "Relix" {

  http = "0.0.0.0:8080"

  api {
    version = 1
    enabled = true
    allow_list = true //list all bindings

    auth {
      enabled = true
      provider = "{{.provider_addr}}"
      token_prefix = "{{.token_prefix}}"
    }
  }

  controller {
    enabled = true
    identifiers = ["controller.internal"]
  }

  agent {
    //It is recommended to keep the maximum number of tokens in your subjects to a reasonable value of 16
    enabled = false
    identifiers = []
  }

  storage {
    data_dir  = ""
  }

  kv_db {
    //https://github.com/dgraph-io/badger
  }

  stream {
    //avoid being over flooded/attacked by rogue dispatcher

    circuit_breaker {
      enabled = true
      max_requests = 100  //uint32
      interval = 1        //time.Duration in seconds
      timeout  = 3000     //time.Duration in seconds
    }

    deduping {
      enabled = true
      strategy = "default"
    }
  }

  axon {
    addr = "192.168.10.2"
    port = 4222
  }


  #Drivers support interpreting node attributes and runtime environment
  bindings {
    route "speak" {
      verb = "POST"
      routing = "cns.brain.{{.route::lobe}}.{{.route::region}}"
      local = "/cns/brain/:lobe/region/:region"
      remote = "{{.remote::audio_device_addr}}/{{.local::audio_quality}}/sonant"
      headers = [
        "hachi-relay: {{.remote::relay_service_addr}}",
        "hachi-token: {{.local::static_token}}"
      ]
    }
  }
}
```
### agent Configuration
```hcl
version = 1

agent "agent" {

  http = "0.0.0.0:8080"

  api {

    version = 1
    enabled = true
    allow_list = true //list all bindings

    auth {
      enabled = true
      provider = "{{.provider_addr}}"
      token_prefix = "{{.token_prefix}}"
    }
  }

  controller {
    enabled = false
  }

  agent {
    //It is recommended to keep the maximum number of tokens in your subjects to a reasonable value of 16
    identifiers = ["{{.local::identifier}}", "{{.local::region}}", "{{.local::functionality}}"]
    enabled = true
  }

  storage {
    data_dir  = ""
  }

  kv_db {
    //https://github.com/dgraph-io/badger
  }

  stream {
    //avoid being over flooded/attacked by rogue dispatcher

    circuit_breaker {
      enabled = true
      max_requests = 100  //uint32
      interval = 1        //time.Duration in seconds
      timeout  = 3000     //time.Duration in seconds
    }

    deduping {
      enabled = true
      strategy = "default"
    }
  }

  axon {
    addr = "192.168.10.2"
    port = 4222
  }


  #Drivers support interpreting node attributes and runtime environment
  bindings {
    route "speak" {
      verb = "POST"
      routing = "cns.brain.{{.route::lobe}}.{{.route::region}}"
      local = "/cns/brain/:lobe/region/:region"
      remote = "{{.remote::audio_device_addr}}/{{.local::audio_quality}}/sonant"
      headers = [
        "hachi-relay: {{.remote::relay_service_addr}}",
        "hachi-token: {{.local::static_token}}"
      ]
    }
  }
}
```
