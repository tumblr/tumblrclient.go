# Tumblr API Client Library

[![Build Status](https://travis-ci.org/tumblr/tumblrclient.go.svg?branch=master)](https://travis-ci.org/tumblr/tumblrclient.go) [![GoDoc](https://godoc.org/github.com/tumblr/tumblrclient.go?status.svg)](https://godoc.org/github.com/tumblr/tumblrclient.go)

This is a concrete implementation of the `ClientInterface` defined in the [Tumblr API](https://github.com/tumblr/tumblr.go) library.

This project utilizes an external OAuth1 library you can find at [github.com/dghubble/oauth1](https://github.com/dghubble/oauth1) 

Install by running `go get github.com/tumblr/tumblrclient.go`

You can instantiate a client by running

```
client := tumblrclient.NewClient(
    "CONSUMER KEY",
    "CONSUMER SECRET",
)
// or
client := tumblr_go.NewClientWithToken(
    "CONSUMER KEY",
    "CONSUMER SECRET",
    "USER TOKEN",
    "USER TOKEN SECRET",
)
```

From there you can use the convenience methods.

## Support/Questions
You can post a question in the [Google Group](https://groups.google.com/forum/#!forum/tumblr-api) or contact the Tumblr API Team at [api@tumblr.com](mailto:api@tumblr.com)

## License

Copyright 2016 Tumblr, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
