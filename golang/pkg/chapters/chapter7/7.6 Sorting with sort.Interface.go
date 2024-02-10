package chapter7

import (
	"fmt"
	htmltmpl "html/template"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

var tracks = []*Track{
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
}

type StringSlice []string

func (ss *StringSlice) Len() int {
	return len(*ss)
}

func (ss *StringSlice) Less(i, j int) bool {
	return (*ss)[i] < (*ss)[j]
}

func (ss *StringSlice) Swap(i, j int) {
	(*ss)[i], (*ss)[j] = (*ss)[j], (*ss)[i]
}

func IndirectStringSliceSorting() {
	names := []string{"Diana", "Dmitriy", "Denis", "Rustam", "Kate"}
	ssNames := StringSlice(names)

	fmt.Println(ssNames)

}

func DirectStringSliceSorting() {
	names := []string{"Diana", "Dmitriy", "Denis", "Rustam", "Kate"}

	sort.Strings(names)
	fmt.Println(names)
}

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

func PrintTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}

	tw.Flush()
}

var (
	templ = htmltmpl.Must(htmltmpl.New("trackList").Parse(
		`
		<h1>Tracklist:</h1>
		<table>
		<tr style = 'text-align:left'>
			<th>Title</th>
			<th>Artist</th>
			<th>Album</th>
			<th>Year</th>
			<th>Length</th>
		</tr>
		{{range .Tracks}}
		<tr>
			<td>{{.Title}}</td>
			<td>{{.Artist}}</td>
			<td>{{.Album}}</td>
			<td>{{.Year}}</td>
			<td>{{.Length}}</td>
		</tr>
		{{end}}
		</table>
		`))
)

// "tracklist.html"
func PrintTracksHTML(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("creating %s; %s", filename, err)
	}
	defer file.Close()

	byLengthTracks := byLength(tracks)
	sort.Sort(&byLengthTracks)
	if err := templ.Execute(file, map[string]interface{}{"Tracks": tracks}); err != nil {
		return fmt.Errorf("writing tracks; %s", err)
	}
	return nil
}

func TrackSortUsing() {
	// Unsorted data
	PrintTracks(tracks)

	fmt.Println()

	byLengthTracks := byLength(tracks)
	sort.Sort(&byLengthTracks)
	// Sorted in the ASC order
	PrintTracks(tracks)

	fmt.Println()

	sort.Sort(sort.Reverse(&byLengthTracks))
	// Sorted in the DESC order
	PrintTracks(tracks)
}

type byLength []*Track

func (bt *byLength) Len() int {
	return len(*bt)
}

func (bt *byLength) Less(i, j int) bool {
	return (*bt)[i].Length < (*bt)[j].Length
}

func (bt *byLength) Swap(i, j int) {
	(*bt)[i], (*bt)[j] = (*bt)[j], (*bt)[i]
}

/* The custom sort order */
type customSort struct {
	tracks []*Track
	less   func(x, y *Track) bool
}

func (cs *customSort) Len() int {
	return len(cs.tracks)
}

func (cs *customSort) Less(i, j int) bool {
	return cs.less(cs.tracks[i], cs.tracks[j])
}

func (cs *customSort) Swap(i, j int) {
	cs.tracks[i], cs.tracks[j] = cs.tracks[j], cs.tracks[i]
}

func CustomSortUsing() {
	sort.Sort(&customSort{tracks: tracks, less: func(x, y *Track) bool {
		if x.Title != y.Title {
			return x.Title < y.Title
		}
		if x.Year != y.Year {
			return x.Year < y.Year
		}
		if x.Length != y.Length {
			return x.Length < y.Length
		}
		return false
	},
	})

	PrintTracks(tracks)
}

func IsSortedUsing() {
	var values = []int{3, 1, 4, 5, 0}

	fmt.Println(values)
	fmt.Println(sort.IntsAreSorted(values)) // false

	// ASC []int sort
	sort.Ints(values)

	fmt.Println(values)
	fmt.Println(sort.IntsAreSorted(values)) // true

	// DESC []int sort
	fmt.Println(sort.Reverse(sort.IntSlice(values)))
	fmt.Println(sort.IntsAreSorted(values)) // false

}

/*
HOMEWORK
*/

/* 7.10 */
func IsPalindrome(s sort.Interface) bool {
	for i := 0; i < s.Len()/2; i++ {
		if s.Less(i, s.Len()-1-i) || s.Less(s.Len()-1-i, i) {
			return false
		}
	}

	return true
}

type SymbolComparableString string

func (scs SymbolComparableString) Len() int {
	return len(scs)
}

func (scs SymbolComparableString) Less(i, j int) bool {
	return (scs)[i] < (scs)[j]
}

func (scs SymbolComparableString) Swap(i, j int) {
	return
}
