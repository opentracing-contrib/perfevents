package perfevents

import (
	"github.com/opentracing/opentracing-go"
	otobserver "github.com/opentracing-contrib/go-observer"
)

// TODO: Add a member to keep the list of all available events, which
// is initialized when NewObserver() is called.
type Observer struct {}

// New observer creates a new observer
func NewObserver() *Observer {
	return &Observer{}
}

// OnStartSpan creates a new Observer for the span
func (o *Observer) OnStartSpan(sp opentracing.Span, operationName string, options opentracing.StartSpanOptions) (otobserver.SpanObserver, bool) {
	return NewSpanObserver(sp, options)
}

// SpanObserver collects perfevent metrics
type SpanObserver struct {
	sp opentracing.Span
	EventDescs []PerfEventInfo
}

// NewSpanObserver creates a new SpanObserver that can emit perfevent
// metrics
func NewSpanObserver(s opentracing.Span, opts opentracing.StartSpanOptions) (*SpanObserver, bool) {
	so := &SpanObserver{
		sp: s,
	}

	req := false
	for k, v := range opts.Tags {
		if k == "perfevents" {
			so.OnSetTag(k, v)
			req = true
		}
	}

	return so, req
}

func (so *SpanObserver) OnSetOperationName(operationName string) {
}

func (so *SpanObserver) OnSetTag(key string, value interface{}) {
	if key == "perfevents" {
		if v, ok := value.(string); ok {
			_, _, so.EventDescs = InitOpenEventsEnableSelf(v)
		}
	}
}

func (so *SpanObserver) OnFinish(options opentracing.FinishOptions) {
	err := EventsRead(so.EventDescs)
	if err != nil {
		return
	}

	// log and close the perf events first, if any, since, we don't
	// want to account for the code to finish up the span.
	for _, event := range so.EventDescs {
		// In any case of an error for an event, event.EventName
		// will contain "" for an event.
		if event.EventName != "" {
			so.sp.LogEvent(event.EventName + ":" +
				FormatDataToString(event))
		}
	}

	EventsDisableClose(so.EventDescs)
}
