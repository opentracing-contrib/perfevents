# perfevents/go
This is the implementation of perfevents package in golang.

## Standalone Installation and Usage
perfevents provides an API and a library which can be used by a go application to profile its
components.

To use this package inside your go application, use :

```go
import "github.com/opentracing-contrib/perfevents/go"
```

To open a list of events :

```go
err, evs, pds := perfevents.InitOpenEventsEnableSelf("cpu-cycles,cache-misses,instructions")
```

This creates 3 event descriptors `pds` for cpu-cycles,
cache-misses and instructions and enables them, so they
start counting.

At any point in time, we can read the event values :

```go
err, evs = perfevents.EventsRead(pds)
```

After we are done monitoring, just close out the events :

```go
err, evs = perfevents.EventsDisableClose(pds)
```

That's it!

## Usage with OpenTracing
[OpenTracing](http://opentracing.io/) is a vendor neutral open standard for distributed tracing, which
basically means, it provides standard and vendor-neutral APIs for popular platforms, i.e., popular
tracing backends implement the OpenTracing API. More on OpenTracing : http://opentracing.io/documentation/.

perfevents can be used with the popular tracing implementations like [Zipkin](http://zipkin.io/) and
[Jaeger](http://jaeger.readthedocs.io/en/latest/) via the [go-observer](https://github.com/opentracing-contrib/go-observer) interface.

### Usage with Zipkin
In the application where zipkin is initialized, a new perfevents observer must be created. This observer
is then assigned to zipkin.

First import the perfevents/go package:

```go
import perfevents "github.com/opentracing-contrib/perfevents/go"
```

Initialize a perfevents observer:

```go
observer := perfevents.NewObserver()
```

And then, pass this new observer as part of initialization of zipkin:

```go
tracer, _ := zipkin.NewTracer(..., zipkin.WithObserver(observer))
```

Now, to start a span with metrics like cache-misses and cycles :

```go
sp := tracer.StartSpan("name", opentracing.Tag{"perfevents", "cpu-cycles,cache-misses"})
```

With this, the results can be seen in zipkin's UI.

## Supported Events
For now, 7 generic hardware events are supported :
* cpu-cycles
* instructions
* cache-references
* cache-misses
* branch-instructions
* branch-misses
* bus-cycles

## Supported Tracers
perfevents is supported with distributed tracers which:
* is OpenTracing compliant and,
* provide the [go-observer](https://github.com/opentracing-contrib/go-observer) interface.

Right now, the above conditions are satisfied by Zipkin and Jaeger, with jaeger satisfying
the second condition indirectly.

## TODOs
- Support the hardware cache events.
- Support the dynamic events exported by the kernel.
