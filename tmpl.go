package main

const indexTempl = `<!DOCTYPE html>
<html>
  <head>
  <meta charset="utf-8">
  <title>learning log</title>
  </head>
<body>
  <h1>Hello World</h1>

    {{.table | raw}}
</body>
</html>`
