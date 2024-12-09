package proc

import (
	"log"
	"os"
	"os/signal"
	"rudder/config"
	"rudder/resource"
	"syscall"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type JobResources struct {
	Config *config.AppConfig
	DB     *resource.Database
	API    *resource.SimpleFINAPI
	Args   *config.Args
}

func makeJob(sched gocron.Scheduler, crontab string, name string, days int, jobResources JobResources) error {
	_, err := sched.NewJob(
		gocron.CronJob(crontab, false),
		gocron.NewTask(
			Update,
			jobResources.Config,
			jobResources.DB,
			jobResources.API,
			jobResources.Args,
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

func StartScheduler(appConfig *config.AppConfig, db *resource.Database, sfinAPI *resource.SimpleFINAPI, args *config.Args) error {
	jobResources := JobResources{
		Config: appConfig,
		DB:     db,
		API:    sfinAPI,
		Args:   args,
	}

	sched, err := gocron.NewScheduler(gocron.WithLocation(appConfig.Location))
	if err != nil {
		return err
	}

	makeJob(sched, "* * * * *", "hourly", appConfig.HourlyPullDays, jobResources)
	makeJob(sched, "0 5 * * *", "daily", appConfig.DailyPullDays, jobResources)
	makeJob(sched, "0 6 * * 6", "weekly", appConfig.WeeklyPullDays, jobResources)

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
