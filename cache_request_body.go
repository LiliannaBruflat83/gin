package gin

import (
	"bytes"
	"io"
)

// CacheDecodedBody returns a middleware that caches the request body.
func CacheDecodedBody() HandlerFunc {
	return func(c *Context) {
		if c.Request != nil && c.Request.Body != nil {
			if _, ok := c.Request.Body.(*repeatableReader); !ok {
				var body []byte
				if c.Request.ContentLength > 0 {
					body = make([]byte, 0, c.Request.ContentLength)
				}
				c.Request.Body = &repeatableReader{
					rc:   c.Request.Body,
					body: body,
					c:    c,
				}
			}
		}
		c.Next()
	}
}

// CacheRequestBody returns a middleware that caches the request body.
func CacheRequestBody() HandlerFunc {
	return CacheDecodedBody()
}

type repeatableReader struct {
	rc       io.ReadCloser
	body     []byte
	offset   int
	readEOF  bool
	c        *Context
}

func (r *repeatableReader) Read(p []byte) (n int, err error) {
	if r.readEOF {
		r.offset = 0
		r.readEOF = false
	}

	if r.rc != nil {
		n, err = r.rc.Read(p)
		if n > 0 {
			r.body = append(r.body, p[:n]...)
		}
		if err == io.EOF {
			r.readEOF = true
			_ = r.rc.Close()
			r.rc = nil
			if r.c != nil {
				r.c.Set(BodyBytesKey, r.body)
			}
		}
		return n, err
	}

	if r.offset >= len(r.body) {
		r.readEOF = true
		return 0, io.EOF
	}
	n = copy(p, r.body[r.offset:])
	r.offset += n
	if r.offset >= len(r.body) {
		r.readEOF = true
		return n, io.EOF
	}
	return n, nil
}

func (r *repeatableReader) Close() error {
	if r.rc != nil {
		return r.rc.Close()
	}
	return nil
}
