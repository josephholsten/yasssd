# yasssd - Yet Another Simple Storage Service Daemon

`yasssd` is a service to provide a basic storage service. It's almost certainly inapropriate to your needs, see the Alternatives section for something more appropriate.

## The Simplest Thing That Could Possibly Work

```
$ ./script/bootstrap # install dependencies
$ ./script/build     # build executable
$ ./yasssd &         # start service

$ curl -X POST                   localhost:8080/register  -d '{"username":"bib","password":"secretiv"}'  # create an account
$ curl -X POST                   localhost:8080/login     -d '{"username":"bib","password":"secretiv"}'  # request session token
{"token":1}
$ curl -X PUT    -H "X-Token: 1" localhost:8080/files/foo -d "foo-file-contents"                         # create
$ curl -X GET    -H "X-Token: 1" localhost:8080/files/foo                                                # read
foo-file-contents
$ curl -X GET    -H "X-Token: 1" localhost:8080/files                                                    # list
[
  "foo"
]
$ curl -X PUT    -H "X-Token: 1" localhost:8080/files/foo -d "foo-file-contents"                         # update
$ curl -X DELETE -H "X-Token: 1" localhost:8080/files/foo                                                # delete
```

## Alternatives

- Amazon S3: Possibly the most hosted popular storage service in the world. Originally had very few features. These days has an incredible number of advanced extensions while still supporting the original API. If you intend you run your own service: you will have a very hard time achieving better reliability, and will need to spend quite a bit up-front to beat S3's pay-as-you-go rate. Its latency and consistency properties are distinctive, usually if you're trying to avoid S3 it's for this reason.
- Minio: If you're here because you thought you would find an S3 compatible service written in go, you almost certainly should be here instead. 
- For other options, see http://www.s3-client.com/s3-compatible-storage-solutions.html
