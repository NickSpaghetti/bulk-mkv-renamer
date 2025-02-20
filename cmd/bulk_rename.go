package main

import (
	"GoMkvRn/data_access"
	"GoMkvRn/models"
	"GoMkvRn/service"
	"errors"
	"fmt"
	"github.com/flytam/filenamify"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	path_flag_name           string = "path"
	path_short                      = "p"
	content_name_flag_name   string = "content-name"
	content_name_short       string = "c"
	file_name_case_flag_name string = "file-name-case"
	file_name_case_short     string = "f"
)

const (
	ask         = "ask"
	defaultCase = "default"
	snakeCase   = "snake_case"
)

type FileNameCase int8

const (
	DefaultCase FileNameCase = iota
	SnakeCase
	Ask
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(bulkRenameCmd)

	bulkRenameCmd.Flags().StringP(path_flag_name, path_short, "", "Specify the path to the directory")
	bulkRenameCmd.Flags().StringP(content_name_flag_name, content_name_short, "", "Specify the name of the movie or tv show")
	bulkRenameCmd.Flags().StringP(file_name_case_flag_name, file_name_case_short, "default", "Specify the case for file names (default, snake_case")

	if err := bulkRenameCmd.MarkFlagRequired("path"); err != nil {
		panic(err)
	}

	if err := bulkRenameCmd.MarkFlagRequired("content-name"); err != nil {
		panic(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "bulk-rename",
	Short: "A tool for bulk file renaming files",
}

var bulkRenameCmd = &cobra.Command{
	Use:   "bfr",
	Short: "Renames files in a given directory",
	Long: `Renames files in a given directory.
               This command takes three flags
					path: 'File path where your .mkv files live'
					content-name: 'The content of your mkv files are your directory'
					file-name-case: The case for file names (default or camelcase)
               This will then go though the files and ask for input to rename what file to what ep or movie`,
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString(path_flag_name)
		contentName, _ := cmd.Flags().GetString(content_name_flag_name)
		fileNameCase, _ := cmd.Flags().GetString(file_name_case_flag_name)
		fileNameCaseId := GetFileNameCaseId(fileNameCase)
		RenameBulkCmd(path, contentName, fileNameCaseId)

	},
}

func GetFileNameCaseId(fileNameCase string) FileNameCase {
	switch strings.ToLower(fileNameCase) {
	case "":
		fallthrough
	case defaultCase:
		return DefaultCase
	case snakeCase:
		return SnakeCase
	case ask:
		return Ask
	default:
		panic(fmt.Sprintf("Unsupported file name case: %s", fileNameCase))
	}
}

func RenameBulkCmd(startPath string, contentName string, fileNameCaseId FileNameCase) {
	tvMazeDataAccess := data_access.TvMazeDataAccess{}
	tvMazeService := service.NewTvMazeService(tvMazeDataAccess)
	searchTerm := contentName
	searchResults := tvMazeService.FindShowIdByName(searchTerm)
	fmt.Printf("Here is a list of shows found with the search term: %s\n", searchTerm)
	for i, searchResult := range searchResults {
		fmt.Printf("%d) %s\n", i, searchResult.Name)
	}
	fmt.Printf("Enter the number of which show you want to select the episodes for. ie %d to %d\n", 0, len(searchResults)-1)
	var selectionIndex int
	fmt.Scanln(&selectionIndex)
	selection := searchResults[selectionIndex]
	fmt.Printf("You have selected %s\n", selection.Name)
	fmt.Println(selection.Id)
	seasons := tvMazeService.ListSeasons(selection.Id)
	fmt.Printf("Which seasons do you want to rename your files for. ie 1 to %d\n", seasons[len(seasons)-1].Number)
	// TODO: Some seasons do not have a name so we should selected by season number (some seasons such as Chopped have seasons as years so index wont work)
	for _, seasons := range seasons {
		fmt.Printf("%d) %s\n", seasons.Number, seasons.Name)
	}
	var seasonNumberSelected int
	fmt.Scanln(&seasonNumberSelected)
	fmt.Printf("You have selected season %d %s\n", seasonNumberSelected, seasons[seasonNumberSelected-1].Name)
	episodes := tvMazeService.ListEpisodesBySeason(seasons[seasonNumberSelected-1].Id)
	files, err := os.ReadDir(startPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	for _, file := range files {
		for _, episode := range episodes {
			// When EpisodeNumber is null it's typically a special
			if episode.EpisodeNumber != nil {
				fmt.Printf("%d) %s\n", *episode.EpisodeNumber, episode.Name)
			}
		}
		fmt.Printf("%s should be renamed to what?\n", file.Name())
		var selectedEpisodeNumber int
		fmt.Scan(&selectedEpisodeNumber)
		selectedEpisode, err := GetEpisodeByEpisodeNumber(selectedEpisodeNumber, episodes)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		fmt.Printf("Do you want to rename %s to %s? Y/n \n", file.Name(), selectedEpisode.Name)
		var confirm string
		fmt.Scan(&confirm)
		if confirm == "Y" || confirm == "y" {
			strEpNumber := strconv.Itoa(selectedEpisodeNumber)
			fileSafeName, err := filenamify.Filenamify(selectedEpisode.Name, filenamify.Options{})
			if err != nil {
				fmt.Printf("Failed to convert file name %s to be safe", selectedEpisode.Name)
				fmt.Println(err)
			}
			isMultipleWhitespacePresent := regexp.MustCompile(`\s`).MatchString(fileSafeName)
			var mutiWhiteSpaceConverter FileNameCase = fileNameCaseId
			if isMultipleWhitespacePresent {
				if fileNameCaseId == Ask {
					fmt.Printf("Detected spaces in the name of your file.  Would you like keep them or convert them?\n")
					fmt.Printf("\n0) Do not Convert\n1) snake_case\n")
					fmt.Scan(&mutiWhiteSpaceConverter)
					fmt.Println(startPath + "E" + strEpNumber + "_" + fileSafeName)
					if mutiWhiteSpaceConverter != DefaultCase || mutiWhiteSpaceConverter != SnakeCase {
						mutiWhiteSpaceConverter = DefaultCase
					}
				}
				var reNameErr error
				switch mutiWhiteSpaceConverter {
				case DefaultCase:
					reNameErr = os.Rename(startPath+file.Name(), filepath.Join(startPath, "E"+strEpNumber+"_"+fileSafeName+filepath.Ext(startPath+file.Name())))
				case SnakeCase:
					reNameErr = os.Rename(startPath+file.Name(), filepath.Join(startPath, "e"+strEpNumber+"_"+strings.ToLower(strings.ReplaceAll(fileSafeName, " ", "_")+filepath.Ext(startPath+file.Name()))))
				}

				if reNameErr != nil {
					fmt.Printf("Failed to rename %s to %s\n", file.Name(), fileSafeName)
					fmt.Println(reNameErr)
					os.Exit(-1)
				}
			} else {
				reNameErr := os.Rename(startPath+file.Name(), filepath.Join(startPath, "E"+strEpNumber+"_"+fileSafeName+filepath.Ext(startPath+file.Name())))
				if reNameErr != nil {
					fmt.Printf("Failed to rename %s to %s\n", file.Name(), fileSafeName)
					fmt.Println(reNameErr)
					os.Exit(-1)
				}
			}
		}
	}
}

func GetEpisodeByEpisodeNumber(episodeNumber int, episodes []models.Episode) (*models.Episode, error) {
	for _, episode := range episodes {
		if *episode.EpisodeNumber == episodeNumber {
			return &episode, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No episode found with the number %d", episodeNumber))
}
