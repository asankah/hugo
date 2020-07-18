package pandoc_config

// Config contains configuration settings for Pandoc.
type Config struct {
	// Input format. Use the 'Extensions' field to specify extensions thereof.
	// Only specify the bare format here. Defaults to 'markdown' if empty. Invoke
	// "pandoc --list-input-formats" to see the list of supported input formats
	// including various Markdown dialects.
	InputFormat string

	// If true, the output format is HTML (i.e. "--to=html"). Otherwise the output
	// format is HTML5 (i.e. "--to=html5").
	UseLegacyHtml bool

	// Equivalent to specifying "--mathjax". For compatibility, this option is
	// always true if none of the other math options are used.
	// See https://pandoc.org/MANUAL.html#math-rendering-in-html
	UseMathjax bool

	// Equivalent to specifying "--mathml".
	// See https://pandoc.org/MANUAL.html#math-rendering-in-html
	UseMathml bool

	// Equivalent to specifying "--webtex".
	// See https://pandoc.org/MANUAL.html#math-rendering-in-html. Uses the default
	// Webtex rendering URL.
	UseWebtex bool

	// Equivalent to specifying "--katex".
	// See https://pandoc.org/MANUAL.html#math-rendering-in-html
	UseKatex bool

	// List of filters to use. These translate to '--filter=' arguments to the
	// pandoc invocation.  The order of elements in `Filters` is preserved when
	// constructing the `pandoc` commandline.
	Filters []string

	// List of Pandoc Markdown extensions to use. No need to include default
	// extensions. Specifying ["foo", "bar"] is equivalent to specifying
	// --from=markdown+foo+bar on the pandoc commandline.
	Extensions []string

	// Random extra arguments.
	ExtraArgs []string
}
