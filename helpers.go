package larago

import "os"

func (l *Larago) createDirIfNotExist(path string) error {
	const mode = 0755
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, mode)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Larago) createFileIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		var file, err = os.Create(path)

		if err != nil {
			return err
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}
	return nil
}
