package chrono

import "github.com/go-co-op/gocron/v2"

type OnceJob struct {
	Name     string
	TaskFunc func()
	Hooks    gocron.EventListener
	err      error
}

func NewOnceJob() *OnceJob {
	return &OnceJob{}
}
