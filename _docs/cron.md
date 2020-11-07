# Cron


Cron-like jobs allowed as part of Unit definition [thanks to robfig/cron](https://godoc.org/github.com/robfig/cron#hdr-Usage).

**Example definition:**

```yaml
# ... unit definition above ...
cron:
    # every hour on the half hour
    - spec: 30 * * * *
    # same as above but with name to detect in UI and logs
    - spec: 30 * * * *
      name: named schedule
    # each hour with custom payload and headers
    - spec: @hourly
      content: |
        hello world
      headers:
        X-Some-Header: test-header
    # each day with content from file
    - spec: @daily
      content_file: /path/to/content
```


Schema:

* `spec` (required, string) - cron tab specification for the time interval. See [online builder](https://crontab.guru/)
* `name` (optional, string) - name for entry to distinguish record in UI or in logs.
* `headers` (optional, map string=>string) - headers that will be used in simulated request.
* `content` (optional, string) - simulated request content.
* `content_file`  (optional, string) - simulated request content file. Has less priority than `content`.

Caveats:

* **Security:** cron job ignores authorizations defined on unit level.
* **Enqueuing:** cron job will be enqueued regardless of status of previous job.
* **Errors:** if cron job can't create request (ex: `content_file` not available) - it will print error to log and 
will try again later at the next schedule.