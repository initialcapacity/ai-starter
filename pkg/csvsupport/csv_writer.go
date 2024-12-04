package csvsupport

import (
	"encoding/csv"
	"io"
	"log/slog"
)

func WriteCSV(file io.Writer, lines [][]string) error {
	w := csv.NewWriter(file)
	defer w.Flush()

	for _, line := range lines {
		err := w.Write(line)
		if err != nil {
			slog.Error("unable to write CSV", "error", err)
			return err
		}
	}

	return nil
}
