package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type All struct {
	Updated                int64   `json:"updated"`
	Cases                  int64   `json:"cases"`
	TodayCases             int64   `json:"todayCases"`
	Deaths                 int64   `json:"deaths"`
	TodayDeaths            int64   `json:"todayDeaths"`
	Recovered              int64   `json:"recovered"`
	Active                 int64   `json:"active"`
	Critical               int64   `json:"critical"`
	CasesPerOneMillion     float64 `json:"casesPerOneMillion"`
	DeathsPerOneMillion    float64 `json:"deathsPerOneMillion"`
	Tests                  int64   `json:"tests"`
	TestsPerOneMillion     float64 `json:"testsPerOneMillion"`
	Population             int64   `json:"population"`
	ActivePerOneMillion    float64 `json:"activePerOneMillion"`
	RecoveredPerOneMillion float64 `json:"recoveredPerOneMillion"`
	AffectedCountries   int     `json:"affectedCountries"`
}

type Country struct {
	Updated                int64   `json:"updated"`
	Cases                  int64   `json:"cases"`
	TodayCases             int64   `json:"todayCases"`
	Deaths                 int64   `json:"deaths"`
	TodayDeaths            int64   `json:"todayDeaths"`
	Recovered              int64   `json:"recovered"`
	Active                 int64   `json:"active"`
	Critical               int64   `json:"critical"`
	CasesPerOneMillion     float64 `json:"casesPerOneMillion"`
	DeathsPerOneMillion    float64 `json:"deathsPerOneMillion"`
	Tests                  int64   `json:"tests"`
	TestsPerOneMillion     float64 `json:"testsPerOneMillion"`
	Population             int64   `json:"population"`
	ActivePerOneMillion    float64 `json:"activePerOneMillion"`
	RecoveredPerOneMillion float64 `json:"recoveredPerOneMillion"`
	CriticalPerOneMillion  float64 `json:"criticalPerOneMillion"`
	Continent              string  `json:"continent"`
	Country                string  `json:"country"`
	CountryInfo            struct {
		ID   int     `json:"_id"`
		Iso2 string  `json:"iso2"`
		Iso3 string  `json:"iso3"`
		Lat  float64 `json:"lat"`
		Long float64 `json:"long"`
		Flag string  `json:"flag"`
	} `json:"countryInfo"`
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("Failed to initialize termui: %v", err)
	}
	defer ui.Close()

	w, h := ui.TerminalDimensions()
	grid := ui.NewGrid()

	grid.SetRect(0, 0, w, h)
	loading := widgets.NewParagraph()
	loading.Text = `
                                           d8b      888
                                           Y8P      888
                                                    888
.d88b.   .d88b.   .d8888b .d88b.  888  888 888  .d88888
d88P"88b d88""88b d88P"   d88""88b 888  888 888 d88" 888
888  888 888  888 888     888  888 Y88  88P 888 888  888
Y88b 888 Y88..88P Y88b.   Y88..88P  Y8bd8P  888 Y88b 888
"Y88888  "Y88P"   "Y8888P "Y88P"    Y88P   888  "Y88888
    888
Y8b d88P
"Y88P"


Worldwide Coronavirus (COVID-19) Statistics for your terminal

[Please wait until information is loading](fg:black,bg:yellow)
	`
	grid.Set(ui.NewRow(1, loading))
	ui.Render(grid)

	time.Sleep(5 * time.Second)
	res, err := http.Get("https://disease.sh/v3/covid-19/all")
	if err != nil {
		log.Fatal(err)
	}
	var all All
	if err := json.NewDecoder(res.Body).Decode(&all); err != nil {
		log.Fatal(err)
	}

	global := widgets.NewParagraph()
	global.Title = "🌐 Global statistics"
	global.Text = fmt.Sprintf("[Infections](fg:blue): %d (%d today)\n", all.Cases, all.TodayCases)
	global.Text += fmt.Sprintf("[Deaths](fg:red): %d (%d today)\n", all.Deaths, all.TodayDeaths)
	global.Text += fmt.Sprintf("[Recoveries](fg:green): %d (%d remaining)\n", all.Recovered, all.Active)
	if all.Critical > 0 {
		global.Text += fmt.Sprintf("[Critical](fg:yellow): %d (%.2f%% of cases)\n", all.Critical, float64(all.Critical)/float64(all.Cases)*100)
	}
	global.Text += fmt.Sprintf("[Mortality rate (IFR)](fg:cyan): %.2f%%\n", float64(all.Deaths)/float64(all.Cases)*100)
	global.Text += fmt.Sprintf("[Mortality rate (CFR)](fg:cyan): %.2f%%\n", float64(all.Deaths)/(float64(all.Recovered)+float64(all.Deaths))*100)
	if all.AffectedCountries > 0 {
		global.Text += fmt.Sprintf("[Affected Countries](fg:magenta): %d\n", all.AffectedCountries)
	}
	global.SetRect(0, 0, 50, 10)
	global.BorderStyle.Fg = ui.ColorYellow
	global.TitleStyle = ui.NewStyle(ui.ColorClear)
	global.TextStyle = ui.NewStyle(ui.ColorClear)

	res2, err := http.Get("https://disease.sh/v3/covid-19/countries")
	if err != nil {
		log.Fatal(err)
	}
	countries := make([]Country, 0)
	if err := json.NewDecoder(res2.Body).Decode(&countries); err != nil {
		log.Fatal(err)
	}

	table := widgets.NewTable()
	tableHeader := []string{"#", "Country", "Total Cases", "Cases (today)", "Total Deaths", "Deaths (today)", "Recoveries", "Active", "Critical", "Mortality"}
	table.Rows = [][]string{tableHeader}

	for i, v := range countries {
		table.Rows = append(table.Rows, []string{
			fmt.Sprintf("%d", i+1),
			v.Country,
			fmt.Sprintf("%d", v.Cases),
			fmt.Sprintf("%d", v.TodayCases),
			fmt.Sprintf("%d", v.Deaths),
			fmt.Sprintf("%d", v.TodayDeaths),
			fmt.Sprintf("%d", v.Recovered),
			fmt.Sprintf("%d", v.Active),
			fmt.Sprintf("%d", v.Critical),
			fmt.Sprintf("%.2f%s", float64(v.Deaths)/float64(v.Cases)*100, "%"),
		})
	}

	table.ColumnWidths = []int{5, 22, 20, 20, 18, 18, 15, 15, 15, 15}
	table.TextAlignment = ui.AlignCenter
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.FillRow = true
	table.RowSeparator = false
	table.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierBold)
	table.BorderLeft = false
	table.BorderRight = false

	instructions := widgets.NewParagraph()
	instructions.Title = "👉 Readability test"

	instructions.Text = `[Please improve me!!](fg:black,bg:yellow) You're allowed to improve this project in the way you want 🙂 Also it'll be great if you implement a:
* Sort feature by columns.
* Auto refresh data.
* Remove old frames like the 👻 behind Global statistics box.

More info about data -> https://corona.lmao.ninja/docs/ and https://github.com/disease-sh/API
	`
	instructions.Border = true
	instructions.BorderStyle.Fg = ui.ColorYellow

	globalWidget := ui.NewRow(0.15, ui.NewCol(1.0, global))
	countriesTable := ui.NewRow(0.70, ui.NewCol(1.0, table))
	instructionsWidget := ui.NewRow(0.15, ui.NewCol(1.0, instructions))
	grid.Set(globalWidget, countriesTable, instructionsWidget)
	ui.Clear()
	ui.Render(grid)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>", "<Escape>":
			return
		}
		ui.Render(grid)
	}
}
