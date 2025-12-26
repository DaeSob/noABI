package logger

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	timeFormat      = "2006-01-02"
	defaultFileMode = fs.FileMode(0644)
	defaultDirMode  = fs.FileMode(0744)
	mb              = 1024 * 1024
)

type Config struct {
	FilePath     string
	MaxFileCount int
	MaxFileSize  int
}

var ConfigDefault = Config{
	MaxFileCount: 30,
	MaxFileSize:  30,
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}
	cfg := config[0]

	if cfg.MaxFileCount == 0 {
		cfg.MaxFileCount = ConfigDefault.MaxFileCount
	}
	if cfg.MaxFileSize == 0 {
		cfg.MaxFileSize = ConfigDefault.MaxFileSize
	}
	return cfg
}

var _ io.Writer = &LogWriter{}

type LogWriter struct {
	mu sync.Mutex

	fp    *os.File
	mode  fs.FileMode
	date  time.Time
	size  int64
	count int

	filePath string
	dir      string
	base     string
	ext      string
	prefix   string

	maxFileCount int
	maxFileSize  int64
}

func NewWriter(cfg ...Config) (io.Writer, error) {
	config := configDefault(cfg...)
	if config.FilePath == "" {
		return os.Stdout, nil
	} else {
		return New(config)
	}
}

func New(cfg Config) (*LogWriter, error) {
	if cfg.FilePath == "" {
		return nil, fmt.Errorf("empty file path")
	}
	cfg = configDefault(cfg)

	logger := &LogWriter{
		filePath:     cfg.FilePath,
		mode:         defaultFileMode,
		maxFileCount: cfg.MaxFileCount,
		maxFileSize:  int64(cfg.MaxFileSize * mb),
		dir:          filepath.Dir(cfg.FilePath),
		base:         filepath.Base(cfg.FilePath),
	}
	logger.ext = filepath.Ext(logger.base)
	logger.prefix = logger.base[:len(logger.base)-len(logger.ext)]

	if err := os.MkdirAll(logger.dir, defaultDirMode); err != nil {
		return nil, err
	}
	if err := logger.newFile(); err != nil {
		return nil, err
	}
	return logger, nil
}

func (w *LogWriter) Path() string { return w.filePath }

func (w *LogWriter) Write(p []byte) (n int, err error) {
	inputLen := int64(len(p))

	// ~~~ first round ~~~
	// pre-process
	switch {
	case w == nil, inputLen == 0: // skip write
		return 0, nil

	case inputLen > w.maxFileSize: // oversize log
		return 0, fmt.Errorf("log length : %d, max file size : %d", inputLen, w.maxFileSize)

	case w.fp == nil:
		if err = w.newFile(); err != nil {
			return 0, err
		}
	}

	// ~~~ second round ~~~
	// check-rotate
	w.mu.Lock()
	defer w.mu.Unlock()

	now := logDateYmdNow()
	bRotate := false
	if now.After(w.date) || inputLen+w.size > w.maxFileSize {
		bRotate = true
	}

	// ~~~ final round ~~~
	// rotate, wirte
	if bRotate {
		w.autoRemove()
		if err = w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = w.fp.Write(p)
	if err == nil {
		w.size += int64(n)
		w.date = now
	}
	return
}
