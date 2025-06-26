package job

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type Job struct {
	ID        int
	Name      string
	Schedule  string
	StartDate time.Time
	EndDate   time.Time
	Function  func()
}

type JobExecutor struct {
	Cron      *cron.Cron
	Jobs      map[int]Job
	JobRunLog map[int][]time.Time
}

func NewJobExecutor() *JobExecutor {
	return &JobExecutor{
		Cron:      cron.New(cron.WithSeconds()),
		Jobs:      make(map[int]Job),
		JobRunLog: make(map[int][]time.Time),
	}
}

func (je *JobExecutor) RegisterJob(job Job) {
	je.Jobs[job.ID] = job
	entryID, err := je.Cron.AddFunc(job.Schedule, func() {
		now := time.Now()
		if now.After(job.StartDate) && now.Before(job.EndDate) {
			fmt.Printf("Running job: %s at %v\n", job.Name, now)
			job.Function()
			je.JobRunLog[job.ID] = append(je.JobRunLog[job.ID], now)
		}
	})
	if err != nil {
		fmt.Printf("Failed to schedule job: %v\n", err)
	} else {
		fmt.Printf("Scheduled job %s with entry ID %v\n", job.Name, entryID)
	}
}

func CalculateJobInstances(Schedule string, startDate, endDate time.Time) []time.Time {
	var instances []time.Time
	schedule, err := cron.ParseStandard(Schedule)
	if err != nil {
		fmt.Printf("Invalid schedule format: %v\n", err)
		return instances
	}

	next := schedule.Next(startDate)
	for next.Before(endDate) {
		instances = append(instances, next)
		next = schedule.Next(next)
	}
	return instances
}
