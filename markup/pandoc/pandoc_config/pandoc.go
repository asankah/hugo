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

	// Equivalent to specifying "--katex".
	// See https://pandoc.org/MANUAL.html#math-rendering-in-html
	UseKatex bool

	// List of filters to use. These translate to '--filter=' or '--lua-filter'
	// arguments to the pandoc invocation.
	//
	// The order of elements in `Filters` is preserved when constructing the
	// `pandoc` commandline.
	//
	// When used with bibliography, the `--citeproc` argument is placed *after*
	// all the filters.
	//
	// Use the prefix 'lua:' or the suffix '.lua' to indicate Lua filters.
	Filters []string

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

	// Metadata. The dictionary keys and values are handled in the obvious way.
	Metadata map[string]interface{}

	// Resource paths in addition to the top of the project directory.
	ResourcePath []string
}

var Default = Config{
	InputFormat:   "markdown",
	UseLegacyHtml: false,
	UseMathjax:    false,
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
	if c.UseKatex {
		return "--katex"
	}
	return "--mathjax"
}

func (c *Config) getMetadataArgs() []string {
	var args []string
	for k, iv := range c.Metadata {
		var v string
		if sv, ok := iv.(string); ok {
			v = sv
		} else if sv, ok := iv.(fmt.Stringer); ok {
			v = sv.String()
		} else {
			v = fmt.Sprintf("%v", iv)
		}
		args = append(args, fmt.Sprintf("-M%s=%s", k, v))
	}
	return args
}

func (c *Config) getFilterArgs() []string {
	var args []string
	for _, filterPath := range c.Filters {
		if strings.HasPrefix(filterPath, "lua:") || strings.HasSuffix(filterPath, ".lua") {
			args = append(args, fmt.Sprintf("--lua-filter=%s", strings.TrimPrefix(filterPath, "lua:")))
		} else {
			args = append(args, fmt.Sprintf("--filter=%s", filterPath))
		}
	}
	return args
}

func (c *Config) getResourcePathArgs() []string {
	return []string{"--resource-path", strings.Join(c.ResourcePath, ":")}
}

// AsPandocArguments returns a list of strings that can be used as arguments to
// a "pandoc" invocation. All the settings contained in Config are represented
// in the returned list of arguments.
func (c *Config) AsPandocArguments() []string {
	args := []string{
		c.getInputArg(),
		c.getOutputArg(),
		c.getMathRenderingArg()}

	args = append(args, c.getMetadataArgs()...)
	args = append(args, c.getFilterArgs()...)
	args = append(args, c.getResourcePathArgs()...)

	return args
}
