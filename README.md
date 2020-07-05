# loggergo
Functions to configure the logrus logging framework for quick setup of other projects

## Usage
```go
package main

import (
	"context"
	
	"github.com/aramonc/loggergo"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

func main() {
	logger := logrus.New()

	// sets the JSON logrus formatter with the following property names
	// 
	// * time -> @timestamp
	// * message -> message
	// * level -> level_name
	logger = loggergo.WithJSONFormatter(logger)

	// sets minimum log record level to debug
	logger = loggergo.WithLevel(logger, "debug")

	_, span := trace.StartSpan(context.Background(), "example")
	defer span

	// sets trace & span ID fields on the logger converting it to an entry
	logger = loggergo.WithTrace(
		logger, 
		span.SpanContext().TraceID.String(), 
		span.SpanContext().SpanID.String(),
	)

	logger.Info("a log record")
}
```
