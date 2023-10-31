package sync

import (
	"time"

	"github.com/medstack/go-pagerduty-slack-sync/internal/compare"
	"github.com/sirupsen/logrus"
)

// Schedules does the sync
func Schedules(config *Config) error {
	logrus.Infof("running schedule sync...")
	slackClient, err := newSlackClient(config.SlackToken)
	if err != nil {
		return err
	}
	slackBotClient, err := newSlackClient(config.SlackBotToken)
	if err != nil {
		return err
	}
	p := newPagerDutyClient(config.PagerDutyToken)

	updateSlackGroup := func(emails []string, groupName string) error {
		slackIDs, err := slackClient.getSlackIDsFromEmails(emails)
		if err != nil {
			return err
		}

		userGroup, err := slackClient.createOrGetUserGroup(groupName)
		if err != nil {
			return err
		}
		members, err := slackClient.Client.GetUserGroupMembers(userGroup.ID)
		if err != nil {
			return err
		}

		if !compare.Array(slackIDs, members) {
			logrus.Infof("member list %s needs updating...", groupName)
			_, err = slackClient.Client.UpdateUserGroupMembers(userGroup.ID, slackIDs[0])
			if err != nil {
				return err
			}

			err = slackBotClient.postMessage("support", slackIDs[0], userGroup.ID)
			if err != nil {
				return err
			}
		}
		return nil
	}

	getEmailsForSchedules := func(schedules []string, from time.Duration, to time.Duration) ([]string, error) {
		var emails []string

		for _, sid := range schedules {
			e, err := p.getEmailsForSchedule(sid, from, to)
			if err != nil {
				return nil, err
			}

			emails = appendIfMissing(emails, e...)
		}

		return emails, nil
	}

	for _, schedule := range config.Schedules {
		logrus.Infof("checking slack group: %s", schedule.CurrentOnCallGroupName)

		currentOncallEngineerEmails, err := getEmailsForSchedules(schedule.ScheduleIDs, config.PagerdutyScheduleFrom, time.Second)
		logrus.Infof("current on-call is %s...", currentOncallEngineerEmails)
		if err != nil {
			logrus.Errorf("failed to get emails for %s: %v", schedule.CurrentOnCallGroupName, err)
			continue
		}

		err = updateSlackGroup(currentOncallEngineerEmails, schedule.CurrentOnCallGroupName)
		if err != nil {
			logrus.Errorf("failed to update slack group %s: %v", schedule.CurrentOnCallGroupName, err)
			continue
		}
	}

	return nil
}

func appendIfMissing(slice []string, items ...string) []string {
out:
	for _, i := range items {
		for _, ele := range slice {
			if ele == i {
				continue out
			}
		}
		slice = append(slice, i)
	}

	return slice
}
