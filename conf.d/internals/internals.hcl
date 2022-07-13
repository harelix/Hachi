dna "internals" {
  tracts {
    //static event from configuration
    stream "webhook_registration_endpoint" {
      async = true
      verb = "POST"
      subject = [""]
      local = "/webhooks/:client"
      remote {
        webhook {
          event = "__internal__.webhook.register#endpoint"
        }
      }
    }

    /*
    stream "edge_devices_registration" {
      async   = false
      subject = ["__internal__.selectors.ed#ed"]
      verb    = "POST"
      local   = "/hrl/edge_devices_registration"
      remote {
        internal {
          directive = "encrypt"
        }
      }
    }
    */

    stream "encrypt" {
      async   = false
      subject = ["__internal__.crypto.encrypt#selector"]
      verb    = "POST"
      local   = "/hrl/encrypt"
      remote {
        internal {
          directive = "encrypt"
        }
      }
    }

    stream "decrypt" {
      async   = false
      subject = ["__internal__.crypto.decrypt#selector"]
      verb    = "POST"
      local   = "/hrl/decrypt"
      remote {
        internal {
          directive = "decrypt"
        }
      }
    }
  }
}
