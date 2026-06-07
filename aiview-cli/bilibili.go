package main

import (
	biliplatform "github.com/jackwener/aiview/internal/platform/bilibili"
	"github.com/jackwener/aiview/internal/platform"
	biliCommands "github.com/jackwener/aiview/commands/bilibili"
	"github.com/spf13/cobra"
)

func init() {
	p, ok := platform.GetPlatform("bilibili")
	if !ok {
		return
	}
	bp := p.(*biliplatform.BilibiliPlatform)

	biliCmd := &cobra.Command{
		Use:   "bilibili",
		Short: "Bilibili platform commands",
		Long:  `Commands for interacting with Bilibili content.`,
	}

	// Auth commands
	biliCmd.AddCommand(biliCommands.NewLoginCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewLogoutCmd(bp.GetAuthStore()))
	biliCmd.AddCommand(biliCommands.NewStatusCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewWhoamiCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))

	// Video commands
	biliCmd.AddCommand(biliCommands.NewVideoCmd(func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewAudioCmd(func() biliCommands.Client { return bp.BuildClient() }))

	// Discovery commands
	biliCmd.AddCommand(biliCommands.NewHotCmd(func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewRankCmd(func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewFeedCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))

	// Search
	biliCmd.AddCommand(biliCommands.NewSearchCmd(func() biliCommands.Client { return bp.BuildClient() }))

	// User commands
	biliCmd.AddCommand(biliCommands.NewUserCmd(func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewUserVideosCmd(func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewDynamicCmd(func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewDynamicPostCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewDynamicDeleteCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewCollectionCmd(func() biliCommands.Client { return bp.BuildClient() }))

	// Collections
	biliCmd.AddCommand(biliCommands.NewFavoritesCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewFollowingCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewHistoryCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewWatchLaterCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))

	// Interactions
	biliCmd.AddCommand(biliCommands.NewLikeCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewCoinCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewTripleCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewUnfollowCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewFavoriteCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))

	// Comment
	biliCmd.AddCommand(biliCommands.NewCommentCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewCommentDeleteCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))

	// Danmaku
	biliCmd.AddCommand(biliCommands.NewDanmakuCmd(func() biliCommands.Client { return bp.BuildClient() }))
	biliCmd.AddCommand(biliCommands.NewDanmakuSendCmd(bp.GetAuthStore(), func() biliCommands.Client { return bp.BuildClient() }))

	// Recommend
	biliCmd.AddCommand(biliCommands.NewRecommendCmd(func() biliCommands.Client { return bp.BuildClient() }))

	// Tags
	biliCmd.AddCommand(biliCommands.NewTagsCmd(func() biliCommands.Client { return bp.BuildClient() }))

	// Suggest
	biliCmd.AddCommand(biliCommands.NewSuggestCmd(func() biliCommands.Client { return bp.BuildClient() }))

	// Fans
	biliCmd.AddCommand(biliCommands.NewFansCmd(func() biliCommands.Client { return bp.BuildClient() }))

	rootCmd.AddCommand(biliCmd)
}