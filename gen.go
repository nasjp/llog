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

	l, err := calc(dir, now)
	if err != nil {
		return errInternal{Err: err}
	}

	t, err := template.New("index").Parse(indexTempl)
	if err != nil {
		return errInternal{Err: err}
	}

	args := map[string]interface{}{
		"header":   l.header(),
		"body":     l.body(now),
		"weekDays": weekDays,
	}

	if err := t.Execute(os.Stdout, args); err != nil {
		return errInternal{Err: err}
	}

	return nil
}

func calc(dir string, today time.Time) (log, error) {
	l := initLog(today)

	if err := filepath.Walk(dir, l.checkLearned); err != nil {
		return nil, errInternal{Err: err}
	}

	return l, nil
}

type log map[string]bool

func initLog(today time.Time) log {
	l := log{}

	for diff := 0; diff <= 7*52+int(today.Weekday()); diff++ {
		l[today.AddDate(0, 0, -diff).Format(format)] = false
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

	d := strings.TrimSuffix(info.Name(), ".md")

	if _, err := time.Parse(format, d); err != nil {
		return nil
	}

	l[d] = true

	return nil
}

func (l log) body(today time.Time) [][]*day {
	m := make([][]*day, 7)

	for i := 0; i < 7; i++ {
		m[i] = make([]*day, (len(l)/7 + 1))
	}

	for d, learned := range l {
		t, _ := time.Parse(format, d)

		index := len(l) - int(today.Sub(t).Hours())/24
		column := (index + 1) / 7
		row := (index + 1) % 7
		m[row][column] = &day{d: d, Learned: learned}
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
	d       string
	Learned bool
}
