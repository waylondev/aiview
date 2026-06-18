package bilibili

import bilibilitypes "github.com/jackwener/aiview/internal/platform/bilibili/bilibilitypes"

// API response types used across the bilibili platform layer.
// These are type aliases to bilibilitypes to avoid import cycles with the commands layer.

type (
	VideoInfo        = bilibilitypes.VideoInfo
	OwnerInfo        = bilibilitypes.OwnerInfo
	VideoStats       = bilibilitypes.VideoStats
	UserInfo         = bilibilitypes.UserInfo
	SearchUserResult = bilibilitypes.SearchUserResult
	SearchVideoResult = bilibilitypes.SearchVideoResult
	CommentInfo      = bilibilitypes.CommentInfo
	AuthorInfo       = bilibilitypes.AuthorInfo
	SubtitleInfo     = bilibilitypes.SubtitleInfo
	SubtitleItem     = bilibilitypes.SubtitleItem
	FavoriteFolder   = bilibilitypes.FavoriteFolder
	FavoriteMedia    = bilibilitypes.FavoriteMedia
	FollowingUser    = bilibilitypes.FollowingUser
	HistoryItem      = bilibilitypes.HistoryItem
	WatchLaterItem   = bilibilitypes.WatchLaterItem
	DynamicItem      = bilibilitypes.DynamicItem
	DynamicStats     = bilibilitypes.DynamicStats
	Credential       = bilibilitypes.BiliCredential
	DanmakuInfo      = bilibilitypes.DanmakuInfo
	VideoTag         = bilibilitypes.VideoTag
	FansUserInfo     = bilibilitypes.FansUserInfo
)