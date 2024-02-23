package main

import (
	"GoMkvRn/data_access"
	"GoMkvRn/models"
	"GoMkvRn/service"
	page "GoMkvRn/ui"
	"GoMkvRn/ui/menu"
	"GoMkvRn/ui/splitfileviewer"
	"errors"
	"flag"
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/flytam/filenamify"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// TODO: Add DI
func main() {
	injectPathPtr := flag.String("path", "", "Path to mkv folder")
	if injectPathPtr == nil || *injectPathPtr == "" {
		startApp()
	}
	//replace with path
	path := "B:\\Documents\\Movies\\BluRay\\FOOD_WARS\\Food_Wars_Season_3_pt3\\"
	RenameBulkCmd(path)
}

func startApp() {
	go func() {
		w := app.NewWindow()
		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops

	router := page.NewRouter()
	router.Register(1, splitfileviewer.New(&router))
	router.Register(3, menu.New(&router))

	for {
		switch e := w.NextEvent().(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			router.Layout(gtx, th)
			e.Frame(gtx.Ops)
		}
	}

}

func RenameBulkCmd(startPath string) {
	tvMazeDataAccess := data_access.TvMazeDataAccess{}
	tvMazeService := service.NewTvMazeService(tvMazeDataAccess)
	searchTerm := "Food Wars"
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
		if confirm == "Y" {
			strEpNumber := strconv.Itoa(selectedEpisodeNumber)
			fileSafeName, err := filenamify.Filenamify(selectedEpisode.Name, filenamify.Options{})
			if err != nil {
				fmt.Printf("Failed to convert file name %s to be safe", selectedEpisode.Name)
				fmt.Println(err)
			}
			isMultipleWhitespacePresent := regexp.MustCompile(`\s`).MatchString(fileSafeName)
			var mutiWhtieSpaceConverter int
			if isMultipleWhitespacePresent {
				fmt.Printf("We detected spaces in the name of your file.  Would you like keep them or convert them?\n")
				fmt.Printf("\n0) Do not Convert\n1) snake_case\n")
				fmt.Scan(&mutiWhtieSpaceConverter)
				fmt.Println(startPath + "E" + strEpNumber + "_" + fileSafeName)
				var reNameErr error
				if mutiWhtieSpaceConverter <= 0 {
					reNameErr = os.Rename(startPath+file.Name(), startPath+"E"+strEpNumber+"_"+fileSafeName+filepath.Ext(startPath+file.Name()))
				} else {
					reNameErr = os.Rename(startPath+file.Name(), startPath+"e"+strEpNumber+"_"+strings.ToLower(strings.ReplaceAll(fileSafeName, " ", "_")+filepath.Ext(startPath+file.Name())))
				}
				if reNameErr != nil {
					fmt.Printf("Failed to rename %s to %s\n", file.Name(), fileSafeName)
					fmt.Println(reNameErr)
					os.Exit(-1)
				}
			} else {
				reNameErr := os.Rename(startPath+file.Name(), startPath+"E"+strEpNumber+"_"+fileSafeName+filepath.Ext(startPath+file.Name()))
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
