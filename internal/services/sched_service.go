package services

import (
	"context"
	"log"
	"rudder/internal/config"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type SchedService struct {
	Config    *config.AppConfig
	Args      config.Args
	SFIN      *SimpleFINService
	Cat       *CategoriesService
	Scheduler gocron.Scheduler
}

func NewSchedService(
	c *config.AppConfig,
	a config.Args,
	sfin *SimpleFINService,
	cat *CategoriesService,
) *SchedService {
	return &SchedService{
		Config: c,
		Args:   a,
		SFIN:   sfin,
		Cat:    cat,
	}
}

// Default options for job names
func withName(name string) gocron.JobOption {
	return gocron.WithName(name)
}

// Default options for event listeners
func withEventListeners(s *SchedService) gocron.JobOption {
	return gocron.WithEventListeners(
		gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
			log.Printf("Running %v job...", jobName)
		}),
		gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
			log.Printf("Finished running %v job", jobName)

			// Run autocat job after sync job
			s.Cat.CategorizeTransactions(context.Background())
		}),
		gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
			log.Printf("Error while running %v job: %v\n", jobName, err)
		}),
		gocron.AfterJobRunsWithPanic(func(jobID uuid.UUID, jobName string, recoverData any) {
			log.Printf("Panic while running %v job: %v\n", jobName, recoverData)
		}),
	)
}

func (s *SchedService) makeSyncJob(
	sched gocron.Scheduler,
	crontab string,
	name string,
	days int,
) error {
	_, err := sched.NewJob(
		gocron.CronJob(crontab, false),
		gocron.NewTask(
			s.SFIN.SyncSimpleFIN,
			context.Background(),
			s.Args.UseCached,
			s.Args.SaveCached,
			days,
		),
		withName(name),
		withEventListeners(s),
	)

	return err
}

func (s *SchedService) Start() error {
	sched, err := gocron.NewScheduler(gocron.WithLocation(s.Config.Location))
	if err != nil {
		return err
	}
	s.Scheduler = sched

	s.makeSyncJob(s.Scheduler, "0 * * * *", "hourly", s.Config.HourlyPullDays)
	s.makeSyncJob(s.Scheduler, "0 5 * * *", "daily", s.Config.DailyPullDays)
	s.makeSyncJob(s.Scheduler, "0 6 * * 6", "weekly", s.Config.WeeklyPullDays)

	s.Scheduler.Start()
	for _, job := range sched.Jobs() {
		name := job.Name()
		nextRun, _ := job.NextRun()
		log.Printf("%v job scheduled for %v\n", name, nextRun)
	}

	return nil
}

func (s *SchedService) Stop() error {
	if s.Scheduler != nil {
		return s.Scheduler.Shutdown()
	}
	return nil
}
