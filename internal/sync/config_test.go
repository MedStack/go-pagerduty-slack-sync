package sync

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewConfigFromEnv_SingleScheduleDefined(t *testing.T) {
	defer SetEnv("SCHEDULE_TEST", "1234,test")()
	defer SetEnv("PAGERDUTY_TOKEN", "token1")()
	defer SetEnv("SLACK_TOKEN", "secretToken1")()
	defer SetEnv("RUN_INTERVAL_SECONDS", "10")()

	config, err := NewConfigFromEnv()

	assert.NoError(t, err)
	assert.Equal(t, "token1", config.PagerDutyToken)
	assert.Equal(t, "secretToken1", config.SlackToken)
	assert.Equal(t, 10, config.RunIntervalInSeconds)
	assert.Equal(t, time.Second, config.PagerdutyScheduleFrom)
	assert.Equal(t, 1, len(config.Schedules))

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{{
		ScheduleIDs:            []string{"1234"},
		CurrentOnCallGroupName: "test",
	}},
		config.Schedules))
}

func Test_NewConfigFromEnv_SingleScheduleDefinedWithDefaultRunInterval(t *testing.T) {
	defer SetEnv("SCHEDULE_ON_CALL", "1234,on-call")()
	defer SetEnv("PAGERDUTY_TOKEN", "token1")()
	defer SetEnv("SLACK_TOKEN", "secretToken1")()

	config, err := NewConfigFromEnv()

	assert.NoError(t, err)
	assert.Equal(t, "token1", config.PagerDutyToken)
	assert.Equal(t, "secretToken1", config.SlackToken)
	assert.Equal(t, 60, config.RunIntervalInSeconds)
	assert.Equal(t, 1, len(config.Schedules))
	assert.Equal(t, "on-call", config.Schedules[0].CurrentOnCallGroupName)

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{{
		ScheduleIDs:            []string{"1234"},
		CurrentOnCallGroupName: "on-call",
	}},
		config.Schedules))
}

func Test_NewConfigFromEnv_SingleScheduleDefinedWithScheduleFrom(t *testing.T) {
	defer SetEnv("SCHEDULE_TEST", "1234,on-call")()
	defer SetEnv("PAGERDUTY_TOKEN", "token1")()
	defer SetEnv("SLACK_TOKEN", "secretToken1")()

	config, err := NewConfigFromEnv()

	assert.NoError(t, err)
	assert.Equal(t, "token1", config.PagerDutyToken)
	assert.Equal(t, "secretToken1", config.SlackToken)
	assert.Equal(t, 60, config.RunIntervalInSeconds)
	assert.Equal(t, time.Second, config.PagerdutyScheduleFrom)
	assert.Equal(t, 1, len(config.Schedules))
	assert.Equal(t, "on-call", config.Schedules[0].CurrentOnCallGroupName)

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{{
		ScheduleIDs:            []string{"1234"},
		CurrentOnCallGroupName: "on-call",
	}},
		config.Schedules))
}

func Test_NewConfigFromEnv_NoSchedulesDefined(t *testing.T) {
	config, err := NewConfigFromEnv()

	assert.Errorf(t, err, "expecting at least one schedule defined as an env var using prefix SCHEDULE")
	assert.Nil(t, config)
}

func SetEnv(key, value string) func() {
	_ = os.Setenv(key, value)
	return func() {
		_ = os.Unsetenv(key)
	}
}
