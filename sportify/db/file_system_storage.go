package db

import (
	"context"
	"fmt"
	"os"
	"sync"
)

type FileSystemStorage struct {
	baseDir             string
	mapExistenceFiles   map[string]struct{}
	muMapExistenceFiles *sync.RWMutex
}

func NewFileSystemStorage(baseDir string) (*FileSystemStorage, error) {
	preFSStorage := &FileSystemStorage{
		baseDir:             baseDir,
		mapExistenceFiles:   make(map[string]struct{}),
		muMapExistenceFiles: &sync.RWMutex{},
	}

	err := preFSStorage.recover()
	if err != nil {
		return nil, fmt.Errorf("to recover: %w", err)
	}

	return preFSStorage, nil
}

func (f *FileSystemStorage) recover() error {
	files, err := os.ReadDir(f.baseDir)
	if err != nil {
		return fmt.Errorf("to read base dir: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			f.muMapExistenceFiles.Lock()
			f.mapExistenceFiles[file.Name()] = struct{}{}
			f.muMapExistenceFiles.Unlock()
		}
	}

	return nil
}

// Check bool in return slice means file exist if it's true.
func (f *FileSystemStorage) Check(_ context.Context, files []string) ([]bool, error) {
	result := make([]bool, len(files))

	f.muMapExistenceFiles.RLock()
	defer f.muMapExistenceFiles.RUnlock()

	for i, filename := range files {
		if _, ok := f.mapExistenceFiles[filename]; ok {
			result[i] = true
		} else {
			result[i] = false
		}

	}

	return result, nil
}

func (f *FileSystemStorage) SaveFile(_ context.Context, content []byte, fileName string) error {
	file, err := os.Create(f.baseDir + "/" + fileName)
	if err != nil {
		return fmt.Errorf("to create file: %w", err)
	}

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("to write file: %w", err)
	}

	f.muMapExistenceFiles.Lock()
	f.mapExistenceFiles[fileName] = struct{}{}
	f.muMapExistenceFiles.Unlock()

	return nil
}
