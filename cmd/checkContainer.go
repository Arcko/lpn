package cmd

import (
	"log"
	docker "lpn/docker"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "checkContainer",
	Short: "Check if there is a container created by lpn (Liferay Portal Nightly)",
	Long:  `Uses docker container inspect to check if there is a container created by lpn (Liferay Portal Nightly)`,
	Run: func(cmd *cobra.Command, args []string) {
		if docker.CheckDockerContainerExists() {
			log.Println("The container [" + docker.DockerContainerName + "] is running.")
		} else {
			log.Println("The container [" + docker.DockerContainerName + "] is NOT running.")
		}
	},
}