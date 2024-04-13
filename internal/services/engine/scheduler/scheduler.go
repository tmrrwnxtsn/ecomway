package scheduler

import "context"

type BackgroundTask interface {
	Start(ctx context.Context)
}

type Scheduler struct {
	backgroundTasks []BackgroundTask
}

func NewScheduler(tasks ...BackgroundTask) *Scheduler {
	return &Scheduler{
		backgroundTasks: tasks,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	for _, task := range s.backgroundTasks {
		task.Start(ctx)
	}
}
