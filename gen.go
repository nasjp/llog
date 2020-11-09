package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const cmdGen = "gen"

func gen(dir string) error {
	l, err := calc(dir)
	if err != nil {
		return err
	}

	t, err := template.New("index").Funcs(map[string]interface{}{"raw": func(text string) template.HTML { return template.HTML(text) }}).Parse(indexTempl)
	if err != nil {
		return errInternal{Err: err}
	}

	if err := t.Execute(os.Stdout, map[string]string{"table": l.html()}); err != nil {
		return errInternal{Err: err}
	}

	return err
}

func calc(dir string) (log, error) {
	l := initLog(time.Now())

	if err := filepath.Walk(dir, l.calc); err != nil {
		return nil, err
	}

	return l, nil
}

type log map[string]bool

const format = "20060102"

func initLog(today time.Time) log {
	days := 7*52 + int(today.Weekday())

	l := log{}

	for diff := 0; diff <= days; diff++ {
		l[today.AddDate(0, 0, -diff).Format(format)] = false
	}

	return l
}

func (l log) calc(path string, info os.FileInfo, err error) error {
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

var weekDays = [...]string{"日", "月", "火", "水", "木", "金", "土"}

func (l log) html() string {
	s := make([]*day, 0, len(l))
	for d, learned := range l {
		s = append(s, &day{d: d, learned: learned})
	}

	sort.Slice(s, func(i, j int) bool {
		return s[i].d < s[j].d
	})

	t := map[int][]*day{}

	for i, d := range s {
		t[i/7] = append(t[i/7], d)
	}

	transposed := make([][]*day, len(t))

	for w, days := range t {
		transposed[w] = days
	}

	html := `<table border="1" width="400", height="180" style="font-size: x-small">
  <thead>
    <tr>
      <th> </th>
`

	for i := 0; i < len(s)/7+1; i++ {
		html += fmt.Sprintf(`      <th>%02d</th>
`, i+1)
	}

	html += `    </tr>
  </thead>
  <tbody>
`

	for k := 0; k < 7; k++ {
		html += `    <tr>
`

		html += fmt.Sprintf(`      <td>%s</td>
`, weekDays[k])

		for j := 0; j < len(s)/7+1; j++ {
			days := transposed[j]
			if len(days) <= k {
				html += `      <td></td>
`

				continue
			}

			html += days[k].tableData()
		}

		html += `    </tr>
`
	}

	return html + `
  </tbody>
</table>`
}

type day struct {
	d       string
	learned bool
}

func (d *day) tableData() string {
	if d.learned {
		return `      <td style="background-color: lime;"> </td>
`
	}

	return `      <td> </td>
`
}
