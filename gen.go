package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	cmdGen = "gen"

	format = "20060102"
)

var weekDays = [...]string{"日", "月", "火", "水", "木", "金", "土"}

func gen(dir string) error {
	now, err := time.Parse(format, time.Now().Format(format))
	if err != nil {
		return errInternal{Err: err}
	}

	l, err := calcLog(dir, now)
	if err != nil {
		return errInternal{Err: err}
	}

	t, err := template.New("index").Parse(indexTempl)
	if err != nil {
		return errInternal{Err: err}
	}

	args := map[string]interface{}{
		"header":   l.header(),
		"body":     l.body(),
		"weekDays": weekDays,
	}

	if err := t.Execute(os.Stdout, args); err != nil {
		return errInternal{Err: err}
	}

	return nil
}

type log map[string]*day

func calcLog(dir string, today time.Time) (log, error) {
	l := newLog(today)

	if err := filepath.Walk(dir, l.checkLearned); err != nil {
		return nil, errInternal{Err: err}
	}

	return l, nil
}

func newLog(today time.Time) log {
	l := log{}

	for diff, showDays := 0, 7*52+int(today.Weekday()); diff <= showDays; diff++ {
		yyyymmdd := today.AddDate(0, 0, -diff).Format(format)
		l[yyyymmdd] = &day{yyyymmdd: yyyymmdd, index: showDays - diff + 1}
	}

	return l
}

func (l log) checkLearned(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	yyyymmdd := strings.TrimSuffix(info.Name(), ".md")

	if _, err := time.Parse(format, yyyymmdd); err != nil {
		return nil
	}

	if _, ok := l[yyyymmdd]; ok {
		l[yyyymmdd].Learned = true
	}

	return nil
}

func (l log) body() [][]*day {
	m := make([][]*day, 7)

	for i := 0; i < 7; i++ {
		m[i] = make([]*day, (len(l)/7 + 1))
	}

	for _, d := range l {
		m[(d.index+1)%7][(d.index+1)/7] = d
	}

	return m
}

func (l log) header() []string {
	h := make([]string, 0, (len(l)/7 + 1))
	for i := 1; i <= (len(l)/7 + 1); i++ {
		h = append(h, fmt.Sprintf("%02d", i))
	}

	return h
}

type day struct {
	yyyymmdd string
	index    int
	Learned  bool
}
