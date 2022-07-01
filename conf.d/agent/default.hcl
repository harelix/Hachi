version = 1

dna "agent" {

  controller {
    enabled = false
  }

  agent {
    //what mode are we in (controller or agent)
    enabled = true
    //It is recommended to keep the maximum number of tokens in your subjects to a reasonable value of 16
    identifiers = ["{{.local::identifier}}", "{{.local::region}}", "{{.local::functionality}}"]
    //how long before an invocation request must finish executing on the target agent - in nillisecond
    invocation_timeout = 10000
  }

  http {
    addr = "0.0.0.0"
    port = 3008
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

  nats {
    addr = "0.0.0.0"
    port = 4222
  }
  
  #drivers support interpreting node attributes and runtime environment
  tracts {

    stream "self_test" {
      async = true
      verb = "GET"
      subject = ["neurostream.controller.to.agents"]
      payload = <<EOF
                {"name" : "testing", "addr" : "{{.remote::service_addr}}"}
                EOF
      local = "/selft"
      remote = "{{.remote::service_addr}}/"
      headers = {
        "hachi-token" = ["{{.local::static_token}}"]
      }
    }

    stream "gossip" {
      async = true
      verb = "GET"
      subject = ["neurostream.agent.to.controller"]
      local = "/test"
      remote = "{{.remote::service_addr}}/"
      headers = {
        "hachi-token" = ["{{.local::static_token}}"]
      }
    }

    stream "speak" {
      verb = "POST"
      subject = ["cns.brain.{{.route::lobe}}.{{.route::region}}","neurostream.agent.to.controller","ORDER.cns.{{.route::region}}"]
      local = "/cns/brain/:lobe/region/:region"
      remote = "{{.remote::audio_device_addr}}/{{.local::audio_quality}}/sonant"
      headers = {
        "hachi-relay" = ["{{.remote::relay_service_addr}}"]
        "hachi-token" = ["{{.local::static_token}}"]
      }
    }
  }
}