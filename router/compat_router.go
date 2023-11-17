package router

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

func CompatRouter(r fiber.Router) {
	r.Route("/compat", func(t fiber.Router) {
		t.Get("/", getCompatibilityList)
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

func getCompatibilityList(c *fiber.Ctx) error {
	f, err := excelize.OpenFile("./data/data.xlsx")
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Get all the rows in the 32:9.
	rows, err := f.GetRows(f.WorkBook.Sheets.Sheet[0].Name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	gameStatuses := make([]GameStatus, 0)

	for i := 7; i < len(rows); i++ {
		row := rows[i]

		if len(row) == 0 || len(row) < 6 || row[0] == "TYPE OF GAMES" {
			continue
		}

		if !strings.Contains(row[0], " (") {
			continue
		}

		splitTitle := strings.Split(row[0], " (")
		title := splitTitle[0]
		year := splitTitle[1]

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
			Id:           i,
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

		fmt.Println()

	}

	return c.JSON(gameStatuses)
}
