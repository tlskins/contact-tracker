package http

import "github.com/sirupsen/logrus"

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Logger
}

func NewLogger() *StandardLogger {
	var baseLogger = logrus.New()
	var standardLogger = &StandardLogger{baseLogger}
	standardLogger.Formatter = &logrus.JSONFormatter{}

	return standardLogger
}

// Declare variables to store log messages as new Events
var (
	httpErrMessage         = Event{1, "An error occured starting HTTP listener at port: %s. Error: %s"}
	httpInfoMessage        = Event{2, "Starting HTTP service at %s"}
	invalidArgMessage      = Event{3, "Invalid arg: %s"}
	invalidArgValueMessage = Event{4, "Invalid value for argument: %s: %v"}
	missingArgMessage      = Event{5, "Missing arg: %s"}
	clientCreationErr      = Event{6, "Error creating the %s client. Error: %s"}
)

func (l *StandardLogger) InvalidArg(argumentName string) {
	l.Errorf(invalidArgMessage.message, argumentName)
}

func (l *StandardLogger) InvalidArgValue(argumentName string, argumentValue string) {
	l.Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
}

func (l *StandardLogger) MissingArg(argumentName string) {
	l.Errorf(missingArgMessage.message, argumentName)
}

func (l *StandardLogger) HttpError(port string, err string) {
	l.Errorf(httpErrMessage.message, port, err)
}

func (l *StandardLogger) HttpInfo(port string) {
	l.Infof(httpInfoMessage.message, port)
}

func (l *StandardLogger) ClientErr(pkgName string, err string) {
	l.Errorf(clientCreationErr.message, pkgName, err)
}
