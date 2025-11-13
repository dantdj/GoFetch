# GoFetch
Basic terminal-based download manager written in Go to experiment with resumable downloads.

Includes a basic system to indicate progress, as well as allowing resuming of failed downloads.

## How does it work?

The HTTP spec has a few interesting headers for supporting this:
* [Accept-Ranges](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Accept-Ranges) - this header is returned on the response when trying to obtain a file. This indicates that the server will accept a request for part of a file
* [Content-Range](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Range) - this header is returned on a `206 Partial Content` response to indicate which bytes the response contains
* [Range](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Range) - this header is included in requests to a server to download a file, and specifies the range of bytes that is being requested

This application largely makes use of the last two headers.