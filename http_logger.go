package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type logrusFormatter struct {
	logger *logrus.Logger
}
type logrusEntry struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

func (formatter *logrusFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := formatter.logger.WithFields(logrus.Fields{
		"method": r.Method,
		"uri":    r.RequestURI,
	})

	return &logrusEntry{
		logger: formatter.logger,
		entry:  entry,
	}
}

func (entry *logrusEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	e := entry.entry
	e = e.WithField("status", status)
	e = e.WithField("elapsed", elapsed.Milliseconds())
	e = e.WithField("size", bytes)
	e = e.WithTime(time.Now())

	e.Info("Incoming HTTP request")
}

func (entry *logrusEntry) Panic(v interface{}, stack []byte) {
	entry.logger.Panic("Request failed")
}
