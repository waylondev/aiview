package commands

import (
	"fmt"
	"strconv"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewLikeCmd creates the like command.
func NewLikeCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "like <BV>",
		Short: "Like video",
		Long:  `Like a video (login and write permission required).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(true)
			if err != nil {
				output.EmitError("not_authenticated", "Login with write permission required, use aiview bilibili login", format)
				return err
			}
			_ = cred

			bvid, err := ExtractBVID(args[0])
			if err != nil {
				output.EmitError("invalid_input", err.Error(), format)
				return err
			}

			if err := client.LikeVideo(bvid, false); err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to like: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(ActionResult{Success: true, Action: "like"}, format)
			}

			fmt.Println("✅ Liked")
			return nil
		},
	}
}

// NewCoinCmd creates the coin command.
func NewCoinCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	var num int

	cmd := &cobra.Command{
		Use:   "coin <BV>",
		Short: "Coin video",
		Long:  `Give coins to a video (login and write permission required).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(true)
			if err != nil {
				output.EmitError("not_authenticated", "Login with write permission required", format)
				return err
			}
			_ = cred

			bvid, err := ExtractBVID(args[0])
			if err != nil {
				output.EmitError("invalid_input", err.Error(), format)
				return err
			}

			if err := client.CoinVideo(bvid, num); err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to give coins: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(ActionResult{Success: true, Action: "coin"}, format)
			}

			fmt.Printf("✅ Gave %d coins\n", num)
			return nil
		},
	}

	cmd.Flags().IntVarP(&num, "num", "n", 1, "Number of coins (1-2)")
	return cmd
}

// NewTripleCmd creates the triple command.
func NewTripleCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "triple <BV>",
		Short: "Triple like",
		Long:  `Like, coin, and favorite a video in one go (login and write permission required).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(true)
			if err != nil {
				output.EmitError("not_authenticated", "Login with write permission required", format)
				return err
			}
			_ = cred

			bvid, err := ExtractBVID(args[0])
			if err != nil {
				output.EmitError("invalid_input", err.Error(), format)
				return err
			}

			if err := client.TripleVideo(bvid); err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to triple: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(ActionResult{Success: true, Action: "triple"}, format)
			}

			fmt.Println("✅ Triple action successful")
			return nil
		},
	}
}

// NewUnfollowCmd creates the unfollow command.
func NewUnfollowCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "unfollow <UID>",
		Short: "Unfollow",
		Long:  `Unfollow a user (login and write permission required).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(true)
			if err != nil {
				output.EmitError("not_authenticated", "Login with write permission required", format)
				return err
			}
			_ = cred

			uid, err := strconv.Atoi(args[0])
			if err != nil {
				output.EmitError("invalid_input", "UID must be a number", format)
				return err
			}

			if err := client.UnfollowUser(uid); err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to unfollow: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(ActionResult{Success: true, Action: "unfollow"}, format)
			}

			fmt.Println("✅ Unfollowed")
			return nil
		},
	}
}