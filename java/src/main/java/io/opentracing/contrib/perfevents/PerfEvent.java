package io.opentracing.contrib.perfevents;

import com.sun.jna.Library;
import com.sun.jna.Native;
import com.sun.jna.ptr.PointerByReference;
import com.sun.jna.Pointer;

/*
* Load the clibrary.
*/
interface CLibrary extends Library {
    CLibrary INSTANCE = (CLibrary) Native.loadLibrary("c", CLibrary.class);
    
    /* Need this for perf_event_open call */
    int syscall(int number, Object... args);

    /*
    * Various IOCTL operations needed on the file descriptor to
    * reset, enable and disable the events.
    */
    int ioctl(int fd, int request, int flags);

    /* Need this to read the values of the counters */
    int read(int fd, long[] values, int size);

    /* Need this to close the perf event */
    int close(int fd);
}

public class PerfEvent {
    /* File descriptor */    
    int fd;
    /* Arch specific properties and methods */
    ArchOps ao;
    /* Perf event attributes */
    PerfEventAttr pea;
    /* Values read from the event */
    long value;

    public PerfEvent() {
        pea = new PerfEventAttr();
    }

    /*
    * Takes an event's name as an argument.
    */
    public PerfEvent(String eventName) {
        EventMap em;
        try {
            em = new EventMap(eventName);
        } catch(IllegalArgumentException iae) {
            System.out.println("Error: event unsupported");
            return;
        }
        pea = new PerfEventAttr();
        pea.type = em.type;
        pea.size = PerfMacros.PERF_EVENT_ATTR_SIZE;
        pea.config = em.event;
        pea.bit_attrs = OtherBits.set_bit(pea.bit_attrs, OtherBits.DISABLED);
	    pea.bit_attrs = OtherBits.set_bit(pea.bit_attrs, OtherBits.EXCLUDE_KERNEL);
        pea.bit_attrs = OtherBits.set_bit(pea.bit_attrs, OtherBits.EXCLUDE_HV);
        try {
            ao = new ArchOps();
        } catch(IllegalStateException ise) {
            System.out.println("System platform not supported");
            return;
        }
    }

    /* pid == 0, cpu == 0 */
    public void selfOpenEvent() {
        this.openEvent(pea, 0, -1, -1, 0);
    }

    public void openEvent(PerfEventAttr pea, int pid, int cpu, int group_fd, int flags) {
        fd = CLibrary.INSTANCE.syscall(ao.perfSyscall, pea, 0, -1, -1, 0);
        if (fd == -1) {
            // throw Exception
            throw new UnsupportedOperationException();
        }
    }

    public void ioctlOnEvent(int ioctlOp) {
        int ret = CLibrary.INSTANCE.ioctl(fd, ioctlOp, 0);
        if (ret < 0) {
            // Exception
            throw new UnsupportedOperationException();
        }
    }

    public void enableEvent() {
        ioctlOnEvent(ao.enableOp);
    }

    public void resetEvent() {
        ioctlOnEvent(ao.resetOp);
    }

    public void disableEvent() {
        ioctlOnEvent(ao.disableOp);
    }

    public void readEvent() {
        long[] value_list = new long[1];
        // 64 is the amount of data we need.
        // TODO: replace 64 with a macro or native Long size.
        int ret = CLibrary.INSTANCE.read(fd, value_list, 64);
        if (ret < 0) {
            throw new UnsupportedOperationException();
        }
        value = value_list[0];
    }

    public void closeEvent() {
        int ret = CLibrary.INSTANCE.close(fd);
        if (ret < 0) {
            throw new UnsupportedOperationException();
        }
    }

    /*
    * startEvent opens, resets and enables an event.
    */
    public void startEvent() {
        try {
            selfOpenEvent();
        } catch (UnsupportedOperationException uoe) {
            System.out.println("error in opening event");
        }
        try {
            resetEvent();
        } catch (UnsupportedOperationException uoe) {
            System.out.println("error in reset'ing the event");
        }
        try {
            enableEvent();
        } catch (UnsupportedOperationException uoe) {
            System.out.println("error in enabling the event");
        }
    }

    /*
    * Wrapper to disable and destroy the perf event.
    */
    public void destroyEvent() {
        try {
            disableEvent();
        } catch (UnsupportedOperationException uoe) {
            System.out.println("error in disabling the event");
        }
        try {
            closeEvent();
        } catch (UnsupportedOperationException uoe) {
            System.out.println("error in closing the event");
        }
    }

    public long getEventData() {
        return value;
    }
}

class EventMap {
    int type;
    int event;

    EventMap(String eventName) {
        switch (eventName) {
            case "cycles":
                this.type = PerfMacros.PERF_TYPE_HARDWARE;
                this.event = PerfMacros.PERF_HW_CPU_CYCLES;
                break;
            case "instructions":
                this.type = PerfMacros.PERF_TYPE_HARDWARE;
                this.event = PerfMacros.PERF_HW_INSTRUCTIONS;
                break;
            case "cache-references":
                this.type = PerfMacros.PERF_TYPE_HARDWARE;
                this.event = PerfMacros.PERF_HW_CACHE_REF;
                break;
            case "cache-misses":
                this.type = PerfMacros.PERF_TYPE_HARDWARE;
                this.event = PerfMacros.PERF_HW_CACHE_MISSES;
                break;
            case "branch-instructions":
                this.type = PerfMacros.PERF_TYPE_HARDWARE;
                this.event = PerfMacros.PERF_HW_BRANCH_INSTRUCTIONS;
                break;
            case "branch-misses":
                this.type = PerfMacros.PERF_TYPE_HARDWARE;
                this.event = PerfMacros.PERF_HW_BRANCH_MISSES;
                break;
            case "bus-cycles":
                this.type = PerfMacros.PERF_TYPE_HARDWARE;
                this.event = PerfMacros.PERF_HW_BUS_CYCLES;
                break;
            default:
                this.type = -1;
                this.event = -1;
                throw new IllegalArgumentException();
        }
    }
}