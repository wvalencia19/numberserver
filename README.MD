# Number Server

## Problem

Using any of the following programming language (taking performance into consideration):
and restricts input to at most 5 concurrent clients. Clients will connect to the Application and
write one or more numbers of 9 digit numbers, each number followed by a server-native newline
sequence, and then close the connection. The Application must write a de- duplicated list of
these numbers to a log file in no particular order.


## Requirements
1. The Application must accept input from at most 5 concurrent clients on TCP/IP port4000.
2. Input lines presented to the Application via its socket must either be composed of exactly
nine decimal digits (e.g.: 314159265 or 007007009) immediately followed by a
server-native newline sequence; or a termination sequence as detailed in #9, below.
3. Numbers presented to the Application must include leading zeros as necessary to
ensure they are each 9 decimal digits.
4. The log file, to be named "numbers.log”, must be created anew and/or cleared when the
Application starts.
5. Only numbers may be written to the log file. Each number must be followed by a
server-native newline sequence.
6. No duplicate numbers may be written to the log file.
7. Any data that does not conform to a valid line of input should be discarded and the client
connection terminated immediately and without comment.
8. Every 10 seconds, the Application must print a report to standard output:
a. The difference since the last report of the count of new unique numbers that have
been received.
b. The difference since the last report of the count of new duplicate numbers that
have been received.
c. The total number of unique numbers received for this run of the Application.
d. Example text for #8: Received 50 unique numbers, 2 duplicates. Unique total:
567231

9. If any connected client writes a single line with only the word "terminate" followed by a
server-native newline sequence, the Application must disconnect all clients and perform
a clean shutdown as quickly as possible.
10. Clearly state all of the assumptions you made in completing the Application.

## Notes
* You may write tests at your own discretion. Tests are useful to ensure your Application
passes Primary Consideration A.
* You may use common libraries in your project such as Apache Commons and Google
Guava, particularly if their use helps improve Application simplicity and readability.
However the use of large frameworks, such as Akka, is prohibited.
* Your Application may not for any part of its operation use or require the use of external
systems, for example Apache Kafka or Redis.
* At your discretion, leading zeroes present in the input may be stripped—or not
used—when writing output to the log or console.
* Robust implementations of the Application typically handle more than 2M numbers per
10-second reporting period on a modern MacBook Pro laptop (e.g.: 16 GiB of RAM and
* 2.5 GHz Intel i7 processor).
* To test if your application is working as expected, you can try to telnet to it through the
port 4000 by executing:

```
> telnet localhost 4000
```
And manually type in the numbers sequentially followed by a newline (enter).

## Solution

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


