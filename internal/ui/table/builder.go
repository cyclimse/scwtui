package table

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/cyclimse/scwtui/internal/resource"
	table "github.com/cyclimse/scwtui/internal/ui/table/custom"
	"github.com/mattn/go-runewidth"
)

const widthOfAnUUID = 36

// nolint: gochecknoglobals
var (
	titles = []string{
		"Status",
		"Name",
		"Type",
		"Project",
		"Locality",
	}

	altTitles = []string{
		"Status",
		"ID",
		"Type",
		"Project ID",
		"Locality",
	}

	columnsWithFixedWidth = map[string]int{
		"Status":     runewidth.StringWidth("Status"),
		"Locality":   runewidth.StringWidth("Locality"),
		"ID":         widthOfAnUUID,
		"Project ID": widthOfAnUUID,
	}
)

func NewBuilder(styles table.Styles) *Build {
	return &Build{
		styles: styles,
	}
}

type BuildParams struct {
	Width             int
	AltView           bool
	Resources         []resource.Resource
	ProjectIDsToNames map[string]string
}

func (b *Build) Build(params BuildParams, opts ...table.Option) table.Model {
	var rows []table.Row
	if !params.AltView {
		rows = b.buildRows(params)
	} else {
		rows = b.buildRowsAlt(params)
	}
	cols := b.buildCols(params)

	return table.New(append(
		opts,
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithStyles(b.styles))...)
}

func (b *Build) buildRows(params BuildParams) []table.Row {
	resources := params.Resources
	rows := make([]table.Row, 0, len(resources))

	for _, r := range resources {
		metadata := r.Metadata()

		rows = append(rows, table.Row{
			lipgloss.PlaceHorizontal(6, lipgloss.Center, string(metadata.Status.Emoji(metadata.Type))),
			metadata.Name,
			metadata.Type.String(),
			params.ProjectIDsToNames[metadata.ProjectID],
			metadata.Locality.String(),
		})
	}

	return rows
}

func (b *Build) buildRowsAlt(params BuildParams) []table.Row {
	resources := params.Resources
	rows := make([]table.Row, 0, len(resources))

	for _, r := range resources {
		metadata := r.Metadata()
		rows = append(rows, table.Row{
			lipgloss.PlaceHorizontal(6, lipgloss.Center, string(metadata.Status.Emoji(metadata.Type))),
			metadata.ID,
			metadata.Type.String(),
			metadata.ProjectID,
			metadata.Locality.String(),
		})
	}

	return rows
}

// reduce the magic numbers.
func (b *Build) buildCols(params BuildParams) []table.Column {
	widthWithPadding := params.Width - 3
	titles := titles
	if params.AltView {
		titles = altTitles
		widthWithPadding -= 2
	}

	fixedColumnsWidth := 0
	fixedColumnsCount := 0
	for _, title := range titles {
		if width, ok := columnsWithFixedWidth[title]; ok {
			fixedColumnsWidth += width
			fixedColumnsCount++
		}
	}

	widthPerColumn := (widthWithPadding - fixedColumnsWidth) / (len(titles) - fixedColumnsCount)
	if widthPerColumn < 0 {
		widthPerColumn = 0
	}

	cols := make([]table.Column, 0, len(titles))
	for _, title := range titles {
		var width int
		if w, ok := columnsWithFixedWidth[title]; ok {
			width = w
		} else {
			width = widthPerColumn
		}

		cols = append(cols, table.Column{
			Title: title,
			Width: width,
		})
	}

	return cols
}

type Build struct {
	styles table.Styles
}
