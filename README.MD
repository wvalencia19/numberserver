# Number Server

## Considerations

To meet the requirements of the test I have made the following decisions:

* Using go routines and channels to improve the performance of the service.
* Using the [concurrency pipelines pattern](https://go.dev/blog/pipelines), which in this case is composed
  by two stages connected by channels, which will allow concurrent processing (writing report and writing log) of the messages received in the socket.
* Sharing the context with the `gracefullyShutdown` function to listen for the `SIGINT` and `SIGTERM` signals and gracefully terminate the process by canceling the context.
* Creation of the numberSet type that simulates a `Set` data type, this with the aim of save in memory the unique values having in account the time complexity; O(1) for adding, removing, and verifying membership of elements.

Other aspects of the design of the solution:

* The logs are categorized by level, it is possible to select the logging level of the service with the LogLevel environment variable.
* It is possible to configure different aspects of the application through environment variables.
* Unit tests are included.
* Including Make as interface to execute commands to compile, test, check(linter) the program.
* Docker configuration.

## Environment Variables

* NUMBERSERVER_TERMINATIONWORD: termination word.
* NUMBERSERVER_MESSAGELENREQUIRED: limit of the digit numbers.
* NUMBERSERVER_LOGLEVEL: log level
* NUMBERSERVER_REPORTPERIODICITY: periodicity of print report to standard output.
* NUMBERSERVER_LOGNAME:name of the log file.
* NUMBERSERVER_SERVER_MAXCONN: max TCP server concurrent connections.
* NUMBERSERVER_SERVER_HOST: TCP server host.
* NUMBERSERVER_SERVER_PORT: TCP server port.

## Execute the program

* Locally

```
go run -race cmd/numberserver/main.go
```

* Docker

```
make compile
docker build --no-cache -t numberserver .
docker run -i -t -p 4000:4000 numberserver
```

## Run tests

### Unit tests
```
make test
```

### Manual test
I have included a client which reads from a file with 2.4 million of rows to test the performance of the service.
```
go run cmd/client/main.go
```
## Run checks
```
make check
```


