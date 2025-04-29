package httpserver

import "strings"

var (
	SwaggerUI   = &swaggerUIBuilder{}
	RapidocUI   = &rapidocBuilder{}
	StoplightUI = &stoplightElementBuilder{}
	RedocUI     = &redocBuilder{}
)

type OpenAPIUIBuilder interface {
	HTML(doc string, title string) string
	Doc() string
}

type swaggerUIBuilder struct {
}

func (s *swaggerUIBuilder) HTML(doc string, title string) string {
	const template = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="description" content="SwaggerUI" />
    <title>{:title}</title>
    <link
      rel="stylesheet"
      href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css"
    />
  </head>
  <style>
    * {
      font-family: Kaiti SC, cursive, sans-serifcursive, sans-serif !important;
    }
  </style>
  <body>
    <div id="swagger-ui"></div>
    <script
      src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"
      crossorigin
    ></script>
    <script
      src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-standalone-preset.js"
      crossorigin
    ></script>
    <script>
      window.onload = () => {
        window.ui = SwaggerUIBundle({
          spec: {:spec},
          dom_id: "#swagger-ui",
          filter: false,
          defaultModelsExpandDepth: -1,
        });
      };
    </script>
  </body>
</html>
`

	html := strings.ReplaceAll(template, "{:title}", title)
	html = strings.ReplaceAll(html, "{:spec}", doc)
	return html
}

func (s *swaggerUIBuilder) Doc() string {
	return "https://swagger.io/docs/open-source-tools/swagger-ui/usage/installation/#unpkg"
}

type rapidocBuilder struct {
}

func (r *rapidocBuilder) HTML(doc string, title string) string {
	const template = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <!-- Important: rapi-doc uses utf8 characters -->
    <title>{:title}</title>
    <script
      type="module"
      src="http://unpkg.com/rapidoc@9.3.8/dist/rapidoc-min.js"
    ></script>
  </head>
  <style>
    * {
      font-family: Kaiti SC, cursive, sans-serifcursive, sans-serif !important;
    }
  </style>
  <body>
    <rapi-doc id="thedoc"> </rapi-doc>

    <script>
      document.addEventListener("DOMContentLoaded", (event) => {
        let docEl = document.getElementById("thedoc");
        docEl.setAttribute("text-color", "");
        docEl.setAttribute("render-style", "read");
        docEl.setAttribute("show-header", "false");
        docEl.setAttribute("theme", "light");
        docEl.setAttribute("bg-color", "#f9f9fa");
        docEl.setAttribute("nav-bg-color", "#3f4d67");
        docEl.setAttribute("nav-text-color", "#a9b7d0");
        docEl.setAttribute("nav-hover-bg-color", "#333f54");
        docEl.setAttribute("nav-hover-text-color", "#fff");
        docEl.setAttribute("nav-accent-color", "#f87070");
        docEl.setAttribute("primary-color", "#5c7096");
        const spec = {:spec};
        docEl.loadSpec(spec);
      });
    </script>
  </body>
</html>
`

	html := strings.ReplaceAll(template, "{:title}", title)
	html = strings.ReplaceAll(html, "{:spec}", doc)
	return html
}

func (r *rapidocBuilder) Doc() string {
	return "https://rapidocweb.com/examples.html"
}

type stoplightElementBuilder struct {
}

func (s *stoplightElementBuilder) HTML(doc string, title string) string {
	const template = `
  <!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1, shrink-to-fit=no"
    />
    <title>{:title} Document</title>

    <script src="https://unpkg.com/@stoplight/elements/web-components.min.js"></script>
    <link
      rel="stylesheet"
      href="https://unpkg.com/@stoplight/elements/styles.min.css"
    />
  </head>
  <body>
    <div class="api-container">
      <elements-api
        id="docs"
        router="hash"
        layout="responsive"
        hideSchemas="true"
      ></elements-api>
    </div>
  </body>
  <style>
    body {
      display: flex;
      flex-direction: column;
      height: 100vh;
    }

    * {
      font-family: Kaiti SC, cursive, sans-serifcursive, sans-serif !important;
    }
    .api-container {
      flex: 1 0 0;
      overflow: hidden;
    }
  </style>
  <script>
    (async () => {
      const docs = document.getElementById("docs");
      docs.apiDescriptionDocument = {:spec};
    })();
  </script>
</html>`

	html := strings.ReplaceAll(template, "{:title}", title)
	html = strings.ReplaceAll(html, "{:spec}", doc)
	return html
}

func (s *stoplightElementBuilder) Doc() string {
	return "https://docs.stoplight.io/docs/elements/a71d7fcfefcd6-elements-in-html"
}

type redocBuilder struct {
}

func (r *redocBuilder) HTML(doc string, title string) string {
	const template = `<!DOCTYPE html>
<html>
  <head>
    <title>{:title}</title>
    <!-- needed for adaptive design -->
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <link
      href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700"
      rel="stylesheet"
    />

    <!--
    Redoc doesn't change outer page styles
    -->
    <style>
      body {
        margin: 0;
        padding: 0;
      }
      * {
        font-family: Kaiti SC, cursive, sans-serifcursive, sans-serif !important;
      }
    </style>
  </head>
  <body>
    <redoc id="doc"></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
    <script>
      (async()=>{
        Redoc.init({:spec}, {}, document.getElementById('doc'))
      })()
    </script>
  </body>
</html>
`

	html := strings.ReplaceAll(template, "{:title}", title)
	html = strings.ReplaceAll(html, "{:spec}", doc)
	return html
}

func (r *redocBuilder) Doc() string {
	return "https://redocly.com/docs/redoc/deployment/html"
}
