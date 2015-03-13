# urit

[![Build Status](https://travis-ci.org/nowk/urit.svg?branch=master)](https://travis-ci.org/nowk/urit)
[![GoDoc](https://godoc.org/gopkg.in/nowk/urit.v0?status.svg)](http://godoc.org/gopkg.in/nowk/urit.v0)

URI Templates in Go

Based on [RFC6570](http://tools.ietf.org/html/rfc6570), provides up-to op-level2 and op-level3 expansions.


## Install

    go get gopkg.in/nowk/urit.v0

## Usage

    var u urit.URI = "{/path,id}/comments"

    u = u.Expand(urit.Variables{
      "path": "posts",
      "id": "1234",
    })

    // -> /posts/123/comments

---

__Inspect__

`Inspect()` allows you to expand and rewrite any unexpanded variables back into the uri. This allows for chaining or continuous expanding of a uri through different stages.

    var u urit.URI = "{/path,id}/comments"

    u = u.Inspect(urit.Variables{
      "path": "posts",
    })

    // -> /posts{/id}/comments

    ... some time later ...

    u = u.Expand(urit.Variables{
      "id": "123",
    })

    // -> /posts/123/comments


`Inspect` is incomplete at this stage since partial expansion on some operators end up requiring using of expression syntax that is not consistent or found in the original proposal.


## License

MIT
