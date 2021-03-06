package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"

	liferay "github.com/mdelapenya/lpn/liferay"
	tablewriter "github.com/olekukonko/tablewriter"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var imagesSize int
var imagesPage int

type imageResponse struct {
	Size         int
	Architecture string
	Variant      string
	Features     string
	OS           string
	OSVersion    string
	OSFeatures   string
}

type resultResponse struct {
	Name        string
	FullSize    int
	Images      []imageResponse
	ID          int64
	Repository  int64
	Creator     int64
	LastUpdater int64
	LastUpdated string
	ImageID     string
	V2          bool
}

type tagsResponse struct {
	Count    int
	Next     string
	Previous string
	Results  []resultResponse
}

func init() {
	rootCmd.AddCommand(tagsCmd)

	subcommands := []*cobra.Command{
		tagsCECmd, tagsCommerceCmd, tagsDXPCmd, tagsNightlyCmd, tagsReleaseCmd}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		subcommand.Flags().IntVarP(&imagesSize, "size", "s", 25, "Sets the number of tags to retrieve.")
		subcommand.Flags().IntVarP(&imagesPage, "page", "p", 1, "Sets the page element where tags exist.")

		tagsCmd.AddCommand(subcommand)

		subcommand.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Runs commands with Debug log level")
		subcommand.VisitParents(addVerboseFlag)
	}
}

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Lists all tags for Liferay Portal Docker image",
	Long: `Lists all tags for Liferay Portal Docker image from the Official repositories (see configuration file).
		For non-official Docker images, the tool lists tags from the unofficial repositories (see configuration file).
	For that, please run this command adding "commerce", "release" or "nightly" subcommands.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("tags requires zero or one argument representing the image tag to be pulled")
		}

		return nil
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		enableDebugLevel()
	},
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var tagsCECmd = &cobra.Command{
	Use:   "ce",
	Short: "Lists all tags for Liferay Portal CE Docker image",
	Long:  `Lists all tags for Liferay Portal CE Docker image from the Official repository`,
	Run: func(cmd *cobra.Command, args []string) {
		ce := liferay.CE{}

		readTags(ce, imagesSize, imagesPage)
	},
}

var tagsCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Lists all tags for Liferay Commerce Docker image",
	Long:  `Lists all tags for Liferay Commerce Docker image from the Official repository`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		readTags(commerce, imagesSize, imagesPage)
	},
}

var tagsDXPCmd = &cobra.Command{
	Use:   "dxp",
	Short: "Lists all tags for Liferay DXP Docker image",
	Long:  `Lists all tags for Liferay DXP Docker image from the Official repository`,
	Run: func(cmd *cobra.Command, args []string) {
		dxp := liferay.DXP{}

		readTags(dxp, imagesSize, imagesPage)
	},
}

var tagsNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Lists all tags for Liferay Portal Nightly Build Docker image",
	Long:  `Lists all tags for Liferay Portal Nightly Build Docker image from one of the unofficial repository`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		readTags(nightly, imagesSize, imagesPage)
	},
}

var tagsReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Lists all tags for Liferay Portal Release Docker image",
	Long:  `Lists all tags for Liferay Portal Release Docker image from one of the unofficial repository`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		readTags(release, imagesSize, imagesPage)
	},
}

func convertToHuman(bytes int) string {
	return fmt.Sprintf("%d MB", (bytes / 1000000))
}

func printTagsAsTable(data [][]string, currentPage int, totalPages int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Image:Tag", "Size"})
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT})
	table.SetFooter([]string{"", fmt.Sprintf("%d of %d", currentPage, totalPages)}) // Add Footer

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}

func readTags(image liferay.Image, count int, page int) {
	tagsPage := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags/?page_size=%d&page=%d", image.GetDockerHubTagsURL(), count, page)

	// Request the HTML page.
	res, err := http.Get(tagsPage)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error getting response from the server")
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		log.WithFields(log.Fields{
			"statusCode": res.StatusCode,
		}).Warn("There are no available tags for that pagination. Please use --page and --size arguments to filter properly")
		return
	}

	if res.StatusCode != 200 {
		log.WithFields(log.Fields{
			"status":     res.Status,
			"statusCode": res.StatusCode,
		}).Fatal("Error getting response from the server")
	}

	// Load the JSON document
	tagsResponse := new(tagsResponse)
	err = json.NewDecoder(res.Body).Decode(tagsResponse)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Error decoding response from the server")
	}

	data := [][]string{}

	for _, t := range tagsResponse.Results {
		// For each item found, get the tag and its size
		tag := t.Name
		size := t.Images[0].Size

		data = append(data, []string{tag, convertToHuman(size)})
	}

	if len(data) > 0 {
		totalPages := int(math.Ceil(float64(tagsResponse.Count) / float64(count)))
		if count > tagsResponse.Count {
			count = tagsResponse.Count
		}

		log.WithFields(log.Fields{
			"images":      tagsResponse.Count,
			"elements":    count,
			"currentPage": page,
			"totalPages":  totalPages,
		}).Infof("There are %d images, showing %d elements in page %d of %d", tagsResponse.Count, count, page, totalPages)

		printTagsAsTable(data, page, totalPages)
	} else {
		log.Info("There are no available tags for that pagination. Please use --page and --size arguments to filter properly")
	}
}
