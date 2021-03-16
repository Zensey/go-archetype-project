package x

import (
	"net/http"

	"github.com/ory/x/logrusx"
)

func LogError(r *http.Request, err error, logger *logrusx.Logger) {
	if logger == nil {
		logger = logrusx.New("", "")
	}

	logger.WithRequest(r).
		WithError(err).Errorln("An error occurred")
}
