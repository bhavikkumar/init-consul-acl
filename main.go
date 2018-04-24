package main

import (
  "os"
  "time"
  "github.com/hashicorp/consul/api"
  "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
  // This comes from https://www.consul.io/docs/guides/acl.html#configuring-acls
  // Defines the environment variable name which sets the ID for the agent acl
  // token.
  AgentTokenID = "AGENT_ACL_TOKEN"

  // Defines the environment variable name which sets the ID for the vault acl
  // token.
  VaultTokenID = "VAULT_ACL_TOKEN"

  // Defines the environment variable name which sets the ID for the read only
  // acl token.
  ReadOnlyID = "READ_ONLY_ACL_TOKEN"
)

func main()  {
  zerolog.SetGlobalLevel(zerolog.InfoLevel)

  client, err := api.NewClient(api.DefaultConfig())
  if err != nil {
      log.Fatal().Err(err).Msg("Unable to create new client")
  }

  // While the consul cluster isn't running, we wait
  for i := 0; i < 5; i++ {
    _, _, err := client.Health().Service("consul", "", true, nil)
    if err == nil {
      break
    }
    log.Info().Msg("Consul cluster still not running, waiting 30 seconds to try again")
    time.Sleep(30 * time.Second)
  }

  acl := client.ACL()
  createAgentToken(acl)
  createVaultToken(acl)
  createReadOnlyToken(acl)
  log.Info().Msg("Finished creating ACL Tokens")
}

func createAgentToken(acl *api.ACL) {
  if tokenId := os.Getenv(AgentTokenID); tokenId != "" {
    aclEntry := api.ACLEntry{
      ID: tokenId,
      Name: "Agent Token",
      Type: api.ACLClientType,
      Rules: `node "" { policy = "write" } service "" { policy = "read" }`,
    }
    id, _, err := acl.Create(&aclEntry, nil)

    if err != nil {
      log.Fatal().Err(err).Msg("Unable to create new Agent ACL")
    }

    if id == "" {
      log.Fatal().Msg("Agent ACL ID was empty after create")
    }
    log.Info().Msg("Agent Token Successfully Created")
    return
  }
  log.Fatal().Msg("Agent ACL Token is mandatory")
}

func createVaultToken(acl *api.ACL) {
  if tokenId := os.Getenv(VaultTokenID); tokenId != "" {
    aclEntry := api.ACLEntry{
      ID: tokenId,
      Name: "Vault Token",
      Type: api.ACLClientType,
      Rules: `key "vault/" { policy = "write" } node "" { policy = "write" } service "vault" { policy = "write" } agent "" { policy = "write" } session "" { policy = "write" }`,
    }
    id, _, err := acl.Create(&aclEntry, nil)

    if err != nil {
      log.Error().Err(err).Msg("Unable to create new vault ACL")
    }

    if id == "" {
      log.Error().Msg("Vault ACL ID was empty after create")
    }
    log.Info().Msg("Vault Token Successfully Created")
    return
  }
  log.Info().Msg("Vault token not provided.")
}

func createReadOnlyToken(acl *api.ACL) {
  if tokenId := os.Getenv(ReadOnlyID); tokenId != "" {
    aclEntry := api.ACLEntry{
      ID: tokenId,
      Name: "Read Only Token",
      Type: api.ACLClientType,
      Rules: `key "" { policy = "read" } node "" { policy = "read" } service "" { policy = "read" } agent "" { policy = "read" } session "" { policy = "read" }`,
    }
    id, _, err := acl.Create(&aclEntry, nil)

    if err != nil {
      log.Error().Err(err).Msg("Unable to create new read only ACL")
    }

    if id == "" {
      log.Error().Msg("Read Only ACL ID was empty after create")
    }
    log.Info().Msg("Read Only Token Successfully Created")
    return
  }
  log.Info().Msg("Read only token not provided.")
}
