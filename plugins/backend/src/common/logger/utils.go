package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func dateFormatYmd(y int, m time.Month, d int) (time.Time, string) {
	str := fmt.Sprintf("%d-%02d-%02d", y, m, d)
	date, _ := time.Parse(timeFormat, str)
	return date, str
}

func subDate(date time.Time) string {
	_, str := dateFormatYmd(date.AddDate(0, 0, -1).Date())
	return str
}

func logDateYmd(t time.Time) time.Time {
	t2, _ := dateFormatYmd(t.Date())
	return t2
}

func logDateYmdNow() time.Time {
	return logDateYmd(time.Now())
}

func (w *LogWriter) fileList() []string {
	var files []string
	matches, _ := filepath.Glob(w.dir + "/" + w.prefix + "*" + w.ext)
	for _, s := range matches {
		fi, _ := os.Stat(s)

		if !fi.IsDir() {
			files = append(files, s)
		}
	}
	return files
}

func (w *LogWriter) newFile() error {
	var err error
	w.fp, err = os.OpenFile(w.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, w.mode)
	if err != nil {
		return err
	}
	fi, err := w.fp.Stat()
	if err != nil {
		return err
	}
	w.size = fi.Size()
	w.count = len(w.fileList())
	w.date = logDateYmd(fi.ModTime())
	return nil
}

func (w *LogWriter) autoRemove() {
	fl := w.fileList()
	count := len(w.fileList())
	if w.maxFileCount == 0 || w.maxFileCount >= count {
		return
	}
	os.Remove(fl[0])
}

func (w *LogWriter) rotateFileName() string {
	t := logDateYmdNow()
	date := subDate(t)

	count := 0
	fl := w.fileList()
	for _, f := range fl {
		if strings.Contains(f, date) {
			count = count + 1
		}
	}
	return filepath.Join(
		w.dir,
		fmt.Sprintf("%s.%s.%d%s",
			w.prefix, date, count, w.ext,
		))
}

func (w *LogWriter) close() (err error) {
	switch w.fp {
	case nil:
	default:
		if err = w.fp.Close(); err == nil {
			w.fp = nil
		}
	}
	return
}

func (w *LogWriter) rotate() error {
	if _, err := os.Stat(w.filePath); err != nil {
		return err
	}
	w.close()

	if err := os.Rename(w.filePath, w.rotateFileName()); err != nil {
		return err
	}
	return w.newFile()
}
