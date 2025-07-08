package controller

import (
	"github.com/roksky/bootstrap-api/job"
)

type Provider interface {
	GetControllers() []Controller
	GetJobs() []job.Job
}
