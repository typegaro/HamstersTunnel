package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/typegaro/HamstersTunnel/internal/cli"
	"os"
)

var cliInstance *cli.CLI

var rootCmd = &cobra.Command{
	Use:   "HamstersTunnel",
	Short: "CLI to manage the HamstersTunnel server",
	Long:  `CLI to manage the HamstersTunnel server.`, //TODO: Edit this
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new service",
	Run: func(cmd *cobra.Command, args []string) {
		ip, _ := cmd.Flags().GetString("ip")
		name, _ := cmd.Flags().GetString("name")
		tcp, _ := cmd.Flags().GetString("tcp")
		udp, _ := cmd.Flags().GetString("udp")
		http, _ := cmd.Flags().GetString("http")
		save, _ := cmd.Flags().GetBool("save")

		if tcp == "" || udp == "" || http == "" {
			fmt.Println(
				"Error: At least one of the required parameters (tcp, udp, http) must be provided.",
			)
			cmd.Usage()
			return
		}

		cliInstance.NewService(ip, name, tcp, udp, http, save)
	},
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "get all the service running",
	Run: func(cmd *cobra.Command, args []string) {
		inactive, _ := cmd.Flags().GetBool("la")

		cliInstance.ListService(inactive)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop service by id",
	Run: func(cmd *cobra.Command, args []string) {
		remote, _ := cmd.Flags().GetBool("remote")
		id, _ := cmd.Flags().GetString("id")

		cliInstance.StopService(id, remote)
	},
}

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove service by id",
	Run: func(cmd *cobra.Command, args []string) {
		remote, _ := cmd.Flags().GetBool("remote")
		id, _ := cmd.Flags().GetString("id")

		cliInstance.RemoveService(id, remote)
	},
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	// Initialize the CLI instance
	cliInstance = &cli.CLI{}

	//new command
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().String("ip", "", "IP address of the server")
	newCmd.MarkFlagRequired("ip")
	newCmd.Flags().String("name", "", "Name of the service")
	newCmd.MarkFlagRequired("name")
	newCmd.Flags().String("tcp", "", "Local TCP ports")
	newCmd.Flags().String("udp", "", "Local UDP ports")
	newCmd.Flags().String("http", "", "Local HTTP ports")
	newCmd.Flags().Bool(
		"save",
		false,
		fmt.Sprintf("If set, the service configuration will be saved in the directory specified by the DEAMON_SERVICE_DIR environment variable (currently: %s)", os.Getenv("DEAMON_SERVICE_DIR")),
	)

	//ls command
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().Bool("inactive", false, "show also inactive services")

	//stop command
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().String("id", "", "Id of the service")
	stopCmd.Flags().Bool("remote", false, "propagate command on the Server")

	//remove command
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().String("id", "", "Id of the service")
	removeCmd.Flags().Bool("remote", false, "propagate command on the Server")

	// Execute the command
	rootCmd.Execute()
}
