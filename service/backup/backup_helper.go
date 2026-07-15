package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

func (s *service) backupDirectory() (string, error) {

	dir := os.Getenv("BACKUP_DIR")

	if dir == "" {
		dir = "backup"
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return dir, nil
}

func (s *service) backupFilename() string {

	return fmt.Sprintf(
		"backup_%s.sql",
		time.Now().Format("20060102_150405"),
	)
}

func (s *service) buildCommand(fullPath string) *exec.Cmd {

	cmd := exec.Command(

		"pg_dump",

		"-h", os.Getenv("DB_HOST"),

		"-p", os.Getenv("DB_PORT"),

		"-U", os.Getenv("DB_USER"),

		"-d", os.Getenv("DB_NAME"),

		"-f", fullPath,
	)

	cmd.Env = append(
		os.Environ(),
		"PGPASSWORD="+os.Getenv("DB_PASS"),
	)

	return cmd
}

func (s *service) fileSize(path string) int64 {

	info, err := os.Stat(path)

	if err != nil {
		return 0
	}

	return info.Size()
}

func (s *service) backupKeep() int {

	keep := 30

	value := os.Getenv("BACKUP_KEEP")

	if value == "" {
		return keep
	}

	number, err := strconv.Atoi(value)

	if err != nil {
		return keep
	}

	return number
}

type backupFile struct {
	Name    string
	Path    string
	ModTime time.Time
}

func (s *service) scanBackupFiles(dir string) ([]backupFile, error) {

	entries, err := os.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	files := make([]backupFile, 0)

	for _, entry := range entries {

		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()

		if err != nil {
			continue
		}

		files = append(files, backupFile{
			Name:    entry.Name(),
			Path:    filepath.Join(dir, entry.Name()),
			ModTime: info.ModTime(),
		})
	}

	return files, nil
}

func (s *service) cleanupOldBackup(dir string) {

	files, err := s.scanBackupFiles(dir)

	if err != nil {
		return
	}

	keep := s.backupKeep()

	if len(files) <= keep {
		return
	}

	sort.Slice(files, func(i, j int) bool {

		return files[i].ModTime.Before(files[j].ModTime)

	})

	for i := 0; i < len(files)-keep; i++ {

		err := os.Remove(files[i].Path)

		if err != nil {

			s.logger.Error(
				"failed remove old backup",
				"file", files[i].Path,
				"error", err,
			)

			continue
		}

		s.logger.Info(
			"old backup removed",
			"file", files[i].Path,
		)

	}
}

func (s *service) executeBackup() (string, string, int64, error) {

	dir, err := s.backupDirectory()

	if err != nil {
		return "", "", 0, err
	}

	filename := s.backupFilename()

	fullPath := filepath.Join(
		dir,
		filename,
	)

	cmd := s.buildCommand(fullPath)

	if err := cmd.Run(); err != nil {
		return "", "", 0, err
	}

	size := s.fileSize(fullPath)

	s.cleanupOldBackup(dir)

	return filename, fullPath, size, nil
}
