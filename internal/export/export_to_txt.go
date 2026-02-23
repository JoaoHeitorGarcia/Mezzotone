package export

import (
	"bufio"
	"os"
)

func ASCIItToTxT(outPath string, exportString string) error {
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(exportString)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	err = w.Flush()
	if err != nil {
		return err
	}
	return nil
}
