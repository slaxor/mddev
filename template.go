package main

const PAGE = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>{{.Filename}}</title>
    <link rel="stylesheet" type="text/css" href="/app.css" />

  </head>
  <body>
    <div id="markdown">
      {{.ParsedContent}}
    </div>
  </body>
  <script src="/app.js"></script>
</html>
`
