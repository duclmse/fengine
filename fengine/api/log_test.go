package api

import (
	"fmt"
	"github.com/duclmse/fengine/pkg/logger"
	"os"
	"testing"
	"time"
)

func TestElapse(t *testing.T) {
	l, err := logger.New(os.Stdout, "debug")
	if err != nil {
		return
	}
	err = Execute(l)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Execute(log logger.Logger) (err error) {
	defer log.Elapse("ExecuteService %s", "a")(time.Now(), &err)
	return task()
}

func task() error {
	time.Sleep(time.Second)
	return nil
}
