# perfevents/java
This is the implementation of perfevents package in Java.

## Standalone Installation and Usage
perfevents provides an API and a package in Java which can be used
by any Java application to profile its components.

Note: To use this package, one must have Java Native Access package in their
classpath. If not available, please download it from [here](https://github.com/java-native-access/jna).

To use this package inside your Java application, use:
```java
import io.opentracing.contrib.perfevents.PerfEvent;
```

To open and profile for any event:
```java
// Get the event mappings for the hardware event "cycles".
PerfEvent pe = new PerfEvent("cycles");
try {
    // Start any event
    pe.startEvent();
} catch (UnsupportedOperationException u) {
    // Handle error here
}
```

This creates a new event with an event descriptor (which must be closed). At
any point in time, the event's values can be fetched:
```java
try {
    /* Read the event value in pe.value */
	pe.readEvent();
} catch (UnsupportedOperationException u) {
    // Handle exception here
}
```
Fetch the event value using:
```java
pe.getEventData()
```

If you are done with this event, it must be closed:
```java
try{
	pe.destroyEvent();
} catch (UnsupportedOperationException u) {
    System.out.println("Error encountered: " + u);
}
```

A sample application containing the profiling code:
```java
import io.opentracing.contrib.perfevents.PerfEvent;

public class TestPerf {
    public static void main(String[] args) {
        /* Capture the cpu cycles */
		PerfEvent pe = new PerfEvent("cycles");
		try {
			/* Start the event */
			pe.startEvent();
			System.out.println("Event opened");

			/* Read the event value */
			pe.readEvent();
			System.out.println("CPU cycles for the print statement : " + pe.getEventData());
			pe.destroyEvent();
		} catch (UnsupportedOperationException u) {
			System.out.println("Error encountered: " + u);
		}
	}
}
```

## Supported Architectures
For now, only these two architectures are supported:
- x86_64/amd64
- ppc64le

## Supported Events
For now, 7 Hardware events are supported:
- cpu-cycles
- instructions
- cache-references
- cache-misses
- branch-instructions
- branch-misses
- bus-cycles

## Limitations
The Java programs mostly contain multiple threads. The perf infrastructure (in kernel) can't distinguish between different threads inside JVM. Hence, the values that we get are per process/task.

## TODOs:
- Add support for an observer implementation intended to be used with OpenTracing.
- Support Hardware cache events and dynamic events exposed by the linux kernel.
