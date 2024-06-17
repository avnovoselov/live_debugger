package configuration

import (
	"io/fs"

	"github.com/BurntSushi/toml"
)

func NewParser(fs fs.FS) *Parser {
	return &Parser{fs: fs}
}

type Parser struct {
	fs fs.FS
}

func (p *Parser) Parse(path string) (config Common, err error) {
	_, err = toml.DecodeFS(p.fs, path, &config)

	return
}
