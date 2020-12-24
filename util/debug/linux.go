// +build linux

package debug

import (
	"debug/elf"
	"debug/gosym"
)

func BuildLineTable(bin string) (*gosym.Table, error) {
	f, err := elf.Open(bin)
	if err != nil {
		return nil, err
	}

	sym := f.Section(".gosym")

	pcln := f.Section(".gopclntab")
	pclntab := gosym.NewLineTable(pcln.Data(), pcln.Addr)

	return gosym.NewTable(sym.Data(), pclntab)
}
