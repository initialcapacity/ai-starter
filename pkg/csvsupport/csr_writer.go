package csvsupport

import (
	"encoding/csv"
	"io"
	"log/slog"
)

func WriteCSV(file io.Writer, header []string, lines chan []string) error {
	w := csv.NewWriter(file)
	defer w.Flush()

	err := w.Write(header)
	if err != nil {
		slog.Error("unable to write CSV", "error", err)
		return err
	}

	for line := range lines {
		lineErr := w.Write(line)
		if lineErr != nil {
			slog.Error("unable to write CSV", "error", lineErr)
			return lineErr
		}
	}

	return nil
}
