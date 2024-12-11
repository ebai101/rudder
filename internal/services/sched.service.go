package services

import (
	"context"
	"log"
	"os"
	"os/signal"
	"rudder/internal/config"
	"syscall"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

func NewSchedService(
	c *config.AppConfig,
	a *config.Args,
	sfin *SimpleFINService,
) *SchedService {
	return &SchedService{
		Config: c,
		Args:   a,
		SFIN:   sfin,
	}
}

type SchedService struct {
	Config *config.AppConfig
	Args   *config.Args
	SFIN   *SimpleFINService
}

func (s SchedService) makeJob(sched gocron.Scheduler, crontab string, name string, days int) error {
	_, err := sched.NewJob(
		gocron.CronJob(crontab, false),
		gocron.NewTask(
			s.SFIN.SyncSimpleFIN,
			context.Background(),
			s.Args.UseCached,
			s.Args.SaveCached,
			days,
		),
		gocron.WithName(name),
		gocron.WithEventListeners(
			gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
				log.Printf("Running %v job...", jobName)
			}),
			gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
				log.Printf("Finished running %v job", jobName)
			}),
			gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
				log.Printf("Error while running %v job: %v\n", jobName, err)
			}),
			gocron.AfterJobRunsWithPanic(func(jobID uuid.UUID, jobName string, recoverData any) {
				log.Printf("Panic while running %v job: %v\n", jobName, recoverData)
			}),
		),
	)
	return err
}

func (s SchedService) StartScheduler() error {
	sched, err := gocron.NewScheduler(gocron.WithLocation(s.Config.Location))
	if err != nil {
		return err
	}

	s.makeJob(sched, "* * * * *", "hourly", s.Config.HourlyPullDays)
	s.makeJob(sched, "0 5 * * *", "daily", s.Config.DailyPullDays)
	s.makeJob(sched, "0 6 * * 6", "weekly", s.Config.WeeklyPullDays)

	sched.Start()
	for _, job := range sched.Jobs() {
		name := job.Name()
		nextRun, _ := job.NextRun()
		log.Printf("%v job scheduled for %v\n", name, nextRun)
	}

	// block until sigint/sigterm
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	log.Println("Shutting down")
	err = sched.Shutdown()
	return err
}
