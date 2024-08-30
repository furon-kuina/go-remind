package main

import (
	"container/heap"
	"fmt"
	"time"
)

// Schedule represents a registered schedule
type Schedule interface {
	// NextRun returns the next time to run a job
	// and whether there are any schedule left to run
	NextRun(now time.Time) (time.Time, bool)
}

// OnceSchedule represents a schedule that is run only once
type OnceSchedule struct {
	time time.Time
}

func (s *OnceSchedule) NextRun(now time.Time) (time.Time, bool) {
	return s.time, now.Before(s.time)
}

// IntervalSchedule represents a schedule that is run in fixed time interval
type IntervalSchedule struct {
	from     time.Time
	until    time.Time
	interval time.Duration
}

func (s *IntervalSchedule) NextRun(now time.Time) (time.Time, bool) {
	if s.until.Before(now) {
		return time.Time{}, false
	}
	if s.from.After(now) {
		return s.from, true
	}
	// use UnixNano since time.Duration counts in nanoseconds
	startUnix := s.from.UnixNano()
	nowUnix := now.UnixNano()
	nextUnix := startUnix + ((nowUnix-startUnix)/int64(s.interval)+1)*int64(s.interval)
	fmt.Println("now:", now.String(), "from:", s.from.String(), "interval: ", int64(s.interval))

	return time.Unix(0, nextUnix), now.Before(s.until)
}

// Job represents a task performed at certain time
// TODO: handle timeout
// TODO: handle cancel gracefully
type Job struct {
	task func(time.Time)
	time time.Time
}

type jobHeap []Job

func (h *jobHeap) Len() int {
	return len(*h)
}

func (h *jobHeap) Less(i, j int) bool {
	return (*h)[i].time.Before((*h)[j].time)
}

func (h *jobHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *jobHeap) Push(x any) {
	j, ok := x.(Job)
	if !ok {
		return
	}
	*h = append(*h, j)
}

func (h *jobHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type scheduler struct {
	h         *jobHeap
	schedules []Schedule
}

func (s *scheduler) Add(job *Job) {
	s.h.Push(job)
}

func (s *scheduler) List() []Job {
	return *s.h
}

func (s *scheduler) Delete(jobId int) {
	heap.Remove(s.h, jobId)
}
