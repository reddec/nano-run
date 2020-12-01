# Unit

* If work dir not defined - temporary directory will be created and removed after execution for each request automatically.


Schema:

* `command` (required, string) -  command to execute (will be executed in a shell)
* `user` (optional, string) - custom user as process owner (only `bin` mode and only for linux), usually requires root privileges
* `interval` (optional, interval) - interval between attempts
* `timeout` (optional, interval) - maximum execution timeout (enabled only for bin mode and only if positive)
* `graceful_timeout` (optional, interval) - maximum execution timeout after which SIGINT will be sent (enabled only for bin mode and only if positive).
Ie: how long to let command react on SIGTERM signal.
* `shell` (optional, string) - shell to execute command in bin mode (default - /bin/sh)
* `environment` (optional, map string=>string) - custom environment for executable (in addition to system)
* `max_request` (optional, integer) - maximum HTTP body size (enabled if positive)
* `attempts` (optional, integer) - maximum number of attempts
* `workers` (optional, integer) - concurrency level - number of parallel requests
* `mode` (optional, string) - execution mode: `bin` or `cgi` 
* `workdir` (optional, string) - working directory for the worker. if empty - temporary one will be generated automatically.
* `authorization` (optional, [Authorization](authorization.md)) - request authorization
* `cron` (optional,[Cron](cron.md)) - scheduled requests
* `private` (optional, bool) - do not expose over API, could be used for cron-only jobs