package main

import (
	douyinPlatform "github.com/jackwener/aiview/internal/platform/douyin"
	"github.com/jackwener/aiview/internal/platform"
	douyinCommands "github.com/jackwener/aiview/commands/douyin"
	"github.com/spf13/cobra"
)

func init() {
	p, ok := platform.GetPlatform("douyin")
	if !ok {
		return
	}
	dp := p.(*douyinPlatform.DouyinPlatform)

	douyinCmd := &cobra.Command{
		Use:   "douyin",
		Short: "Douyin platform commands",
		Long:  `Commands for interacting with Douyin (抖音) content.`,
	}

	// Hot search
	douyinCmd.AddCommand(douyinCommands.NewHotCmd(func() douyinCommands.Client { return dp.BuildClient() }))
	douyinCmd.AddCommand(douyinCommands.NewTrendingCmd(func() douyinCommands.Client { return dp.BuildClient() }))
	douyinCmd.AddCommand(douyinCommands.NewSearchCmd(func() douyinCommands.Client { return dp.BuildClient() }))
	douyinCmd.AddCommand(douyinCommands.NewVideoCmd(func() douyinCommands.Client { return dp.BuildClient() }))
	douyinCmd.AddCommand(douyinCommands.NewUserCmd(func() douyinCommands.Client { return dp.BuildClient() }))

	rootCmd.AddCommand(douyinCmd)
}