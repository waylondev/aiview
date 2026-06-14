package commands

import (
	"fmt"
	"time"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/scheduler"
	"github.com/spf13/cobra"
)

// NewScheduleCmd creates the schedule command.
func NewScheduleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Manage scheduled tasks",
		Long:  "Add, list, or remove scheduled data collection tasks.",
	}

	cmd.AddCommand(newScheduleAddCmd())
	cmd.AddCommand(newScheduleListCmd())
	cmd.AddCommand(newScheduleRemoveCmd())

	return cmd
}

func newScheduleAddCmd() *cobra.Command {
	var interval, command, id string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new scheduled task",
		Long:  "Add a new scheduled task to run periodically.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse interval
			dur, err := scheduler.ParseInterval(interval)
			if err != nil {
				return fmt.Errorf("invalid interval: %w", err)
			}

			// Create scheduler (in-memory for now)
			s := scheduler.New()

			// Generate ID if not provided
			if id == "" {
				id = fmt.Sprintf("job_%d", time.Now().Unix())
			}

			// Add job
			if err := s.AddJob(id, dur, command); err != nil {
				return fmt.Errorf("add job: %w", err)
			}

			fmt.Printf("✓ Scheduled task added: %s (every %s)\n", id, interval)
			fmt.Printf("  Command: %s\n", command)
			fmt.Printf("  Note: Tasks are in-memory and will be lost on restart.\n")

			return nil
		},
	}

	cmd.Flags().StringVar(&interval, "every", "1h", "Interval (e.g., 30s, 5m, 1h)")
	cmd.Flags().StringVar(&command, "command", "", "Command to execute")
	cmd.Flags().StringVar(&id, "id", "", "Job ID (auto-generated if not provided)")
	cmd.MarkFlagRequired("command")

	return cmd
}

func newScheduleListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all scheduled tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create scheduler (in-memory, so this will be empty)
			s := scheduler.New()
			jobs := s.ListJobs()

			if len(jobs) == 0 {
				fmt.Println("No scheduled tasks.")
				return nil
			}

			fmt.Println("\nScheduled Tasks:")
			fmt.Println("────────────────────────────────────────────────────────")
			fmt.Printf("%-20s %-15s %-30s %s\n", "ID", "Interval", "Command", "Next Run")
			fmt.Println("────────────────────────────────────────────────────────")

			for _, job := range jobs {
				fmt.Printf("%-20s %-15s %-30s %s\n",
					job.ID,
					job.Interval.String(),
					job.Command,
					job.NextRun.Format("2006-01-02 15:04:05"),
				)
			}

			return nil
		},
	}

	return cmd
}

func newScheduleRemoveCmd() *cobra.Command {
	var id string

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a scheduled task",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create scheduler (in-memory)
			s := scheduler.New()

			if err := s.RemoveJob(id); err != nil {
				return aiverr.APIError("schedule", fmt.Sprintf("remove job: %v", err))
			}

			fmt.Printf("✓ Scheduled task removed: %s\n", id)
			return nil
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "Job ID to remove")
	cmd.MarkFlagRequired("id")

	return cmd
}
