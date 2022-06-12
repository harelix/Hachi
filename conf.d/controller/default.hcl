version = 1

dna "controller" {

  controller {
    enabled = true
    identifiers = ["controller.internal", "cns.brain.frontal.broca"]
    invocation_timeout = -1
  }

  agent {
    //It is recommended to keep the maximum number of tokens in your subjects to a reasonable value of 16
    enabled = false
    identifiers = []
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

  nats {
    addr = "0.0.0.0"
    port = 4222
  }

  #Drivers support interpreting node attributes and runtime environment
  tracts {
    #Don't remove these endpoints as they provide selectors encryption capability to Hachi
    stream "encrypt" {
      async = false
      subject = ["__internal__.crypto.encrypt#selector"]
      verb = "POST"
      local = "/hrl/encrypt"
    }
    stream "decrypt" {
      async = false
      subject = ["__internal__.crypto.decrypt#selector"]
      verb = "POST"
      local = "/hrl/decrypt"
    }
    #Don't remove these endpoints as they provide selectors encryption capability to Hachi

    stream "gossip" {
      async = true
      verb = "POST"
      subject = ["neurostream.controller.to.agents"]
      local = "/test"
      remote = "{{.remote::service_addr}}/"
      headers = {
        "hachi-token" = ["{{.local::static_token}}"]
      }
    }

    stream "speak" {
      async = true
      verb = "POST"
      subject = ["cns.brain.{{.route::lobe}}.{{.route::region}}","ORDER.cns","neurostream.controller.to.agents"]
      local = "/api/v2/transactions/upload/:transtype"
      remote = "{{.storix_addr}}/p97/gift"
      headers = {
        "hachi-relay-x" = ["{{.remote::relay_service_addr}}", "{{.local::static_token}}"]
        "hachi-token" = ["{{.local::static_token}}"]
      }
    }

    stream "speak" {
      async = true
      verb = "POST"
      subject = ["cns.brain.{{.route::lobe}}.{{.route::region}}","ORDER.cns","neurostream.controller.to.agents"]
      local = "/cns/brain/:lobe/region/:region"
      remote = "{{.remote::audio_device_addr}}/{{.local::audio_quality}}/sonant"
      headers = {
        "hachi-relay-x" = ["{{.remote::relay_service_addr}}", "{{.local::static_token}}"]
        "hachi-token" = ["{{.local::static_token}}"]
      }
    }
  }
}