package chrono

import (
	"testing"
	"time"

	"github.piwriw.go-tools/pkg/gron/gocron"
)

func Test(t *testing.T) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		t.Fatal(err)
	}
	job, err := scheduler.AddCronJob(gocron.DayTimeToCron(time.Now().Add(time.Second*10)), func() {
		t.Log("hello")
	})
	if err != nil {
		t.Fatal(err)
	}
	scheduler.Start()
	t.Log(job.ID())
	t.Log(job.Name())
	run, err := job.NextRun()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(run.Format("2006-01-02 15:04:05"))
	lastRun, err := job.LastRun()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lastRun.Format("2006-01-02 15:04:05"))

	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 10):
	}
}
