package sync

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	scheduleKeyPrefix  = "SCHEDULE_"
	pagerDutyTokenKey  = "PAGERDUTY_TOKEN"
	slackTokenKey      = "SLACK_TOKEN"
	slackBotTokenKey   = "SLACK_BOT_TOKEN"
	runInterval        = "RUN_INTERVAL_SECONDS"
	pdScheduleFromKey  = "PAGERDUTY_SCHEDULE_FROM"
	runIntervalDefault = 60
)

// Config is used to configure application
// PagerDutyToken - token used to connect to pagerduty API
// SlackToken - token used to connect to Slack API
type Config struct {
	Schedules             []Schedule
	PagerDutyToken        string
	SlackToken            string
	SlackBotToken		  string
	RunIntervalInSeconds  int
	PagerdutyScheduleFrom time.Duration
}

// Schedule models a PagerDuty schedule that will be synced with Slack
// ScheduleIDs - All PagerDuty schedule ID's to sync
// CurrentOnCallGroupName - Slack group name for current person on call
type Schedule struct {
	ScheduleIDs            []string
	CurrentOnCallGroupName string
}

// NewConfigFromEnv is a function to generate a config from env varibles
// PAGERDUTY_TOKEN - PagerDuty Token
// SLACK_TOKEN - Slack Token
// SCHEDULE="id" e.g. 1234
func NewConfigFromEnv() (*Config, error) {
	config := &Config{
		PagerDutyToken:       os.Getenv(pagerDutyTokenKey),
		SlackToken:           os.Getenv(slackTokenKey),
		SlackBotToken:        os.Getenv(slackBotTokenKey),
		RunIntervalInSeconds: runIntervalDefault,
	}

	runInterval := os.Getenv(runInterval)
	v, err := strconv.Atoi(runInterval)
	if err == nil {
		config.RunIntervalInSeconds = v
	}

	pagerdutyScheduleFrom, err := getPagerdutyScheduleTime(pdScheduleFromKey)
	if err != nil {
		return nil, err
	}
	config.PagerdutyScheduleFrom = pagerdutyScheduleFrom

	for _, key := range os.Environ() {
		if strings.HasPrefix(key, scheduleKeyPrefix) {
			value := strings.Split(key, "=")[1]
			scheduleValues := strings.Split(value, ",")
			if len(scheduleValues) != 2 {
				return nil, fmt.Errorf("expecting schedule value to be a comma separated scheduleId,name but got %s", value)
			}
			config.Schedules = appendSchedule(config.Schedules, scheduleValues[0], scheduleValues[1])
		}
	}

	if len(config.Schedules) == 0 {
		return nil, fmt.Errorf("expecting at least one schedule defined as an env var using prefix SCHEDULE")
	}

	return config, nil
}

func appendSchedule(schedules []Schedule, scheduleID string, groupName string) []Schedule {
	newScheduleList := make([]Schedule, len(schedules))
	updated := false

	for i, s := range schedules {
		if s.CurrentOnCallGroupName != groupName {
			newScheduleList[i] = s

			continue
		}

		updated = true

		newScheduleList[i] = Schedule{
			ScheduleIDs:            append(s.ScheduleIDs, scheduleID),
			CurrentOnCallGroupName: groupName,
		}
	}

	if !updated {
		newScheduleList = append(newScheduleList, Schedule{
			ScheduleIDs:            []string{scheduleID},
			CurrentOnCallGroupName: groupName,
		})
	}

	return newScheduleList
}

func getPagerdutyScheduleTime(pdScheduleKey string) (time.Duration, error) {
	result := time.Second

	pdSchedule, ok := os.LookupEnv(pdScheduleKey)
	if !ok {
		return result, nil
	}

	v, err := time.ParseDuration(pdSchedule)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s as time.Duration: %w", pdSchedule, err)
	}

	return v, nil
}
