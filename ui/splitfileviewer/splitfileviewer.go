package splitfileviewer

import (
	"GoMkvRn/data_access"
	"GoMkvRn/models"
	"GoMkvRn/service"
	page "GoMkvRn/ui"
	"GoMkvRn/ui/icon"
	"GoMkvRn/utils"
	"errors"
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/explorer"
	"github.com/flytam/filenamify"
	"github.com/maruel/natural"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// Page holds the state for a page demonstrating the features of
// the Menu component.
type Page struct {
	redButton, greenButton, blueButton        widget.Clickable
	balanceButton, accountButton, cartButton  widget.Clickable
	leftFillColor                             color.NRGBA
	leftContextArea                           component.ContextArea
	leftMenu, rightMenu, showMenu, seasonMenu component.MenuState
	menuInit                                  bool
	menuDemoList                              widget.List
	menuDemoListStates                        []component.ContextArea
	inputShowList                             widget.List
	inputShowListStates                       []component.ContextArea
	inputSeasonList                           widget.List
	inputSeasonListStates                     []component.ContextArea
	widget.List

	*page.Router
}

// New constructs a Page with the provided router.
func New(router *page.Router) *Page {
	return &Page{
		Router: router,
	}
}

var _ page.Page = &Page{}

func (p *Page) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (p *Page) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (p *Page) NavItem() component.NavItem {
	return component.NavItem{
		Name: "File Viewer",
		Icon: icon.FileViewerIcon,
	}
}

func getBulkFiles(paths ...string) map[string][]utils.BulkFile {
	var pathMap = map[string][]utils.BulkFile{}
	for _, path := range paths {
		_, ok := pathMap[path]
		if ok {
			continue
		}
		fl, err := utils.GetFiles(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(len(fl))
		sort.Slice(fl, func(i, j int) bool {
			return natural.Less(fl[i].Name, fl[j].Name)
		})
		pathMap[path] = fl
	}
	return pathMap
}

var bulkFileMap = getBulkFiles("B:\\Documents\\Movies\\BluRay\\FOOD_WARS\\Season_2_pt1", "B:\\Documents\\Movies\\BluRay\\FOOD_WARS\\Season_2_pt2")

func toSortedFiles(fileMap map[string][]utils.BulkFile) []utils.BulkFile {
	var files []utils.BulkFile
	for _, v := range fileMap {
		for _, bulkFile := range v {
			files = append(files, bulkFile)
		}
	}
	sort.Slice(files, func(i, j int) bool {
		return natural.Less(files[i].Path, files[j].Path)
	})
	return files
}

func getEpisodes(seasonId int64) []models.Episode {
	tvMazeDataAccess := data_access.TvMazeDataAccess{}
	tvMazeService := service.NewTvMazeService(tvMazeDataAccess)
	eps := tvMazeService.ListEpisodesBySeason(seasonId)
	var validEps []models.Episode
	for _, ep := range eps {
		if ep.EpisodeNumber != nil {
			validEps = append(validEps, ep)
		}
	}
	return validEps
}

func getShows(showName string) []models.ShowSearchResult {
	tvMazeDataAccess := data_access.TvMazeDataAccess{}
	tvMazeService := service.NewTvMazeService(tvMazeDataAccess)
	return tvMazeService.FindShowIdByName(showName)
}

func getSeasons(showId int64) []models.Season {
	tvMazeDataAccess := data_access.TvMazeDataAccess{}
	tvMazeService := service.NewTvMazeService(tvMazeDataAccess)
	return tvMazeService.ListSeasons(showId)
}

func bulkRename(bulk map[string][]utils.BulkFile, episodes []models.Episode, isSakeCase bool) error {
	bulkFiles := toSortedFiles(bulk)
	if len(bulkFiles) != len(episodes) {
		return errors.New(fmt.Sprintf("The files on the left and the episodes on the right do not have the count L: %d R: %d.", len(bulkFiles), len(episodes)))
	}
	for i, file := range bulkFiles {
		ep := episodes[i]
		epFileSafeName, err := filenamify.Filenamify(ep.Name, filenamify.Options{})
		if err != nil {
			return err
		}
		var renameError error = nil
		var epNumber int = -1
		if ep.EpisodeNumber != nil {
			epNumber = *ep.EpisodeNumber
		}

		if isSakeCase {
			snakeCasePath := filepath.Join(file.ParentDir, fmt.Sprintf("e%d_%s%s", epNumber, strings.ToLower(strings.ReplaceAll(epFileSafeName, " ", "_")), filepath.Ext(file.Path)))
			fmt.Printf("Converting %s to %s\n", file.Path, snakeCasePath)
			//renameError = os.Rename(file.Path, snakeCasePath)
		} else {
			unmodifiedPath := filepath.Join(file.ParentDir, fmt.Sprintf("E%d %s%s", epNumber, epFileSafeName, filepath.Ext(file.Path)))
			fmt.Printf("Converting %s to %s\n", file.Path, unmodifiedPath)
			//renameError = os.Rename(file.Path, unmodifiedPath)
		}
		if renameError != nil {
			return renameError
		}
	}

	return nil
}

var eps = getEpisodes(29087)
var epsBtnsUp = make([]widget.Clickable, len(eps))
var epsBtnsDown = make([]widget.Clickable, len(eps))
var convertBtn = widget.Clickable{}
var columns = layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}
var chkSnakeCase = widget.Bool{Value: true}
var expl = explorer.Explorer{}
var inputBulkFiles = widget.Editor{}
var btnLoadFiles = widget.Clickable{}
var inputShow = widget.Editor{}
var inputSeason = widget.Editor{}
var btnLoadShow = widget.Clickable{}
var radShowOpt = widget.Enum{}
var radShowBtnStyles []material.RadioButtonStyle
var radSeasonOpt = widget.Enum{}
var radSeasonBtnStyles []material.RadioButtonStyle
var btnLoadSeason = widget.Clickable{}

func (p *Page) Layout(gtx C, th *material.Theme) D {
	if !p.menuInit {
		p.leftMenu = component.MenuState{
			Options: []func(gtx C) D{
				func(gtx C) D {
					return layout.Inset{
						Left:  unit.Dp(16),
						Right: unit.Dp(16),
					}.Layout(gtx, material.Body1(th, "Menus support arbitrary widgets.\nThis is just a label!\nHere's a loader:").Layout)
				},
				component.Divider(th).Layout,
				func(gtx C) D {
					return layout.Inset{
						Top:    unit.Dp(4),
						Bottom: unit.Dp(4),
						Left:   unit.Dp(16),
						Right:  unit.Dp(16),
					}.Layout(gtx, func(gtx C) D {
						gtx.Constraints.Max.X = gtx.Dp(unit.Dp(24))
						gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(24))
						return material.Loader(th).Layout(gtx)
					})
				},
				component.SubheadingDivider(th, "Colors").Layout,
				component.MenuItem(th, &p.redButton, "Red").Layout,
				component.MenuItem(th, &p.greenButton, "Green").Layout,
				component.MenuItem(th, &p.blueButton, "Blue").Layout,
			},
		}
		p.rightMenu = component.MenuState{
			Options: []func(gtx C) D{
				func(gtx C) D {
					item := component.MenuItem(th, &p.balanceButton, "Balance")
					item.Icon = icon.AccountBalanceIcon
					item.Hint = component.MenuHintText(th, "Hint")
					return item.Layout(gtx)
				},
				func(gtx C) D {
					item := component.MenuItem(th, &p.accountButton, "Account")
					item.Icon = icon.AccountBoxIcon
					item.Hint = component.MenuHintText(th, "Hint")
					return item.Layout(gtx)
				},
				func(gtx C) D {
					item := component.MenuItem(th, &p.cartButton, "Cart")
					item.Icon = icon.CartIcon
					item.Hint = component.MenuHintText(th, "Hint")
					return item.Layout(gtx)
				},
			},
		}
	}
	if p.redButton.Clicked(gtx) {
		p.leftFillColor = color.NRGBA{R: 200, A: 255}
	}
	if p.greenButton.Clicked(gtx) {
		p.leftFillColor = color.NRGBA{G: 200, A: 255}
	}
	if p.blueButton.Clicked(gtx) {
		p.leftFillColor = color.NRGBA{B: 200, A: 255}
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return material.Editor(th, &inputBulkFiles, "Enter directories you want to bulk load.  Separate them by new lines").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					if btnLoadFiles.Clicked(gtx) {
						bulkDirs := strings.Split(inputBulkFiles.Text(), "\n")
						bulkFileMap = getBulkFiles(bulkDirs...)
					}
					return material.Button(th, &btnLoadFiles, "Load").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return material.Editor(th, &inputShow, "Enter the name of the TV show you want to look for.").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					if btnLoadShow.Clicked(gtx) {
						clear(radShowBtnStyles)
						for _, show := range getShows(strings.TrimSpace(inputShow.Text())) {
							radShowBtnStyles = append(radShowBtnStyles, material.RadioButton(th, &radShowOpt, strconv.FormatInt(show.Id, 10), show.Name))
							fmt.Println(show)
						}
					}
					return material.Button(th, &btnLoadShow, "Find Show").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					if len(radShowBtnStyles) == 0 {
						return D{}
					}
					return material.List(th, &p.inputShowList).Layout(gtx, len(radShowBtnStyles), func(gtx C, index int) D {
						p.inputShowList.Axis = layout.Vertical
						if len(p.inputShowListStates) < index+1 {
							p.inputShowListStates = append(p.inputShowListStates, component.ContextArea{})
						}
						state := &p.inputShowListStates[index]
						return layout.Stack{}.Layout(gtx,
							layout.Stacked(func(gtx C) D {
								gtx.Constraints.Min.X = gtx.Constraints.Max.X
								if radShowOpt.Update(gtx) {
									showId, err := strconv.ParseInt(radShowOpt.Value, 10, 64)
									if err != nil {
										fmt.Println(err)
										os.Exit(1)
									}
									for _, season := range getSeasons(showId) {
										radSeasonBtnStyles = append(radSeasonBtnStyles, material.RadioButton(th, &radShowOpt, strconv.FormatInt(season.Id, 10), season.Name))
										fmt.Println(season)
									}
								}
								return radShowBtnStyles[index].Layout(gtx)
							}),
							layout.Expanded(func(gtx C) D {
								return state.Layout(gtx, func(gtx C) D {
									gtx.Constraints.Min.X = 0
									return component.Menu(th, &p.showMenu).Layout(gtx)
								})
							}),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					if len(radSeasonBtnStyles) == 0 {
						return D{}
					}
					return material.List(th, &p.inputSeasonList).Layout(gtx, len(radSeasonBtnStyles), func(gtx C, index int) D {
						p.inputSeasonList.Axis = layout.Vertical
						if len(p.inputSeasonListStates) < index+1 {
							p.inputSeasonListStates = append(p.inputSeasonListStates, component.ContextArea{})
						}
						state := &p.inputSeasonListStates[index]
						return layout.Stack{}.Layout(gtx,
							layout.Stacked(func(gtx C) D {
								gtx.Constraints.Min.X = gtx.Constraints.Max.X
								if radSeasonOpt.Update(gtx) {
									seasonId, err := strconv.ParseInt(radSeasonOpt.Value, 10, 64)
									if err != nil {
										os.Exit(1)
									}
									eps = getEpisodes(seasonId)
								}
								return radSeasonBtnStyles[index].Layout(gtx)
							}),
							layout.Expanded(func(gtx C) D {
								return state.Layout(gtx, func(gtx C) D {
									gtx.Constraints.Min.X = 0
									return component.Menu(th, &p.seasonMenu).Layout(gtx)
								})
							}),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return material.CheckBox(th, &chkSnakeCase, "Output as snake_case").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					if convertBtn.Clicked(gtx) {
						fmt.Println("Convert button clicked")
						err := bulkRename(bulkFileMap, eps, chkSnakeCase.Value)
						if err != nil {
							fmt.Println(err)
							return layout.UniformInset(unit.Dp(8)).Layout(gtx, material.Body1(th, err.Error()).Layout)
						}
					}
					return material.Button(th, &convertBtn, "Convert").Layout(gtx)
				}),
			)
		}),
		layout.Flexed(0.5, func(gtx C) D {
			p.menuDemoList.Axis = layout.Vertical
			bulkFiles := toSortedFiles(bulkFileMap)
			return material.List(th, &p.menuDemoList).Layout(gtx, len(bulkFiles), func(gtx C, index int) D {
				if len(p.menuDemoListStates) < index+1 {
					p.menuDemoListStates = append(p.menuDemoListStates, component.ContextArea{})
				}
				state := &p.menuDemoListStates[index]
				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.UniformInset(unit.Dp(8)).Layout(gtx, material.Body1(th, bulkFiles[index].Path).Layout)
					}),
					layout.Expanded(func(gtx C) D {
						return state.Layout(gtx, func(gtx C) D {
							gtx.Constraints.Min.X = 0
							return component.Menu(th, &p.rightMenu).Layout(gtx)
						})
					}),
				)
			})
		}),
		layout.Flexed(0.5, func(gtx C) D {
			p.menuDemoList.Axis = layout.Vertical
			return material.List(th, &p.menuDemoList).Layout(gtx, len(eps), func(gtx C, index int) D {
				if len(p.menuDemoListStates) < index+1 {
					p.menuDemoListStates = append(p.menuDemoListStates, component.ContextArea{})
				}
				state := &p.menuDemoListStates[index]
				return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return columns.Layout(gtx,
							layout.Flexed(0.5, func(gtx C) D {
								gtx.Constraints.Min.X = gtx.Constraints.Max.X
								return layout.UniformInset(unit.Dp(8)).Layout(gtx, material.Body1(th, fmt.Sprintf("%d %s", *eps[index].EpisodeNumber, eps[index].Name)).Layout)
							}),
							layout.Flexed(0.20, func(gtx C) D {
								if index > 0 && epsBtnsUp[index].Clicked(gtx) {
									eps[index], eps[index-1] = eps[index-1], eps[index]
								}
								icon, _ := widget.NewIcon(icons.HardwareKeyboardArrowUp)
								return material.IconButton(th, &epsBtnsUp[index], icon, "move title Up").Layout(gtx)
							}),
							layout.Flexed(0.20, func(gtx C) D {
								if index >= 0 && epsBtnsDown[index].Clicked(gtx) {
									eps[index], eps[index+1] = eps[index+1], eps[index]
								}
								icon, _ := widget.NewIcon(icons.HardwareKeyboardArrowDown)
								return material.IconButton(th, &epsBtnsDown[index], icon, "move title Down").Layout(gtx)
							}))
					}),
					layout.Rigid(func(gtx C) D {
						return state.Layout(gtx, func(gtx C) D {
							gtx.Constraints.Min.X = 0
							return component.Menu(th, &p.rightMenu).Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}
