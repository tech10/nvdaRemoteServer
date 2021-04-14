package server

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func file_rewrite(file string, data []byte) error {
	var ferr error
	file, ferr = fileOps(file)
	if ferr != nil {
		return errors.New("Unable to create or open the file " + file + "\n" + ferr.Error())
	}
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

func fullPath(old_path string) string {
	var err error
	var path string
	path, err = filepath.Abs(old_path)
	if err != nil {
		return old_path
	}
	var e_path string
	e_path, err = filepath.EvalSymlinks(path)
	if err == nil {
		return e_path
	}
	e_path = ""
	n_path := ""
	err = nil
	for _, v := range strings.Split(path, PS) {
		e_path += v + PS
		if err != nil {
			continue
		}
		n_path, err = filepath.EvalSymlinks(e_path)
		if err == nil {
			e_path = n_path + PS
		}
	}
	return strings.TrimSuffix(e_path, PS)
}

func fileOps(file string) (string, error) {
	path := fullPath(file)
	return path, cdir(filepath.Dir(path))
}

func cdir(dir string) error {
	if !createDir {
		return nil
	}
	err := os.MkdirAll(dir, 0755)
	if err == nil {
		return nil
	}
	return errors.New("Unable to create directory " + dir + "\n" + err.Error())
}

func file_read(file string) ([]byte, error) {
	file = fullPath(file)
	return ioutil.ReadFile(file)
}
