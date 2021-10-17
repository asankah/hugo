// Copyright 2019 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package pandoc converts content to HTML using Pandoc as an external helper.
package pandoc

import (
	"strings"

	"github.com/cli/safeexec"
	"github.com/gohugoio/hugo/htesting"
	"github.com/mitchellh/mapstructure"

	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/markup/bibliography"
	"github.com/gohugoio/hugo/markup/internal"
	"github.com/gohugoio/hugo/markup/pandoc/pandoc_config"

	"github.com/gohugoio/hugo/markup/converter"

	"fmt"
	"os"
	"path"
)

type paramer interface {
	Param(interface{}) (interface{}, error)
}

type searchPaths struct {
	Paths []string
}

func (s *searchPaths) NormalizePath(in_path string) (string, error) {
	if path.IsAbs(in_path) {
		return in_path, nil
	}

	for _, p := range s.Paths {
		fp := path.Join(p, in_path)
		if _, err := os.Stat(fp); err == nil {
			return fp, nil
		}
	}
	return "", fmt.Errorf("Can't find %s", in_path)
}

func (s *searchPaths) AsResourcePath() string {
	return strings.Join(s.Paths, ":")
}

// Provider is the package entry point.
var Provider converter.ProviderProvider = provider{}

type provider struct {
}

func (p provider) New(cfg converter.ProviderConfig) (converter.Provider, error) {
	return converter.NewProvider("pandoc", func(ctx converter.DocumentContext) (converter.Converter, error) {
		return &pandocConverter{
			docCtx: ctx,
			cfg:    cfg,
		}, nil
	}), nil
}

type pandocConverter struct {
	docCtx converter.DocumentContext
	cfg    converter.ProviderConfig
}

func (c *pandocConverter) Convert(ctx converter.RenderContext) (converter.Result, error) {
	return converter.Bytes(c.getPandocContent(ctx.Src)), nil
}

func (c *pandocConverter) Supports(feature identity.Identity) bool {
	return false
}

// getPandocContent calls pandoc as an external helper to convert pandoc markdown to HTML.
func (c *pandocConverter) getPandocContent(src []byte) []byte {
	logger := c.cfg.Logger
	pandoc_path := getPandocExecPath()
	if pandoc_path == "" {
		logger.Println("pandoc not found in $PATH: Please install.\n",
			"                 Leaving pandoc content unrendered.")
		return src
	}

	searchPathSet := searchPaths{
		Paths: []string{path.Dir(c.docCtx.Filename), "static", "."},
	}

	var pandocConfig pandoc_config.Config = c.cfg.MarkupConfig.Pandoc
	var bibConfig bibliography.Config = c.cfg.MarkupConfig.Bibliography

	if pageParameters, ok := c.docCtx.Document.(paramer); ok {
		if bibParam, err := pageParameters.Param("bibliography"); err == nil {
			mapstructure.WeakDecode(bibParam, &bibConfig)
		}

		if pandocParam, err := pageParameters.Param("pandoc"); err == nil {
			mapstructure.WeakDecode(pandocParam, &pandocConfig)
		}
	}

	arguments := pandocConfig.AsPandocArguments(&searchPathSet)

	if bibConfig.Source != "" {
		sourcePath, err := searchPathSet.NormalizePath(bibConfig.Source)
		if err != nil {
			logger.Errorf("Can't find bibliography: %s", bibConfig.Source)
		} else {
			arguments = append(arguments, "--bibliography", sourcePath)
		}
	}

	if bibConfig.CitationStyle != "" {
		citationPath, err := searchPathSet.NormalizePath(bibConfig.CitationStyle)
		if err != nil {
			logger.Errorf("Can't find citation style: %s", bibConfig.CitationStyle)
		} else {
			arguments = append(arguments, "--csl", citationPath)
		}
	}

	arguments = append(arguments, "--resource-path", searchPathSet.AsResourcePath())

	return internal.ExternallyRenderContent(c.cfg, c.docCtx, src, pandoc_path, arguments)
}

func getPandocExecPath() string {
	path, err := safeexec.LookPath("pandoc")
	if err != nil {
		return ""
	}

	return path
}

// Supports returns whether Pandoc is installed on this computer.
func Supports() bool {
	if htesting.SupportsAll() {
		return true
	}
	return getPandocExecPath() != ""
}
