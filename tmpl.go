package main

const indexTempl = `<!DOCTYPE html>
<html>
  <head>
  <meta charset="utf-8">
  <title>Learning Log</title>
  </head>
<body>
  <h1>Learning Log</h1>
	<table border="1" width="400", height="180" style="font-size: x-small">
  <thead>
    <tr>
      <th> </th>
      {{range .header}}
      <th>{{.}}</th>
      {{end}}
    </tr>
  </thead>
  <tbody>
    {{ $wd := .weekDays }}
    {{range $i, $v := .body}}
    <tr>
      <td>{{index $wd $i}}</td>
      {{range $j, $w := $v}}
      <td {{with $w }} {{if $w.Learned}} style="background-color: lime;" {{end}} {{end}}> </td>
      {{end}}
    </tr>
    {{end}}
  </tbody>
</body>
</html>`
