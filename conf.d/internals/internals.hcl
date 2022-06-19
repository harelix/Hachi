dna "internals" {
  tracts {
    stream "encrypt" {
      async   = false
      subject = ["__internal__.crypto.encrypt#selector"]
      verb    = "POST"
      local   = "/hrl/encrypt"
    }

    stream "decrypt" {
      async   = false
      subject = ["__internal__.crypto.decrypt#selector"]
      verb    = "POST"
      local   = "/hrl/decrypt"
    }
  }
}
