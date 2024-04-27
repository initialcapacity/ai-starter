package collector

import "log/slog"

func Collect() error {
	slog.Info("Starting to collect data")
	defer slog.Info("Finished collecting data")
	return nil
}
