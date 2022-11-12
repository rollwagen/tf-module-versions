package printer

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/rollwagen/tf-module-versions/internal/tf"
)

type JSONPrinter struct{}

func (JSONPrinter) PrintReport(modules []tf.Module, writer io.Writer) error {
	b, _ := json.MarshalIndent(modules, "", "  ")
	_, err := fmt.Fprintln(writer, string(b))

	return err
}
