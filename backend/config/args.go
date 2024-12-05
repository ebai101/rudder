package config

import "flag"

type Args struct {
	SaveCached  bool
	UseCached   bool
	DaysToFetch int
	Scheduled   bool
}

func ParseArgs() Args {
	saveCached := flag.Bool("saveCached", false, "save response data to JSON")
	useCached := flag.Bool("useCached", false, "use cached response data instead of fetching from SimpleFIN")
	daysToFetch := flag.Int("days", 7, "number of days to fetch")
	scheduled := flag.Bool("sched", false, "run cron tasks instead of single update")
	flag.Parse()

	return Args{
		SaveCached:  *saveCached,
		UseCached:   *useCached,
		DaysToFetch: *daysToFetch,
		Scheduled:   *scheduled,
	}
}
