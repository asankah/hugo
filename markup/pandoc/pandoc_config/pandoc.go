package pandoc_config

import (
	"fmt"
	"strings"
)

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

	// List of input format extensions to use. Specifying ["foo", "bar"] is
	// equivalent to specifying --from=markdown+foo+bar on the pandoc commandline
	// assuming InputFormat is "markdown".
	InputExtensions []string

	// List of output format extensions to use. Specifying ["foo", "bar"] is
	// equivalent to specifying --to=html5+foo+bar on the pandoc commandline,
	// assuming UseLegacyHTML is false. Invoke "pandoc --list-extensions=html5" to
	// or "pandoc --list-extensions=html5" to see the list of extensions that can
	// be specified here.
	OutputExtensions []string

	// Extra commandline options passed to the pandoc invocation. These options
	// are appended to the commandline after the format and filter options.
	// Arguments are passed in literally. Hence must have the "--" or "-" prefix
	// where applicable.
	ExtraArgs []string
}

func (c *Config) getInputArg() string {
	var b strings.Builder
	b.WriteString("--from=")
	if len(c.InputFormat) > 0 {
		b.WriteString(c.InputFormat)
	} else {
		b.WriteString("markdown")
	}

	for _, extension := range c.InputExtensions {
		b.WriteString("+")
		b.WriteString(extension)
	}
	return b.String()
}

func (c *Config) getOutputArg() string {
	var b strings.Builder
	b.WriteString("--to=")
	if c.UseLegacyHtml {
		b.WriteString("html")
	} else {
		b.WriteString("html5")
	}

	for _, extension := range c.OutputExtensions {
		b.WriteString("+")
		b.WriteString(extension)
	}
	return b.String()
}

func (c *Config) getMathRenderingArg() string {
	switch {
	case c.UseMathml:
		return "--mathml"
	case c.UseWebtex:
		return "--webtex"
	case c.UseKatex:
		return "--katex"
	default:
		return "--mathjax"
	}
}

func (c *Config) getFilterArgs() []string {
	var args []string
	for _, filter := range c.Filters {
		args = append(args, fmt.Sprintf("--filter=%s", filter))
	}
	return args
}

// AsPandocArguments returns a list of strings that can be used as arguments to
// a "pandoc" invocation. All the settings contained in Config are represented
// in the returned list of arguments.
func (c *Config) AsPandocArguments() []string {
	args := []string{
		c.getInputArg(),
		c.getOutputArg(),
		c.getMathRenderingArg()}

	args = append(args, c.getFilterArgs()...)
	args = append(args, c.ExtraArgs...)

	return args
}
