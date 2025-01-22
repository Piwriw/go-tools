package chrono

import (
	"log/slog"

	"github.com/google/uuid"
)

var (
	defaultBeforeJobRuns = func(jobID uuid.UUID, jobName string) {
		slog.Info("BeforeJobRuns", "jobID", jobID, "jobName", jobName)
	}
	defaultBeforeJobRunsSkipIfBeforeFuncErrors = func(jobID uuid.UUID, jobName string) error {
		slog.Info("BeforeJobRunsSkipIfBeforeFuncErrors", "jobID", jobID, "jobName", jobName)
		return nil
	}
	defaultAfterJobRuns = func(jobID uuid.UUID, jobName string) {
		slog.Info("BeforeJobRunsAfterFuncErrors", "jobID", jobID, "jobName", jobName)
	}
	defaultAfterJobRunsWithError = func(jobID uuid.UUID, jobName string, err error) {
		slog.Info("AfterJobRunsWithError", "jobID", jobID, "jobName", jobName)
	}
	defaultAfterJobRunsWithPanic = func(jobID uuid.UUID, jobName string, recoverData any) {
		slog.Info("AfterJobRunsWithPanic", "jobID", jobID, "jobName", jobName)
	}
	defaultAfterLockError = func(jobID uuid.UUID, jobName string, err error) {
		slog.Info("AfterLockError", "jobID", jobID, "jobName", jobName)
	}
)
