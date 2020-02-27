package gologger

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestGoLogger(t *testing.T) {
	t.Parallel()
	var logger = Init(DebugLevel)
	logger.AddAction("print", func(logMsg LogMsg) {
		var counter = 0
		for k := 1; k <= 4; k++ {
			counter++
			t.Logf("hello - %s - %s", logMsg.Time, logMsg.Parameters)
			time.Sleep(50000000)
		}
	})
	logger.AddAction("superprint", func(logMsg LogMsg) {
		var counter = 0
		for k := 1; k <= 4; k++ {
			counter++
			t.Logf("superprinthello - %s - %s", logMsg.Time, logMsg.Parameters)
			time.Sleep(50000000)
		}
	})
	logger.SetLogLevelActions(DebugLevel, []string{"print", "superprint"})
	for i := 0; i < 5; i++ {
		for ix := 0; ix < 5; ix++ {
			logger.Debug(fmt.Sprintf("blya - %d invoke - %d", i, ix))
		}
		time.Sleep(5000000000)
	}
	logger.WaitLogsDone()
}

func TestTime(t *testing.T) {
	t.Log("Hello")

	time.Sleep(20)
	t.Log("Bye")

	//t.Error(logger)
}

func BenchmarkTime(t *testing.B) {

	t.Log("Hello")

	time.Sleep(20)
	t.Log("Bye")

	//t.Error(logger)
}
