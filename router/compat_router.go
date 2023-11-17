package router

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

func CompatRouter(r fiber.Router) {
	r.Route("/compat", func(t fiber.Router) {
		t.Get("/", getCompatibilityList)
		t.Post("/search", searchCompatibilityList)
	})
}

type NativeStatus string

const (
	Yes      NativeStatus = "yes"
	No       NativeStatus = "no"
	Untested NativeStatus = "untested"
)

type GameStatus struct {
	Id           int          `json:"id"`
	Title        string       `json:"title"`
	Year         string       `json:"year"`
	Category     string       `json:"category"`
	NativeStatus NativeStatus `json:"native_status"`
	HorScal      bool         `json:"hor_scal"`
	VertScal     bool         `json:"vert_scal"`
	FovCorrect   bool         `json:"fovCorrect"`
	Solution     string       `json:"solution"`
	SolutionLink string       `json:"solution_link"`
	PreviewLink  string       `json:"preview_link"`
}

const offset int = 20

func getGameStatuses(startIndex int) ([]GameStatus, error) {
	gameStatuses := make([]GameStatus, 0)

	f, err := excelize.OpenFile("./data/data.xlsx")
	if err != nil {
		fmt.Println(err)
		return gameStatuses, err
	}

	// Get all the rows in the 32:9.
	rows, err := f.GetRows(f.WorkBook.Sheets.Sheet[0].Name)
	if err != nil {
		fmt.Println(err)
		return gameStatuses, err
	}

	var start int = startIndex + 7
	if start < 7 {
		start = 7
	}
	end := start + offset
	if end > len(rows)-1 {
		end = len(rows) - 1
	}

	fmt.Println("start:", start)
	fmt.Println("end:", end)

	finalRowIndex := 0
	for i := start; finalRowIndex < end; i++ {
		row := rows[i]

		if len(row) == 0 || len(row) < 6 || row[0] == "TYPE OF GAMES" || !strings.Contains(row[0], " (") {
			continue
		}

		id := finalRowIndex + 1

		splitTitle := strings.Split(row[0], " (")
		title := splitTitle[0]

		year, _ := strings.CutSuffix(splitTitle[1], ")")

		var nativeStatus NativeStatus
		nativeStatusStr := strings.Trim(strings.ToLower(row[2]), " ")
		if nativeStatusStr == "yes" || nativeStatusStr == "no" {
			nativeStatus = NativeStatus(nativeStatusStr)
		} else {
			nativeStatus = NativeStatus("untested")
		}

		horScal := strings.Trim(strings.ToLower(row[3]), " ") == "-"
		vertScal := strings.Trim(strings.ToLower(row[4]), " ") == "-"
		fovCorrent := strings.Trim(strings.ToLower(row[5]), " ") == "-"
		solution := strings.Trim(strings.ToLower(row[6]), " ")
		previewLink := strings.Trim(strings.ToLower(row[7]), " ")

		newGameStatus := GameStatus{
			Id:           id,
			Title:        title,
			Year:         year,
			Category:     row[1],
			NativeStatus: nativeStatus,
			HorScal:      horScal,
			VertScal:     vertScal,
			FovCorrect:   fovCorrent,
			Solution:     solution,
			SolutionLink: solution,
			PreviewLink:  previewLink,
		}
		gameStatuses = append(gameStatuses, newGameStatus)
		finalRowIndex++
	}
	return gameStatuses, nil
}

type PageInfo struct {
	nextPage  int
	lastRowId int
}

func getCompatibilityList(c *fiber.Ctx) error {

	startIndex := 0
	pageStr := c.Query("page")
	var page int = 0
	if len(pageStr) > 0 {
		page, err := strconv.Atoi(pageStr)
		if err == nil {
			startIndex = page * offset
		}
	}
	gameStatuses, err := getGameStatuses(startIndex)
	if err != nil {
		return err
		// return c.SendStatus(500)
	}

	fmt.Println("Length of status: {}", len(gameStatuses))
	//   print " hx-get='/api/compat/?page={{.PageInfo.nextPage}}' hx-trigger='revealed' hx-swap='afterend' "

	if c.Query("template") == "true" {

		var nextPage int = page + 1
		// if page < 1 {
		// 	nextPage = 1
		// } else {
		// 	nextPage = page + 1
		// }
		lastId := gameStatuses[len(gameStatuses)-1].Id
		pageInfo := PageInfo{
			lastRowId: lastId,
			nextPage:  nextPage,
		}
		return c.Render("partials/compatList", fiber.Map{"GameStatuses": gameStatuses, "PageInfo": pageInfo})
	}

	return c.JSON(gameStatuses)
}

func searchCompatibilityList(c *fiber.Ctx) error {
	gameStatuses, err := getGameStatuses(-1)
	if err != nil {
		return c.SendStatus(500)
	}
	searchTerm := c.FormValue("search", "$null")
	for _, game := range gameStatuses {
		hasSearchTerm := false
		if searchTerm == "$null" {
			hasSearchTerm = true
		}

		rowsToCheck := []string{game.Title, game.Category, game.Solution, game.SolutionLink, game.PreviewLink}

		for j := 0; j < len(rowsToCheck) && !hasSearchTerm; j++ {
			if strings.Contains(strings.ToLower(rowsToCheck[j]), strings.ToLower(searchTerm)) {
				hasSearchTerm = true
			}
		}
		if !hasSearchTerm {
			continue
		}
	}

	return c.Render("partials/compatList", fiber.Map{"GameStatuses": gameStatuses})
}
