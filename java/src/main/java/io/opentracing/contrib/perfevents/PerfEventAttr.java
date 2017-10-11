package io.opentracing.contrib.perfevents;

import com.sun.jna.Structure;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Arrays;
import com.sun.jna.Pointer;

public class PerfEventAttr extends Structure implements Structure.ByReference {
    public int type;
    public int size;
    public long config;
    public long sample_period;
    public long sample_type;
    public long read_format;
    public long bit_attrs;
    public int wakeup_events;
    public int bp_type;
    public long bp_addr;
    public long bp_len;
    public long branch_sample_type;
    public long sample_regs_user;
    public int sample_stack_user;
    public int clockid;
    public long sample_regs_intr;
    public int aux_watermark;
    public int reserved_2;

    public PerfEventAttr() {
    }

    public PerfEventAttr(Pointer p) {
	    super(p);
	    read();
    }

    protected List<String> getFieldOrder() {
	return Arrays.asList(new String[] {
		"type",
		"size",
		"config",
		"sample_period",
		"sample_type",
		"read_format",
		"bit_attrs",
		"wakeup_events",
		"bp_type",
		"bp_addr",
		"bp_len",
		"branch_sample_type",
		"sample_regs_user",
		"sample_stack_user",
		"clockid",
		"sample_regs_intr",
		"aux_watermark",
		"reserved_2"
	    });
    }
}

class OtherBits {
    static final int DISABLED = 0;
    static final int INHERIT = 1;
    static final int PINNED = 2;
    static final int EXCLUSIVE = 3;
    static final int EXCLUDE_USER = 4;
    static final int EXCLUDE_KERNEL = 5;
    static final int EXCLUDE_HV = 6;
    static final int EXCLUDE_IDLE = 7;
    static final int MMAP = 8;
    static final int COMM = 9;
    static final int FREQ = 10;
    static final int INHERIT_STAT = 11;
    static final int ENABLE_ON_EXEC = 12;
    static final int TASK = 13;
    static final int WATERMARK = 14;
    static final int PRECISE_IP1 = 15;
    static final int PRECISE_IP2 = 16;
    static final int MMAP_DATA = 17;
    static final int SAMPLE_ID_ALL = 18;
    static final int EXCLUDE_HOST = 19;
    static final int EXCLUDE_GUEST = 20;
    static final int EXCLUDE_CALLCHAIN_KERNEL = 21;
    static final int EXCLUDE_CALLCHAIN_USER = 22;
    static final int MMAP2 = 23;
    static final int COMM_EXEC = 24;
    static final int USE_CLOCKID = 25;
    static final int CONTEXT_SWITCH = 26;
    static final int RESERVED_1 = 27;
    
    static long set_bit(long value, int pos) {
	    long result;
	    result = value | (1 << pos);
	    return result;
    }
}

/*
* Required macros and definitions for the correct operation of
* perf_event_open on a given platform.
*/
class PerfMacros{
    /*
    * PMU events type definitions (from linux/perf_event.h).
    * Only HARDWARE type is supported as of now.
    */
    public static final int PERF_TYPE_HARDWARE = 0;
    public static final int PERF_TYPE_SOFTWARE = 1;

    /*
    * List of generic events supported (from linux/perf_event.h).
    * All of these events belong to HARDWARE type events.
    */
    public static final int PERF_HW_CPU_CYCLES = 0;
    public static final int PERF_HW_INSTRUCTIONS = 1;
    public static final int PERF_HW_CACHE_REF = 2;
    public static final int PERF_HW_CACHE_MISSES = 3;
    public static final int PERF_HW_BRANCH_INSTRUCTIONS = 4;
    public static final int PERF_HW_BRANCH_MISSES = 5;
    public static final int PERF_HW_BUS_CYCLES = 6;

    /*
    * perf_event_open system call numbers for supported platforms.
    */
    public static final int X86_PERF_EVENT_OPEN = 298;
    public static final int PPC_PERF_EVENT_OPEN = 319;

    /*
    * perf_event_attr size
    */
    public static final int PERF_EVENT_ATTR_SIZE = 112;

    /*
    * IOCTL operation values for different architectures.
    */
    /* For x86_64 */
    public static final int PERF_IOC_RESET_X86 = 0x2403;
    public static final int PERF_IOC_ENABLE_X86 = 0x2400;
    public static final int PERF_IOC_DISABLE_X86 = 0x2401;
    /* For ppc64 */
    public static final int PERF_IOC_RESET_PPC = 0x20002403;
    public static final int PERF_IOC_ENABLE_PPC = 0x20002400;
    public static final int PERF_IOC_DISABLE_PPC = 0x20002401;
}