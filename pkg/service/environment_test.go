package service_test

import (
	"testing"

	"github.com/moov-io/base/log"
	"github.com/stretchr/testify/assert"

	"github.com/moovfinancial/backendhiring/pkg/service"
)

func Test_Environment_Startup(t *testing.T) {
	a := assert.New(t)

	env := &service.Environment{
		Logger: log.NewDefaultLogger(),
	}

	env, err := service.NewEnvironment(env)
	a.Nil(err)

	t.Cleanup(env.Shutdown)
}
