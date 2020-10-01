# UI Authorization

By default, there is no authorization (anonymous user will be used).

If list of `auth.users` is not empty, all authorized users will be allowed.

## OAuth2

**This is mostly recommended way**

Defined in the section: `auth.oauth2`

* `title` - text that will be used for login button
* `secret` - OAuth2 client secret
* `key` - OAuth2 client ID
* `callback_url` - redirect URL, must point to your sever plus `/ui/auth/oauth2/callback`
* `auth_url` - authenticate URL, different for each provider
* `token_url` - issue token URL, different for each provider
* `profile_url` - URL that should return user JSON profile on GET request with authorization by token
* `login_field` - filed name (should be string) in profile that identifies user (ex: `login`, `username` or `email`)
* `scopes` (optional) - list of OAuth2 scopes


Gitea example:

```yaml
auth:
  oauth2:
    title: Gitea
    secret: "oauth secret"
    key: "oauth key"
    callback_url: "https://YOUR-SERVER/ui/auth/oauth2/callback"
    auth_url: "https://gitea-server/login/oauth/authorize"
    token_url: "https://gitea-server/login/oauth/access_token"
    profile_url: "https://gitea-server/api/v1/user"
    login_field: "login"
    scopes:
      - nano-run
  users:
    - "reddec"
```