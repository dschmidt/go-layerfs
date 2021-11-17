# go-layerfs

This is a simple wrapper around multiple `fs.FS` instances, recursively merging them together dynamically.


If you have two directories, of which one is called `examples/upper` and looks like this:

```
.
├── dir1
│   ├── f11.txt (content: foo)
│   └── f12.txt (content: foo)
├── f1.txt (content: foo)
└── f2.txt (content: foo)
```

and one is called `examples/lower` which looks like this:

```
.
├── dir1
│   ├── f12.txt (content: bar)
│   └── f13.txt (content: bar)
├── f2.txt (content: bar)
└── f3.txt (content: bar)
```

Then `go-layerfs.New(os.DirFS("examples/upper"), os.DirFS("examples/lower"))` will give you a new filesystem with this content:

```
.
├── dir1
│   ├── f11.txt (content: foo)
│   ├── f12.txt (content: foo)
│   └── f13.txt (content: bar)
├── f1.txt (content: foo)
├── f2.txt (content: foo)
└── f3.txt (content: bar)
```

# Example usage

You can run `examples/file_server.go` like this:

```bash
go run examples/file_server.go

2021/11/17 22:59:22 Listening on :8090...
```

Then requests via `httpie` should give you these results:

```
http GET http://localhost:8090/files
HTTP/1.1 200 OK
Content-Length: 123
Content-Type: text/html; charset=utf-8
Date: Wed, 17 Nov 2021 22:03:21 GMT
Last-Modified: Wed, 17 Nov 2021 21:55:53 GMT

<pre>
<a href="dir1/">dir1/</a>
<a href="f1.txt">f1.txt</a>
<a href="f2.txt">f2.txt</a>
<a href="f3.txt">f3.txt</a>
</pre>
```

```
http GET http://localhost:8090/files/f1.txt
HTTP/1.1 200 OK
Accept-Ranges: bytes
Content-Length: 3
Content-Type: text/plain; charset=utf-8
Date: Wed, 17 Nov 2021 22:05:29 GMT
Last-Modified: Wed, 17 Nov 2021 21:56:26 GMT

foo
```

```
http GET http://localhost:8090/files/f3.txt
HTTP/1.1 200 OK
Accept-Ranges: bytes
Content-Length: 3
Content-Type: text/plain; charset=utf-8
Date: Wed, 17 Nov 2021 22:05:56 GMT
Last-Modified: Wed, 17 Nov 2021 21:56:30 GMT

bar
```
