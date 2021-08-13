package logger

type Logger interface {
	Debugf(message string, args ...interface{})
}

type TestLogger struct {
}

func NewTestLogger() *TestLogger {
	testLogger := &TestLogger{}
	return testLogger
}

func (l *TestLogger) Debugf(message string, args ...interface{}) {
}
