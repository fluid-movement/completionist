package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	FolderName = "completionist"
	FileName = "todos.json"
)

type Storage interface {
	Load() (*TodoList, error)
	Save(list *TodoList) error
}

type JSONStorage struct {
	Path string
}

func (s *JSONStorage) Load() (*TodoList, error) {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return &TodoList{}, nil
		}
		return nil, err
	}

	var list TodoList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}

	return &list, nil
}

func (s *JSONStorage) Save(list *TodoList) error {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.Path, data, 0644)
}

func InitializeStorage() (Storage, error) {
	path, err := storagePath()
	if err != nil {
		return nil, err
	}

	return &JSONStorage{Path: path}, nil
}

func storagePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, FolderName, FileName), nil
}
