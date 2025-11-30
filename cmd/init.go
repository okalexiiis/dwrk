package cmd

import (
	"github.com/okalexiiis/dwrk/cmd/clone"
	"github.com/okalexiiis/dwrk/cmd/config"
	"github.com/okalexiiis/dwrk/cmd/list"
	"github.com/okalexiiis/dwrk/cmd/new"
	"github.com/okalexiiis/dwrk/cmd/open"
)

func init() {
	RootCmd.AddCommand(list.ListCmd)
	RootCmd.AddCommand(new.NewCmd)
	RootCmd.AddCommand(open.OpenCmd)
	RootCmd.AddCommand(clone.CloneCmd)
	RootCmd.AddCommand(config.ConfigCmd)
}
