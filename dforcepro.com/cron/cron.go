package cron

import (
	"fmt"

	"github.com/robfig/cron"
)

type CronJob interface {
	Enable() bool
	GetJobs() []JobSpec
}

type JobSpec struct {
	Spec string
	Job  func()
}

func StartCronJob(cronJobs ...CronJob) {
	c := cron.New()
	for _, cron := range cronJobs {
		if !cron.Enable() {
			continue
		}
		jobs := cron.GetJobs()
		for _, job := range jobs {
			fmt.Println(job.Spec)
			c.AddFunc(job.Spec, job.Job)
		}
	}
	c.Start()
	select {}
}
