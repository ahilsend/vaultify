# vaultify

[![Build Status](https://travis-ci.org/ahilsend/vaultify.svg?branch=master)](https://travis-ci.org/ahilsend/vaultify)
[![](https://img.shields.io/badge/docker%20build-automated-blue.svg)](https://hub.docker.com/r/ahilsend/vaultify "docker build - automated")

Vaultify templates file from vault secrets and auto renews leases

## Running vaultify

`vaultify` has three commands, `template`, `renew-leases`, and `run`

### Template

The `template` command reads a template, renders the vault secrets into it, and stores the result in a file. In addition it also stores the secret lease information in a secrets file to be able to renew the leases.

template.yaml example:
```yaml
credentials:
    <{- $admin := vault "database/creds/maindb-admin" }>
    username: <{ $admin.Data.username | quote }>
    password: <{ $admin.Data.password | quote }>
```

Running `vaultify template`:
```bash
vaultify template --vault https://vault.vault:8200 \
                  --role maindb-admin \
                  --template-file template.yaml \
                  --output-file /app/config.yaml \
                  --secrets-output-file /app/secrets.json \
                  -vv
```

### Renew-leases

The `renew-leases` command renews leases that for created by `template` command and stored in a secrets file.

Running `vaultify renew-leases`:
```bash
vaultify renew-leases --vault https://vault.vault:8200 \
                      --role maindb-admin \
                      --secrets-output-file /app/secrets.json \
                      -vv
```


### Run

Running vaultify and continuously renew leases:

```bash
vaultify run --vault https://vault.vault:8200 \
             --role maindb-admin \
             --template-file template.yaml \
             --output-file /app/config.yaml \
             -vv
```

Note that running only this might not work for all work loads. If you run your application in kubernetes and your configuration needs to be rendered before the application starts, you should run the `template` command in a initContainer and the `renew-leases` command in a side-car.
