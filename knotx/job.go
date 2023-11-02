package knotx

import "github.com/hootuu/tome/bk"

type JobStatus int

const (
	Submitted JobStatus = 0
	Confirmed JobStatus = 1
	Failed    JobStatus = -1
)

type Job struct {
	Tie    *bk.Tie
	Status JobStatus
}
