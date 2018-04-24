# init-consul-acl
This initialises Hashicorp Consul cluster with the following ACLs:
 * Agent ACL Token
 * Vault ACL Token

 For more information about Consul ACLs https://www.consul.io/docs/guides/acl.html#configuring-acls and for Vault ACLs can be found at https://www.vaultproject.io/docs/configuration/storage/consul.html.

## Building the project
This project uses dep so it must be on your path to begin with.
```
dep ensure
go build ./...
docker build -t init-consul-acl .
```

## Running the container
The container can be run using the following command and passing environment variables where required.
```
docker run --restart=no -e AGENT_ACL_TOKEN=uuidgen -e VAULT_ACL_TOKEN=uuidgen -e READ_ONLY_ACL_TOKEN=uuidgen bhavikk/init-consul-acl:latest
```
