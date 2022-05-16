package log

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"runtime/debug"
	"strings"

	gerrors "github.com/go-errors/errors"
	perrors "github.com/pkg/errors"
)

const (
	errorSkipFrames = 4
)

// "github.com/pkg/errors" supports this interface for retrieving stack trace on an error
type stackTracer interface {
	StackTrace() perrors.StackTrace
}

// ErrorValues represents values from the system environment
type ErrorValues struct {
}

func newErrorValues() *ErrorValues {
	return &ErrorValues{}
}

func (df ErrorValues) getErrorValues(err error) logrus.Fields {
	errorMessage := strings.TrimSpace(err.Error())

	stats := &debug.GCStats{}
	stack := df.getErrorStackTrace(err)
	debug.ReadGCStats(stats)

	fields := logrus.Fields{}
	fields["exception"] = map[string]any{
		"error": errorMessage,
		"trace": stack,
		"gc_stats": map[string]any{
			"last_gc":        stats.LastGC,
			"num_gc":         stats.NumGC,
			"pause_total":    stats.PauseTotal,
			"pause_history":  stats.Pause,
			"pause_end":      stats.PauseEnd,
			"page_quantiles": stats.PauseQuantiles,
		},
	}

	return fields
}

func (df ErrorValues) getErrorStackTrace(err error) string {
	// is it the standard google error type?
	var se *gerrors.Error
	if errors.As(err, &se) {
		return string(se.Stack())
	}

	// does it support a Stack interface?
	var ews stackTracer
	if errors.As(err, &ews) {
		return df.getStackTracer(ews)
	}

	// skip 4 frames that belong to glamplify
	return df.getCurrentStack(errorSkipFrames)
}

func (df ErrorValues) getStackTracer(ews stackTracer) string {
	frames := ews.StackTrace()

	buf := bytes.Buffer{}
	for _, f := range frames {
		s := fmt.Sprintf("%+s:%d\n", f, f)
		buf.WriteString(s)
	}

	return buf.String()
}

func (df ErrorValues) getCurrentStack(skip int) string {
	stack := make([]uintptr, gerrors.MaxStackDepth)
	length := runtime.Callers(skip, stack[:])
	stack = stack[:length]

	buf := bytes.Buffer{}
	for _, pc := range stack {
		frame := gerrors.NewStackFrame(pc)
		buf.WriteString(frame.String())
	}

	return buf.String()
}
