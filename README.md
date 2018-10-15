# vaultify

[![Build Status](https://travis-ci.org/ahilsend/vaultify.svg?branch=master)](https://travis-ci.org/ahilsend/vaultify)
[![](https://img.shields.io/badge/docker%20build-automated-blue.svg)](https://hub.docker.com/r/ahilsend/vaultify "docker build - automated")

Vaultify templates file from vault secrets and auto renews leases

## Exammple

template.yaml
```yaml
credentials:
    <{- $admin := vault "database/creds/maindb-admin" }>
    username: <{ $admin.Data.username | quote }>
    password: <{ $admin.Data.password | quote }>
```

Running vaultify and continuously renew leases
```bash
vaultify run --vault https://vault.vault:8200 \
             --role maindb-admin \
             --template-file template.yaml \
             --output-file /app/config.yaml \
             -vv
```
