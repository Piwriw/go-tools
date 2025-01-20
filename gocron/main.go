package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func main() {
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
	}

	// add a job to the scheduler
	job, err := s.NewJob(
		gocron.CronJob("54 14 * * *", false),
		gocron.NewTask(
			func(a string, b int) {
				fmt.Println(a, b, time.Now().Format(""))
			},
			"hello",
			1,
		),
	)
	if err != nil {
		// handle error
	}
	// each job has a unique id
	fmt.Println(job.ID())

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute * 5):
	}

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		// handle error
	}
}
