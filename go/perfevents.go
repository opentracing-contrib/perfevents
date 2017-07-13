// Copyright (c) 2017 IBM Corp. Rights Reserved.
// This project is licensed under the Apache License 2.0, see LICENSE.

package perfevents

import (
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// To create an event in the kernel which can monitor/profile for
// a process, we need to use the system call perf_event_open.
// To use this system call, we need to define its arguments.
// perf_event_open is called with perf_event_attr, pid, cpu and
// other arguments. perf_event_attr is a structure containing
// the attributes for an event. This basically encapsulates the type
// of the event, the config value for the event, whether it should
// start in enabled/disabled mode etc. The second argument is the pid
// i.e., for which process we want to profile this event for. It
// can be the current process, for which its value must be 0 or it
// can be any other process's pid, if this process has appropriate
// perfmissions to monitor that process.
// The third argument is "cpu". If we want to restrict our event
// counting to say only cpu 1 for the current process, then, we
// need to specify 1 for the cpu argument. If we don't want to
// restrict the event counting to any one cpu, then, we specify
// the value of this argument to be -1.

// PerfEventAttr structure is derived from linux/perf_event.h
// This struct defines various attributes for a perf event. This is what
// is set and sent to the linux kernel to create an event.
//
// Note : For the unions, only one member has been taken, viz.
// sample_period, wakeup_events, config1 and config2 are the members
// of a union. There are two members in each of the mentioned unions
// and both of them have the same size in all of the unions.
type PerfEventAttr struct {
	type_hw            uint32
	size_s             uint32
	config             uint64
	sample_period      uint64
	sample_type        uint64
	read_format        uint64
	properties         uint64
	wakeup_events      uint32
	bp_type            uint32
	config1            uint64
	config2            uint64
	branch_sample_type uint64
	sample_regs_user   uint64
	sample_stack_user  uint32
	clockid            int32
	sample_regs_intr   uint64
	aux_watermark      uint32
	reserved_2         uint32
}

// Bit fields for the PerfEventAttr.properties value derived from
// linux/perf_event.h
// Each of these bits specify how do we want to start our counter.
const (
	DISABLED                 = 0 // Starts from bit value 0
	INHERIT                  = 1
	PINNED                   = 2
	EXCLUSIVE                = 3
	EXCLUDE_USER             = 4
	EXCLUDE_KERNEL           = 5
	EXCLUDE_HV               = 6
	EXCLUDE_IDLE             = 7
	MMAP                     = 8
	COMM                     = 9
	FREQ                     = 10
	INHERIT_STAT             = 11
	ENABLE_ON_EXEC           = 12
	TASK                     = 13
	WATERMARK                = 14
	PRECISE_IP1              = 15
	PRECISE_IP2              = 16
	MMAP_DATA                = 17
	SAMPLE_ID_ALL            = 18
	EXCLUDE_HOST             = 19
	EXCLUDE_GUEST            = 20
	EXCLUDE_CALLCHAIN_KERNEL = 21
	EXCLUDE_CALLCHAIN_USER   = 22
	MMAP2                    = 23
	COMM_EXEC                = 24
	USE_CLOCKID              = 25
	CONTEXT_SWITCH           = 26
	RESERVED_1               = 27
)

// PMU hardware type definitions (from linux/perf_event.h)
// Only HARDWARE type is supported as of now.
const (
	PERF_TYPE_HARDWARE = 0
	PERF_TYPE_SOFTWARE = 1
)

// List of generic events supported (from linux/perf_event.h)
// All of these events belong to HARDWARE type events.
const (
	PERF_HW_CPU_CYCLES          = 0
	PERF_HW_INSTRUCTIONS        = 1
	PERF_HW_CACHE_REF           = 2
	PERF_HW_CACHE_MISSES        = 3
	PERF_HW_BRANCH_INSTRUCTIONS = 4
	PERF_HW_BRANCH_MISSES       = 5
	PERF_HW_BUS_CYCLES          = 6
)

// EventConfigType : The configuration struct for an event
type EventConfigType struct {
	typeHw uint32
	config uint64
}

var PerfIOCError = errors.New("error in IOCTL call")
var PerfOpenError = errors.New("error in opening event")
var PerfCloseError = errors.New("error in closing event")
var PerfTooManyEvents = errors.New("too many events requested")
var PerfUnsupportedEvent = errors.New("event(s) not supported")
var PerfFdError = errors.New("incorrect file descriptor for event")
var PerfReadError = errors.New("error in reading event data")

// Initializes the event list.
// For now, we support only 7 generic hardware events.
func initEventList() map[string]EventConfigType {
	return map[string]EventConfigType{
		"cpu-cycles":          {PERF_TYPE_HARDWARE, PERF_HW_CPU_CYCLES},
		"instructions":        {PERF_TYPE_HARDWARE, PERF_HW_INSTRUCTIONS},
		"cache-references":    {PERF_TYPE_HARDWARE, PERF_HW_CACHE_REF},
		"cache-misses":        {PERF_TYPE_HARDWARE, PERF_HW_CACHE_MISSES},
		"branch-instructions": {PERF_TYPE_HARDWARE, PERF_HW_BRANCH_INSTRUCTIONS},
		"branch-misses":       {PERF_TYPE_HARDWARE, PERF_HW_BRANCH_MISSES},
		"bus-cycles":          {PERF_TYPE_HARDWARE, PERF_HW_BUS_CYCLES},
	}
}

// Sets up the perf event attributes for a particular eventConfig having
// the type of the event and the config value.
func setupPerfEventAttr(eventConfig EventConfigType) PerfEventAttr {
	var eventAttr PerfEventAttr
	eventAttr.type_hw = eventConfig.typeHw
	eventAttr.config = eventConfig.config
	eventAttr.size_s = uint32(unsafe.Sizeof(eventAttr))
	eventAttr.properties = setBit(eventAttr.properties, DISABLED)
	eventAttr.properties = setBit(eventAttr.properties, EXCLUDE_KERNEL)
	eventAttr.properties = setBit(eventAttr.properties, EXCLUDE_HV)

	return eventAttr
}

// Fetches the event attributes for a specified event string.
func fetchPerfEventAttr(event string) (PerfEventAttr, error) {
	var eventAttr PerfEventAttr
	evList := initEventList()
	evConf, ok := evList[event]
	if ok == false {
		return eventAttr, PerfUnsupportedEvent
	}
	return setupPerfEventAttr(evConf), nil
}

// Perf IOCTL operations for x86
const (
	PERF_IOC_RESET_X86   = 0x2403
	PERF_IOC_ENABLE_X86  = 0x2400
	PERF_IOC_DISABLE_X86 = 0x2401
)

// Perf IOCTL operations for powerpc
const (
	PERF_IOC_RESET_PPC   = 0x20002403
	PERF_IOC_ENABLE_PPC  = 0x20002400
	PERF_IOC_DISABLE_PPC = 0x20002401
)

// PerfIOCOps stores the correct IOC operations respective
// to the underlying architecture.
type PerfIOCOps struct {
	reset uint64
	enable uint64
	disable uint64
}

// PerfEventInfo holds the file descriptor for a perf event.
// EventName : name of the perf event
// Fd : File descriptor opened by the perf_event_open syscall.
// Data : Contains the event data after performing a read on Fd.
type PerfEventInfo struct {
	EventName string
	Fd        int
	Data      uint64
	IOCOps    PerfIOCOps
}

func findMachineInfo() (string, error) {
	var buf syscall.Utsname
	err := syscall.Uname(&buf)
	if err != nil {
		return "", err
	}
	machineName := make([]byte, len(buf.Machine))

	i := 0
	for ; i < len(buf.Machine); i++ {
		if buf.Machine[i] == 0 {
			break
		}
		machineName[i] = uint8(buf.Machine[i])
	}

	str := string(machineName[0:i])
	return str, nil
}

// InitIOCOps initializes the Perf IOCTL functions respective to
// the underlying architecture
func (event *PerfEventInfo) InitIOCOps() error {
	machine, err := findMachineInfo()
	if err != nil {
		return err
	}
	if machine == "x86_64" {
		event.IOCOps = PerfIOCOps{reset: PERF_IOC_RESET_X86, enable: PERF_IOC_ENABLE_X86, disable: PERF_IOC_DISABLE_X86}
	} else if machine == "ppc64le" {
		event.IOCOps = PerfIOCOps{reset: PERF_IOC_RESET_PPC, enable: PERF_IOC_ENABLE_PPC, disable: PERF_IOC_DISABLE_PPC}
	} else {
		return errors.New("InitIOCOps: machine not supported")
	}
	return nil
}

// FetchPerfEventAttr is the same as that of the independent one, just to
// maintain consistency, this method is defined
func (event *PerfEventInfo) FetchPerfEventAttr(eventName string) (error, PerfEventAttr) {
	var eventAttr PerfEventAttr
	eventAttr, err := fetchPerfEventAttr(eventName)
	if err == PerfUnsupportedEvent {
		event.Fd = -1
		event.Data = 0
	}
	return err, eventAttr
}

// InitOpenEventEnable fetches the perf event attributes for event
// "string", opens the event, resets and then enables the event.
func (event *PerfEventInfo) InitOpenEventEnable(eventName string, pid int, cpu int, group_fd int, flags uint64) error {
	err, eventAttr := event.FetchPerfEventAttr(eventName)
	if err != nil {
		return err
	}
	err = event.InitIOCOps()
	if (err != nil) {
		return err
	}

	err = event.OpenEvent(eventAttr, pid, cpu, group_fd, flags)
	if err != nil {
		return err
	}
	event.EventName = eventName
	err = event.ResetEvent()
	if err != nil {
		return err
	}

	err = event.EnableEvent()
	if err != nil {
		return err
	}

	return nil
}

// InitOpenEventEnableSelf opens, enables an event for self process
func (event *PerfEventInfo) InitOpenEventEnableSelf(eventName string) error {
	return event.InitOpenEventEnable(eventName, 0, -1, -1, 0)
}

func filterOutDuplicates(events string) map[string]int {
	names := strings.Split((events), ",")
	count := 0
	eventList := make(map[string]int)
	for i := 0; i < len(names); i++ {
		eventList[names[i]] = count
		count++
	}
	return eventList
}

// InitOpenEventsEnableSelf opens, enables an event list provided in
// "events" string.
// "events" is a comma separated list of supported events.
// In case of an error, where it couldn't create some or all of the required
// events in "events", it sends the error and the error'ed events along
// with the events which it managed to create.
func InitOpenEventsEnableSelf(events string) (error, []string, []PerfEventInfo) {
	eventList := filterOutDuplicates(events)
	tmp := make([]PerfEventInfo, len(eventList))
	eventListNA := make([]string, 0, len(eventList))

	i := 0

	for key, _ := range eventList {
		err := (&tmp[i]).InitOpenEventEnableSelf(key)
		if err != nil {
			eventListNA = append(eventListNA, key)
		}
		i++
	}

	var eventDescs = make([]PerfEventInfo, len(eventList)-len(eventListNA))
	j := 0
	for i := 0; i < len(eventList); i++ {
		if tmp[i].Fd != -1 {
			eventDescs[j] = tmp[i]
			j++
		}
	}

	if len(eventListNA) != 0 {
		return PerfUnsupportedEvent, eventListNA, eventDescs
	}
	return nil, eventListNA, eventDescs
}

// EventsRead : Read the event count for a slice of event descriptors in
// "eventsInfo'
func EventsRead(eventsInfo []PerfEventInfo) error {
	eventListNA := make([]string, len(eventsInfo))

	for i := 0; i < len(eventsInfo); i++ {
		err := (&eventsInfo[i]).ReadEvent()
		if err != nil {
			// Error in reading this event
			eventListNA = append(eventListNA, eventsInfo[i].EventName)
		}
	}
	if len(eventListNA) != 0 {
		errEvents := strings.Join(eventListNA, ",")
		return errors.New("couldn't read events' data for: " + errEvents)
	}
	return nil
}

// EventsDisableClose : Disable and close all the events in the slice
// "eventsInfo'
func EventsDisableClose(eventsInfo []PerfEventInfo) error {
	eventListNA := make([]string, len(eventsInfo))
	for _, eventInfo := range eventsInfo {
		err := (&eventInfo).DisableClose()
		if err != nil {
			eventListNA = append(eventListNA, eventInfo.EventName)
			return err
		}
	}
	if len(eventListNA) != 0 {
		errEvents := strings.Join(eventListNA, ",")
		return errors.New("couldn't close events: " + errEvents)
	}
	return nil
}

// DisableClose disables the event and then closes it.
func (event *PerfEventInfo) DisableClose() error {
	// File descriptor not set?
	if event.Fd < 0 {
		return PerfFdError
	}

	err := event.DisableEvent()
	if err != nil {
		return err
	}

	errClose := syscall.Close(int(event.Fd))
	if errClose != nil {
		return PerfCloseError
	}

	return nil
}

// OpenEvent opens an event
func (event *PerfEventInfo) OpenEvent(eventAttr PerfEventAttr, pid int, cpu int, group_fd int, flags uint64) error {
	// File descriptor already set?
	if event.Fd > 0 {
		return PerfFdError
	}
	fd, _, err := syscall.Syscall6(syscall.SYS_PERF_EVENT_OPEN, uintptr(unsafe.Pointer(&eventAttr)), uintptr(pid), uintptr(cpu), uintptr(group_fd), uintptr(flags), uintptr(0))
	if err > 0 {
		return PerfOpenError
	}
	if int(fd) == -1 {
		return PerfOpenError
	}
	event.Fd = int(fd)
	return nil
}

// ResetEvent resets an event
func (event *PerfEventInfo) ResetEvent() error {
	if event.Fd < 0 {
		return PerfFdError
	}
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(event.Fd), uintptr(event.IOCOps.reset), uintptr(0), uintptr(0), uintptr(0), uintptr(0))
	if err != 0 {
		return PerfIOCError
	}
	return nil
}

// EnableEvent enables an event
func (event *PerfEventInfo) EnableEvent() error {
	if event.Fd < 2 {
		return PerfFdError
	}
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(event.Fd), uintptr(event.IOCOps.enable), uintptr(0), uintptr(0), uintptr(0), uintptr(0))
	if err != 0 {
		return PerfIOCError
	}
	return nil
}

// DisableEvent disables an event
func (event *PerfEventInfo) DisableEvent() error {
	if event.Fd < 2 {
		return PerfFdError
	}
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(event.Fd), uintptr(event.IOCOps.disable), uintptr(0), uintptr(0), uintptr(0), uintptr(0))
	if err != 0 {
		return PerfIOCError
	}
	return nil
}

// ReadEvent reads the event count
func (event *PerfEventInfo) ReadEvent() error {
	readBuf := make([]byte, 8, 10)
	_, err := syscall.Read(event.Fd, readBuf)
	if err != nil {
		return PerfReadError
	}
	data := binary.LittleEndian.Uint64(readBuf)
	event.Data = data
	return nil
}

func setBit(properties uint64, bitPos uint64) uint64 {
	properties |= (1 << bitPos)
	return properties
}

// FormatDataToString converts the data for an event to string
func FormatDataToString(pi PerfEventInfo) string {
	return strconv.Itoa(int(pi.Data))
}
