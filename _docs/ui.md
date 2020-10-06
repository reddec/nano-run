# UI

UI available by default at port 8989, browser will be automatically redirected to `/ui/`.

Disable ui by `disable_ui: yes` flag in the config.

If UI directory (from `ui_directory`, default `ui`), embedded UI will be used. 

UI directory supports static files under `static` directory under `/static` prefix.

All *.html files scanned under `ui_directory` directory as [Go HTML template](https://golang.org/pkg/html/template/)
with additional functions from [Sprig](http://masterminds.github.io/sprig/). Recursive scanning not supported.

UI embedded by [go-bindata](https://github.com/go-bindata/go-bindata).