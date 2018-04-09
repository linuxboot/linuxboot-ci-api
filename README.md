# Linuxboot CI API

API to submit job to the Linuxboot CI platform.

# API

__NB.__
 > API is in work in progress state and can change in the future

## Submit a job

```
POST /v1/jobs
```

__Body__

```
{
    "repository": {
	    "url": "https://github.com/ggiamarchi/linuxboot-ci-test.git"
    }
}
```

or, to build a specific branch

```
{
    "repository": {
	    "url": "https://github.com/ggiamarchi/linuxboot-ci-test.git",
        "branch": "master"
    }
}
```

__Response__

```
HTTP/1.1 200 ACCEPTED
Content-Type: application/json; charset=utf-8
...

{
    "id":199,
    "repository":{
        "url":"https://github.com/ggiamarchi/linuxboot-ci-test.git",
        "branch":null
    },
    "submitDate":"2018-04-09T17:01:30.942378478+02:00",
    "status":"PENDING"
}
```

## Get job status

```
GET /v1/jobs/:jobId
```

__Response__

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
...

{
    "id":199,
    "repository":{
        "url":"https://github.com/ggiamarchi/linuxboot-ci-test.git",
        "branch":null
    },
    "submitDate":"2018-04-09T17:01:30.942378478+02:00",
    "status":"RUNNING"
}
```

## Get job logs

```
GET /v1/jobs/:jobId/logs[?raw=true]
```

If `raw=true` query parameter is set, API response will be plain text instead of JSON.

__Response__

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
...

{
    "log": "..."
}
```

or with `raw=true`

```
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
...

[2018-04-09 10:21:00] Setting up virtual machine configuration...
[2018-04-09 10:21:25] Running virtual machine...
[2018-04-09 10:21:28] Waiting for virtual machine network...
...
```
