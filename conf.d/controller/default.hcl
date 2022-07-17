version = 1

dna "controller" {

  controller {
    enabled = true

    identifiers {
      core = "controller"
      descriptors = []
    }
    invocation_timeout = -1
  }

  agent {
    //It is recommended to keep the maximum number of tokens in your subjects to a reasonable value of 16
    enabled = false
  }

  nats {
    addr = "0.0.0.0"
    port = 4222
  }

  http {
    addr = "0.0.0.0"
    port = 3007
  }

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

  storage {
    data_dir  = ""
  }

  kv_db {
    //https://github.com/dgraph-io/badger
  }

  hrl {
    crypto {
      provider = "AES"
      encrypt_endpoint = ""
      decrypt_endpoint = ""
    }
  }

  stream {
    //avoid being over-flooded/attacked by rogue dispatchers
    circuit_breaker {
      enabled = true
      max_requests = 100  //uint32
      interval = 1        //time.Duration in seconds
      timeout  = 3000     //time.Duration in seconds
    }

    //nats
    deduping {
      enabled = true
      strategy = "default"
    }
  }

  #Drivers support interpreting node attributes and runtime environment
  tracts {

    stream "simple_async_result" {
      async = true
      verb = "POST"
      selectors = ["stations.>.{{.route::type}}", "agent.internal.>"]
      local = "/metrics/type/:type/filter/:filter"
      remote {
        http {
          url = "metrics_service.remote_server:8080/metrics/{{.route::type}}/{{.route::filter}}"
        }
      }
      headers = {
        "hachi-relay-x" = ["{{.remote::relay_service_addr}}", "{{.route::type}}"]
        "hachi-token" = ["{{.local::static_token}}"]
      }
    }
    //event from url
    stream "trigger_webhook" {
      async = true
      verb = "POST"
      selectors = ["agents.selector.{{.route::selector}}"]
      local = "/events/:event/selector/:selector"
      remote {
        webhook {
          event = "{{.route::event}}"
        }
      }
    }

    //static event from configuration
    stream "trigger_webhook_const_event" {
      async = true
      verb = "POST"
      selectors = ["agents.selector.{{.route::selector}}"]
      local = "/selector/:selector"
      remote {
        webhook {
          event = "EVENT.DATA.SOME_DATA_CHANGE"
        }
      }
    }

    stream "gossip" {
      async = true
      verb = "POST"
      selectors = ["neurostream.controller.to.agents"]
      local = "/test"
      remote {
        http {
          url = "{{.remote::service_addr}}/"
        }
      }
      headers = {
        "hachi-token" = ["{{.local::static_token}}"]
      }
    }

    stream "speak" {
      async = true
      verb = "POST"
      selectors = ["cns.brain.{{.route::lobe}}.{{.route::region}}","ORDER.cns","neurostream.controller.to.agents"]
      local = "/transactions/upload/:transtype"
      remote {
        http {
          url = "{{.storix_addr}}/p97/gift"
        }
      }
      headers = {
        "hachi-relay-x" = ["{{.remote::relay_service_addr}}", "{{.local::static_token}}"]
        "hachi-token" = ["{{.local::static_token}}"]
      }
    }

    stream "shout" {
      async = true
      verb = "POST"
      selectors = ["cns.brain.{{.route::lobe}}.{{.route::region}}","ORDER.cns","neurostream.controller.to.agents"]
      local = "/cns/brain/:lobe/region/:region"
      remote {
        http {
          url = "{{.remote::audio_device_addr}}/{{.local::audio_quality}}/sonant"
        }
      }
      headers = {
        "hachi-relay-x" = ["{{.remote::relay_service_addr}}", "{{.local::static_token}}"]
        "hachi-token" = ["{{.local::static_token}}"]
      }
    }
  }
}