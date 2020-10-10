## BIN

Binary modes just executes any script in shell (`/bin/sh` by default). You can override shell per-unit
by `shell: /path/to/shell` configuration param.

Linux's hosts should not worry about "leaking" processes (when process
creates number of background subprocesses) - it wil be cleaned automatically.

**For non-Linux users**

To handle a graceful timeout, child should be able to forward signal: basically, use `exec` before last command.

Danger (but will work), signals may not be handled by foo

```yaml
command: "V=1 RAIL=2 foo bar -c -y -z"
```

Good

```yaml
command: "V=1 RAIL=2 exec foo bar -c -y -z"
```