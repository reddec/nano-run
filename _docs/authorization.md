# Authorization

By-default - authorization disabled. Multiple policies allowed. 
To allow request at least one policy should be passed.
Each authorization policy can enabled by `enable: yes` param.

Section in `server.yaml`: `authorization`

## JWT

*section: `authorization.jwt`*

[Overview](https://jwt.io/)

HMAC 256 signature validation against secret key

Configurable parameters:

* `header` (optional, string, default: `Authorization`) - header that contains JWT
* `secret` (required, string) - secret key to validate signature

Example minimal unit config

```yaml
command: 'echo hello world'
authorization:
    jwt:
      enable: yes
      secret: '$eCrEtKey'
```

## Query token

*section: `authorization.query_token`*

Plain token in a query string. Will be matched against list of allowed tokens.

For example, client can invoke endpoint by addition token query: `http://example.com/app/?token=deadbeaf`

Configurable parameters:

* `param` (optional, string, default: `token`) - query param where token should be placed
* `tokens` (required, []string) - list of allowed tokens

Example minimal unit config with 3 tokens

```yaml
command: 'echo hello world'
authorization:
    query_token:
      enable: yes
      tokens:
        - my-token-1
        - his-token-2
        - deadbeaf
```

## Header token

*section: `authorization.header_token`*

Plain token in a header. Will be matched against list of allowed tokens.

For example, client can invoke endpoint by curl: 

    curl -H 'X-Api-Token: deadbeaf' http://example.com/app/

Configurable parameters:

* `header` (optional, string, default: `X-Api-Token`) - header name where token should be placed
* `tokens` (required, []string) - list of allowed tokens

Example minimal unit config with 3 tokens

```yaml
command: 'echo hello world'
authorization:
    header_token:
      enable: yes
      tokens:
        - my-token-1
        - his-token-2
        - deadbeaf
```

## Basic

*section: `authorization.basic`*

Basic authentication. [Overview](https://en.wikipedia.org/wiki/Basic_access_authentication)

For example, client can invoke endpoint by curl: 

    curl -u 'alice:admin' http://example.com/app/

To [calculate](https://unix.stackexchange.com/a/419855) hash you may use `htpasswd` (Debian/Ubuntu: `sudo apt install apache2-utils`)

    htpasswd -bnBC 10 "" password | tr -d ':'

where `passsword` is a desired password for the user.

Configurable parameters:

* `users` (string->string, required) - map of users and their hashed password by bcrypt

Example minimal config:

```yaml
command: 'echo hello world'
authorization:
    basic:
      enable: yes
      users:
        alice: '$2y$10$cUe3n8NHaxee.AaGzT8wF.nirPnjv5YLEQGTsLiiMiUAknM2aF2FS'
        bob: '$2y$10$iSczi.MlKTrMv3h0Zf.GDeW1NS6ZWxBgtj4ytrKKDrR2s2wIxq5Qa'
```
