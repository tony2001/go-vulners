package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/pterm/pterm"
	v3 "github.com/tony2001/go-vulners/api/v3"
)

type Command struct {
	fs     *flag.FlagSet
	client *v3.ClientWithResponses
	cfg    Config
}

func (c *Command) Name() string {
	return c.fs.Name()
}

func (c *Command) Init(args []string) error {
	return c.fs.Parse(args)
}

type SearchCommand struct {
	Command
	query string
	size  int
	skip  int
}

func NewSearchCommand(cfg Config, client *v3.ClientWithResponses) *SearchCommand {
	sc := &SearchCommand{
		Command: Command{
			fs:     flag.NewFlagSet("search", flag.ContinueOnError),
			client: client,
			cfg:    cfg,
		},
	}

	sc.fs.StringVar(&sc.query, "query", "", "Search query")
	sc.fs.IntVar(&sc.size, "size", 20, "Number of results requested")
	sc.fs.IntVar(&sc.skip, "skip", 0, "Number of results to skip")
	return sc
}

func (sc *SearchCommand) Run() error {
	request := v3.SearchJSONRequestBody{
		ApiKey: sc.cfg.ApiKey,
		Query:  sc.query,
		Size:   &sc.size,
		Skip:   &sc.skip,
	}

	res, err := sc.client.SearchWithResponse(context.Background(), request)
	if err != nil {
		return err
	}

	if res.JSON200 != nil {
		data, err := res.JSON200.Data.AsSearchResponseDataSchema()
		if err != nil {
			return err
		}

		searchItems := make([]pterm.BulletListItem, 0, len(data.Search))
		for i := range data.Search {
			result := data.Search[i]
			searchItems = append(searchItems, []pterm.BulletListItem{
				{Level: 0, Text: result.Id, TextStyle: pterm.NewStyle(pterm.FgGreen), Bullet: fmt.Sprintf("%d", i+1), BulletStyle: pterm.NewStyle(pterm.FgGreen)},
				{Level: 1, Text: fmt.Sprintf("Score: %f", result.Score), TextStyle: pterm.NewStyle(pterm.FgGreen), Bullet: ""},
				{Level: 1, Text: fmt.Sprintf("Type: %s", result.Type), TextStyle: pterm.NewStyle(pterm.FgGreen), Bullet: ""},
				{Level: 1, Text: fmt.Sprintf("Description: %s", result.FlatDescription), TextStyle: pterm.NewStyle(pterm.FgGreen), Bullet: ""},
				{Level: 1, Text: "", Bullet: ""},
			}...)
		}

		pterm.DefaultBulletList.WithItems(searchItems).Render()
	}
	return nil
}
