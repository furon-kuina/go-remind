package main

import (
	"container/heap"
	"testing"
	"time"
)

func TestOnceSchedule(t *testing.T) {
	cases := []struct {
		name        string
		schedule    time.Time
		now         time.Time
		wantNext    time.Time
		wantIsValid bool
	}{
		{
			name:        "test valid case",
			schedule:    time.Date(2023, 11, 17, 14, 30, 0, 0, time.UTC),
			now:         time.Date(2022, 11, 17, 14, 30, 0, 0, time.UTC),
			wantNext:    time.Date(2023, 11, 17, 14, 30, 0, 0, time.UTC),
			wantIsValid: true,
		},
		{
			name:        "test invalid case",
			schedule:    time.Date(2023, 11, 17, 14, 30, 0, 0, time.UTC),
			now:         time.Date(2024, 11, 17, 14, 30, 0, 0, time.UTC),
			wantNext:    time.Date(2023, 11, 17, 14, 30, 0, 0, time.UTC),
			wantIsValid: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := &OnceSchedule{
				time: tt.schedule,
			}
			next, isValid := s.NextRun(tt.now)
			if isValid != tt.wantIsValid {
				t.Errorf("isValid didn't match: want %t, got %t", tt.wantIsValid, isValid)
				return
			}
			if !next.Equal(tt.wantNext) {
				t.Errorf("next didn't match: want %s, got %s", tt.wantNext.String(), next.String())
			}
		})
	}
}

func TestIntervalSchedule(t *testing.T) {
	oneDay, err := time.ParseDuration("24h")
	if err != nil {
		t.Errorf("time.ParseDuration failed: %v", err)
	}
	oneWeek, err := time.ParseDuration("168h")
	if err != nil {
		t.Errorf("time.ParseDuration failed: %v", err)
	}

	cases := []struct {
		name        string
		from        time.Time
		until       time.Time
		interval    time.Duration
		now         time.Time
		wantNext    time.Time
		wantIsValid bool
	}{
		{
			name:        "test daily interval",
			from:        time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			until:       time.Date(2024, 11, 17, 14, 30, 0, 0, time.UTC),
			interval:    oneDay,
			now:         time.Date(2023, 1, 1, 0, 0, 0, 1, time.UTC),
			wantNext:    time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			wantIsValid: true,
		},
		{
			name:        "test daily interval",
			from:        time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			until:       time.Date(2024, 11, 17, 14, 30, 0, 0, time.UTC),
			interval:    oneDay,
			now:         time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantNext:    time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			wantIsValid: true,
		},
		{
			name:        "test weekly interval",
			from:        time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC),
			until:       time.Date(2024, 11, 17, 14, 30, 0, 0, time.UTC),
			interval:    oneWeek,
			now:         time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantNext:    time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC),
			wantIsValid: true,
		},
		{
			name:        "test from",
			from:        time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC),
			until:       time.Date(2024, 11, 17, 14, 30, 0, 0, time.UTC),
			interval:    oneWeek,
			now:         time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			wantNext:    time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC),
			wantIsValid: true,
		},
		{
			name:        "test until",
			from:        time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC),
			until:       time.Date(2023, 11, 17, 14, 30, 0, 0, time.UTC),
			interval:    oneWeek,
			now:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			wantNext:    time.Time{},
			wantIsValid: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := &IntervalSchedule{
				from:     tt.from,
				until:    tt.until,
				interval: tt.interval,
			}
			next, isValid := s.NextRun(tt.now)
			if isValid != tt.wantIsValid {
				t.Errorf("isValid didn't match: want %t, got %t", tt.wantIsValid, isValid)
				return
			}
			if !next.Equal(tt.wantNext) {
				t.Errorf("next didn't match: want %s, got %s", tt.wantNext.String(), next.String())
			}

		})
	}
}

func TestJobHeap(t *testing.T) {
	h := &jobHeap{}
	heap.Init(h)
	heap.Push(h, Job{
		time: time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC),
	})
	heap.Push(h, Job{
		time: time.Date(2021, 12, 25, 0, 0, 0, 0, time.UTC),
	})
	jobx := heap.Pop(h)
	job, ok := jobx.(Job)
	if !ok {
		t.Fatalf("not Job")
	}
	if !job.time.Equal(time.Date(2021, 12, 25, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("want %s, got %s", time.Date(2021, 12, 25, 0, 0, 0, 0, time.UTC).String(), job.time.String())

	}
	heap.Push(h, Job{
		time: time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC),
	})

	heap.Push(h, Job{
		time: time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
	})
	jobx = heap.Pop(h)
	job, ok = jobx.(Job)
	if !ok {
		t.Fatalf("not Job")
	}
	if !job.time.Equal(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("want %s, got %s", time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC).String(), job.time.String())
	}
}
