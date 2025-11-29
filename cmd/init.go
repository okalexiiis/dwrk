package cmd

import (
	"github.com/okalexiiis/dwrk/cmd/list"
	"github.com/okalexiiis/dwrk/cmd/new"
	"github.com/okalexiiis/dwrk/cmd/open"
)

func init() {
	RootCmd.AddCommand(list.ListCmd)
	RootCmd.AddCommand(new.NewCmd)
	RootCmd.AddCommand(open.OpenCmd)
}
