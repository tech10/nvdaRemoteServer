package server

import (
	"errors"
	"os"
)

func file_rewrite(file string, data []byte) error {
	w, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.New("Unable to create or open the file " + file + "\n" + err.Error())
	}
	_, err = w.Write(data)
	if err != nil {
		return errors.New("Unable to write to the file " + file + "\n" + err.Error())
	}
	err = w.Close()
	if err != nil {
		return errors.New("The file at " + file + " was unable to close. Information may not have been written to it correctly.\n" + err.Error())
	}
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return !info.IsDir()
}
