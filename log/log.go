package log

import (
	"fmt"

	"github.com/cloudwego/kitex/pkg/klog"
)

const (
	prefix = "KitexTool"
)

func WithPrefix(v interface{}) string {
	return fmt.Sprintf("%s %v", prefix, v)
}

func WithPrefixList(v ...interface{}) []interface{} {
	if len(v) > 0 {
		v[0] = WithPrefix(v[0])
	}

	return v
}

func Fatal(v ...interface{}) {
	klog.Fatal(WithPrefixList(v...)...)
}

func Error(v ...interface{}) {
	klog.Error(WithPrefixList(v...)...)
}

func Warn(v ...interface{}) {
	klog.Warn(WithPrefixList(v...)...)
}

func Notice(v ...interface{}) {
	klog.Notice(WithPrefixList(v...)...)
}

func Info(v ...interface{}) {
	klog.Info(WithPrefixList(v...)...)
}

func Debug(v ...interface{}) {
	klog.Debug(WithPrefixList(v...)...)
}

func Trace(v ...interface{}) {
	klog.Trace(WithPrefixList(v...)...)
}

func Fatalf(format string, v ...interface{}) {
	klog.Fatalf(WithPrefix(format), v...)
}

func Errorf(format string, v ...interface{}) {
	klog.Errorf(WithPrefix(format), v...)
}

func Warnf(format string, v ...interface{}) {
	klog.Warnf(WithPrefix(format), v...)
}

func Noticef(format string, v ...interface{}) {
	klog.Noticef(WithPrefix(format), v...)
}

func Infof(format string, v ...interface{}) {
	klog.Infof(WithPrefix(format), v...)
}

// Debugf calls the default logger's Debugf method.
func Debugf(format string, v ...interface{}) {
	klog.Debugf(WithPrefix(format), v...)
}

func Tracef(format string, v ...interface{}) {
	klog.Tracef(WithPrefix(format), v...)
}
