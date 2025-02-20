package logger

import (
	"testing"

	"multitenant/pkg/config"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	logger, err := NewLogger(&config.Config{
		Logger: config.Logger{},
	})
	require.NoError(t, err)
	require.NotNil(t, logger)
}
