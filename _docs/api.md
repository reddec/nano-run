# API

Base url: `/{application}`


## Start process

|||
|------------|------|
| **Method** | POST |
| **Path**   | /    |

Saves and enqueue request. Return 303 See Other on success with headers.

* `X-Correlation-Id` - with ID of request
* `Location` - URL to status

Clients can use cUrl flag `-L` to automatically follow redirect

    curl -L -d '' 'http://127.0.0.1:8989/app/'

## Get status


|||
|------------|------------------|
| **Method** | GET              |
| **Path**   | /{correlationId} |

Returns JSON with full meta-data of request and following headers:

* `Content-Version` - number of attempts
* `Last-Modified` - latest of time of creation, time of last attempt or completion time
* `Location` - URL to the complete attempt
* `X-Status` - status of request processing: `complete` or `processing`
* `X-Last-Attempt` - id of last attempt
* `X-Last-Attempt-At` - time of last attempt
* `X-Correlation-Id` - with ID of request

Body example:

```json
{
  "created_at": "2020-09-10T17:11:33.598542177+08:00",
  "complete_at": "2020-09-10T17:11:33.616550544+08:00",
  "attempts": [
    {
      "code": 200,
      "headers": {},
      "id": "51748767-e89b-48a1-8b00-b9c1f0fdc9bb",
      "created_at": "2020-09-10T17:11:33.614030236+08:00"
    }
  ],
  "headers": {
    "Accept": [
      "*/*"
    ],
    "Content-Length": [
      "0"
    ],
    "Content-Type": [
      "application/x-www-form-urlencoded"
    ],
    "User-Agent": [
      "curl/7.68.0"
    ]
  },
  "uri": "/date/",
  "method": "POST",
  "complete": true
}
```

## Get complete

|||
|------------|---------------------------|
| **Method** | GET                       |
| **Path**   | /{correlationId}/complete |

Will redirect to the complete attempt or 404


## Force complete

|||
|------------|---------------------------|
| **Method** | DELET                     |
| **Path**   | /{correlationId}          |

Forcefully mark request as complete and stop processing (including re-queue)


## Get attempt

|||
|------------|--------------------------------------|
| **Method** | GET                                  |
| **Path**   | /{correlationId}/attempt/{attemptId} |

Get result of processing request for the defined attempt

Returns body, code and headers same as processor returned with additional headers:

* `X-Status` - status of request processing: `complete` or `processing`
* `X-Processed` - `true` to distinguish result



## Get request

|||
|------------|--------------------------|
| **Method** | GET                      |
| **Path**   | /{correlationId}/request |

Get request same as it was POSTed for the defined ID

Returns body and headers same as processor got from client with additional headers:

* `X-Method` - request method (currently always `POST`)
* `X-Request-Uri` - request URI
* `Last-Modified` - time of creation
