package printer

import (
	"io"

	"github.com/rollwagen/tf-module-versions/internal/tf"
)

// ModuleVersionPrinter is an interface that knows how to print Module information
type ModuleVersionPrinter interface {
	// PrintReport receives a report, formats it and prints it to a writer.
	PrintReport([]tf.Module, io.Writer) error
}
