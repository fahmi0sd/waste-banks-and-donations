package scheduler

import (
	"log/slog"

	backupService "github.com/fahmi0sd/waste-banks-and-donations/service/backup"
	"github.com/robfig/cron/v3"
)

func StartBackupScheduler(
	logger *slog.Logger,
	service backupService.Service,
) {

	c := cron.New()

	_, err := c.AddFunc(
		"0 2 * * *",
		func() {

			logger.Info("running scheduled backup")

			if err := service.TriggerScheduler(); err != nil {

				logger.Error(
					"scheduled backup failed",
					"error", err,
				)

			}

		},
	)

	if err != nil {

		logger.Error(
			"failed register scheduler",
			"error", err,
		)

		return

	}

	c.Start()

	logger.Info("backup scheduler started")
}
