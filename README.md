# perfevents
perfevents is a repo to monitor and profile programs or their segments.
It also abstracts out the platform differences and gives a generic
interface to start, read and close the platform counters.

A set of events and what we can do with them can be found at :
https://perf.wiki.kernel.org/index.php/Tutorial#Events

Currently, this exports APIs in golang (perfevents/go).