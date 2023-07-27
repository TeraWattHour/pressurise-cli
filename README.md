# Pressurise CLI

This repository contains the CLI for generating Pressurise projects.
The client module is available [here](https://github.com/terawatthour/pressurise).

This project is in early stages of development, any input on how to make it better is very much welcome.

## Concept

The main goal of this library is to simplify development of full-stack apps
in Go. To accomplish this we use the powerful `html/template` package together
with some simple code generation. Really, this native library is all you need
in most of the cases. Not many languages offer such a great tool, do they?

How your project looks is very similar to Astro javascript apps - but in Go.
Every `.html` file in the `app/` directory is mapped to its generated route - that's all what happens under the hood.

## Usage

1. `go install github.com/terawatthour/pressurise-cli`
2. `pressurise-cli build .`

## Example

- app/layout.htm - bare in mind that files that are not routes may not
  contain code blocks

  ```html
  <!DOCTYPE html>
  <html lang="en">
    <head>
      <meta charset="UTF-8" />
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <title>Document</title>
    </head>
    <body>
      {{ block "content" . }} fallback {{ end }}
    </body>
  </html>
  ```

- app/index.html - this route will map to url `/` and will be an extension of
  `./layout.htm` HTML template

  ```
  ---
  // extends command takes in one argument which is a relative path
  // to the extended template, this file may not be a route
  !extends ./layout.htm
  !component ./footer.htm

  import (
    "fmt"
  )

  // this is executed on every request,
  // in this code block w (http.ResponseWriter),
  // r (*http.Request) and every structure
  // declared in the package `main` are available

  fmt.Println("hello from my page", r.Method)

  // `method` is available in the template whilst r.Method is not,
  // in short: you can use every structure declared in this code block
  // in the HTML template below
  method := r.Method

  ---

  {{ define "content" }}
    <h1>This was accessed with method: {{ .method }} </h1>

    {{ template "footer" . }}
  {{ end }}

  ```

## TODOs

- Add tests (who are we kidding I'm not doing that)
- Add docs
