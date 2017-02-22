# perfevents/go
This is the implementation of perfevents package in golang.

## Installation and Usage
To use this package inside your go application, use :
import "github.com/opentracing-contrib/perfevents/go"

To open a list of events :
err, evs, pds := perfevents.InitOpenEventsEnableSelf("cpu-cycles,cache-misses,instructions")

This creates 3 event descriptors pds for cpu-cycles,
cache-misses and instructions and enables them, so they
start counting.

At any point of time, we can read the event values :
err, evs = perfevents.EventsRead(pds)

After we are done monitoring, just close out the events :
err, evs = perfevents.EventsDisableClose(pds)

That's it!

## Events
For now, 7 generic hardware events are supported :
cpu-cycles
instructions
cache-references
cache-misses
branch-instructions
branch-misses
bus-cycles

## TODOs
- Support the hardware cache events.
- Support the dynamic events exported by the kernel.
