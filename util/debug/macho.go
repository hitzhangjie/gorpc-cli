// +build darwin

package debug

import (
	"debug/gosym"
	"debug/macho"
)

func BuildLineTable(bin string) (*gosym.Table, error) {
	f, err := macho.Open(bin)
	if err != nil {
		return nil, err
	}

	pcln := f.Section("__gopclntab")
	dat, err := pcln.Data()
	if err != nil {
		return nil, err
	}
	pclntab := gosym.NewLineTable(dat, pcln.Addr)

	sym := f.Section("__gosymtab")
	dat, err = sym.Data()
	if err != nil {
		return nil, err
	}
	return gosym.NewTable(dat, pclntab)
}
