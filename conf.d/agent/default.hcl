version = 1.0

dna "agent" {

  controller {
    enabled = false
    invocation_timeout = -1
  }

  agent {
    //It is recommended to keep the maximum number of tokens in your subjects to a reasonable value of 16
    enabled = true
    identifiers {
      core = ""
      descriptors = ["stations", "north", "galil" , "large", "happy", "mobile"]
    }
  }

  nats {
    addr = "0.0.0.0"
    port = 4222
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
    stream "register" {
      async = true
      verb = "POST"
      selectors = ["neurostream.controller.to.agents"]
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