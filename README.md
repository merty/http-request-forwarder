# http-request-forwarder

[![Go Report Card](https://goreportcard.com/badge/github.com/merty/http-request-forwarder)](https://goreportcard.com/report/github.com/merty/http-request-forwarder)

A reverse HTTP proxy to duplicate incoming HTTP traffic to given hosts.

Useful for replicating production traffic to run tests on the development environment or simply to transmit an HTTP request to multiple servers using a single HTTP request.

## Usage

In order to forward the incoming traffic on port `8080` to ports `8081`, `8082` and `8083` on the same host:

```bash
$ ./http-request-forwarder -l 8080 -h "localhost:8081,localhost:8082,localhost:8083"
```

You may also provide a timeout duration in milliseconds using the `-t` flag, which defaults to 3 seconds:

```bash
$ ./http-request-forwarder -l 8080 -h "localhost:8081,localhost:8082,localhost:8083" -t 1000
```

## Notes

Requests always return the status code `200 OK` as different target hosts may return different status codes. This behavior may change in later versions.

## Changelog

**0.1.0**

* Initial release.

## Author

Mert Yazıcıoğlu - [Website](https://mertyazicioglu.com) &middot; [GitHub](https://github.com/merty) &middot; [Twitter](https://twitter.com/_mert)

## License

This project is released under the MIT License. See the `LICENSE` file for details.