package printer

import (
	"fmt"
	"io"

	"github.com/rollwagen/tf-module-versions/internal/tf"
)

type ErrorformatPrinter struct{}

func (ErrorformatPrinter) PrintReport(modules []tf.Module, writer io.Writer) error {
	// -efm="%f:%l: %m"
	//   %f	file name
	//   %l	line number
	//   %m	error message
	for _, module := range modules {
		if !module.HasNewerVersion() {
			continue
		}

		f := module.Location.FileName
		l := module.Location.Line

		var m string
		if module.UsedVersion == "nil" {
			m = fmt.Sprintf(":warning: No version :bangbang: Latest version for %s is %s", module.Name, module.AvailableVersion)
		} else {
			m = fmt.Sprintf("Newer version %s available for module %s", module.AvailableVersion, module.Name)
		}

		_, err := fmt.Fprintf(writer, "%s:%d: %s\n", f, l, m)
		if err != nil {
			return err
		}
	}

	return nil
}
