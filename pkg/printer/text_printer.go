package printer

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rollwagen/tf-module-versions/internal/tf"
)

type TextPrinter struct{}

func (TextPrinter) PrintReport(modules []tf.Module, writer io.Writer) error {
	if len(modules) == 0 {
		fmt.Println("No Gitlab terraform modules used. Nothing to validate.")
		return nil
	}

	red := func(s string) string {
		return color.New(color.FgRed).Sprintf("%s", s)
	}
	magenta := func(s string) string {
		return color.New(color.FgHiMagenta).Sprintf("%s", s)
	}
	yellow := func(s string) string {
		return color.New(color.FgYellow).Sprintf("%s", s)
	}
	green := func(s string) string {
		return color.New(color.FgHiGreen).Sprintf("%s", s)
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleDefault)
	t.SetOutputMirror(writer)
	t.AppendHeader(table.Row{"Name", "File", "Line", "Version used", "Version available", "Status"})

	for _, m := range modules {
		status := "?"

		if m.HasNewerVersion() {
			status = magenta("⚠")
			switch m.NewerVersion() {
			case "MINOR":
				status = yellow("⚠")
			case "MAJOR":
				status = red("⚠")
			}
		}

		if m.HasSameVersion() {
			status = green("✔")
		}

		used := m.UsedVersion
		if m.UsedVersion == "nil" {
			used = red("✖")
			status = red("")
		}

		t.AppendRow(table.Row{
			m.Name,
			m.Location.FileName,
			m.Location.Line,
			used,
			m.AvailableVersion,
			status,
		},
		)
	}
	t.Render()

	return nil
}
