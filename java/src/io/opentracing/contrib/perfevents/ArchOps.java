package io.opentracing.contrib.perfevents;

class ArchOps {
    String archName;
    int resetOp;
    int enableOp;
    int disableOp;
    int perfSyscall;

    ArchOps() {
        archName = System.getProperty("os.arch");
        switch (archName) {
            case "amd64":
                perfSyscall = PerfMacros.X86_PERF_EVENT_OPEN;
                resetOp = PerfMacros.PERF_IOC_RESET_X86;
                enableOp = PerfMacros.PERF_IOC_ENABLE_X86;
                resetOp = PerfMacros.PERF_IOC_DISABLE_X86;
                break;
            case "ppc64le":
                perfSyscall = PerfMacros.PPC_PERF_EVENT_OPEN;
                resetOp = PerfMacros.PERF_IOC_RESET_PPC;
                enableOp = PerfMacros.PERF_IOC_ENABLE_PPC;
                disableOp = PerfMacros.PERF_IOC_DISABLE_PPC;
                break;
            default:
                throw new IllegalStateException();
        }
    }
}