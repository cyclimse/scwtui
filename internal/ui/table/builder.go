package table

import (
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/cyclimse/scwtui/internal/resource"
	table "github.com/cyclimse/scwtui/internal/ui/table/custom"
	"github.com/mattn/go-runewidth"
	"github.com/xeonx/timeago"
)

const widthOfAnUUID = 36

// nolint: gochecknoglobals
var (
	titles = []string{
		"Status",
		"Name",
		"Type",
		"Project",
		"Created",
		"Locality",
	}

	altTitles = []string{
		"Status",
		"ID",
		"Type",
		"Project ID",
		"Created At",
		"Locality",
	}

	columnsWithFixedWidth = map[string]int{
		"Status":     runewidth.StringWidth("Status"),
		"Locality":   runewidth.StringWidth("Locality"),
		"ID":         widthOfAnUUID,
		"Project ID": widthOfAnUUID,
		"Created At": runewidth.StringWidth(time.RFC3339),
		"Type":       runewidth.StringWidth(resource.TypeContainerNamespace.String()),
		"Created":    runewidth.StringWidth("about an hour ago"),
	}
)

func NewBuilder(styles table.Styles) *Build {
	b := &Build{
		styles:       styles,
		timeagoCache: make(map[time.Time]string),
		mutex:        &sync.Mutex{},
	}

	timer := time.NewTicker(time.Minute)

	go func() {
		// purge the cache every minute to allow timeago to update.
		// eg. "a few seconds ago" -> "a minute ago"
		for range timer.C {
			b.mutex.Lock()
			b.timeagoCache = make(map[time.Time]string)
			b.mutex.Unlock()
		}
	}()

	return b
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
			b.formattedTimeAgo(metadata.CreatedAt),
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

		createdAt := ""
		if metadata.CreatedAt != nil {
			createdAt = metadata.CreatedAt.Format(time.RFC3339)
		}

		rows = append(rows, table.Row{
			lipgloss.PlaceHorizontal(6, lipgloss.Center, string(metadata.Status.Emoji(metadata.Type))),
			metadata.ID,
			metadata.Type.String(),
			metadata.ProjectID,
			createdAt,
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

	widthPerColumn := 0
	if len(titles) > fixedColumnsCount {
		widthPerColumn = (widthWithPadding - fixedColumnsWidth) / (len(titles) - fixedColumnsCount)
	}

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

	if widthPerColumn == 0 {
		cols[len(cols)-1].Width += widthWithPadding - fixedColumnsWidth
	}

	return cols
}

func (b Build) formattedTimeAgo(t *time.Time) string {
	if t == nil {
		return ""
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if formatted, ok := b.timeagoCache[*t]; ok {
		return formatted
	}

	formatted := timeago.English.Format(*t)
	b.timeagoCache[*t] = formatted

	return formatted
}

type Build struct {
	styles       table.Styles
	timeagoCache map[time.Time]string
	mutex        *sync.Mutex
}
