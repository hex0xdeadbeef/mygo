package logic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	cfg "golang/pkg/projects/chapter4/b_xkcdtool/config"
	lg "golang/pkg/projects/chapter4/b_xkcdtool/logger"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ComicBody struct {
	Number int `json:"num,omitempty"`

	Link     string `json:"link,omitempty"`
	ImageURL string `json:"img,omitempty"`
	News     string `json:"news,omitempty"`

	Month string `json:"month,omitempty"`
	Day   string `json:"day,omitempty"`
	Year  string `json:"year,omitempty"`

	Title         string `json:"title,omitempty"`
	SafeTitle     string `json:"safe_title,omitempty"`
	Transcription string `json:"transcript,omitempty"`

	AlternativeText string `json:"alt,omitempty"`
}

func (comic *ComicBody) String() string {
	return fmt.Sprintf("Number: %d\nURL: %s\nTitle: %s\nDay: %s\nMonth: %s\nYear: %s\nAlternative Text: %s\nLink: %s\nImageURL: %s\nNews: %s\n\n",
		comic.Number,
		getRawURL(),
		comic.Title,
		comic.Day,
		comic.Month,
		comic.Year,
		comic.AlternativeText,
		comic.Link,
		comic.ImageURL,
		comic.News,
	)
}

type Empty struct{}

var (
	Logger *log.Logger

	methods map[string]Empty

	fileName, method *string
	comicNumber      *int
)

func init() {
	Logger = lg.Init()

	methods = map[string]Empty{
		"write": {},
		"read":  {},
	}

	fileName = flag.String("filename", "comics.txt", "file to be read from/written into\nthe default value is comics.txt")
	method = flag.String("method", "", "method to be executed (write, read)\nthe default value is \"\"")
	comicNumber = flag.Int("comicnumber", -1, "number of comic you're gonna interact with\nthe default value is -1")
}

func Execute() {
	Logger.Println("Session started")
	defer Logger.Printf("Session closed\n\n")

	flag.Parse()
	Logger.Printf("USER DATA: \nfileName: %s\nmethod: %s\ncomicNumber: %d\n", *fileName, *method, comicNumber)

	err := validate()
	if err != nil {
		Logger.Println(err.Error())
		return
	}

	switch *method {
	case "read":
		comicContentBuf, err := read()
		if err != nil {
			Logger.Println(err)
			fmt.Print(err)
			return
		}
		Logger.Println("successfully read")
		fmt.Println(comicContentBuf.String())

	case "write":
		err := write()
		if err != nil {
			Logger.Println(err)
			fmt.Print(err)
			return
		}
		Logger.Println("successfully written")

	}
}

/* Validations */
func validate() error {
	err := isConsoleParamsCountValid()
	if err != nil {
		return err

	}

	err = isMethodValid()
	if err != nil {
		return err
	}

	return nil
}

func isConsoleParamsCountValid() error {
	if len(os.Args[1:]) == cfg.ValidParamsCount {
		return nil
	}
	return errors.New("invalid valid params count")
}

func isMethodValid() error {
	for validMethod := range methods {
		if validMethod == *method {
			return nil
		}
	}
	return errors.New("no valid method")
}

/* Main methods */
func read() (*bytes.Buffer, error) {
	file, err := os.OpenFile(*fileName, os.O_RDONLY, os.ModeIrregular)
	if err != nil {
		return nil, errors.New("search file opening failed")
	}

	readFileBuf := bufio.NewScanner(file)
	readFileBuf.Split(bufio.ScanLines)

	comicContentBuffer := bytes.Buffer{}
	var readText string

	titleFound := false
	for readFileBuf.Scan() {
		readText = readFileBuf.Text()
		if !titleFound {
			if strings.HasPrefix(readText, cfg.ComicNumberToken+strconv.Itoa(*comicNumber)) {
				_, err = comicContentBuffer.WriteString(readFileBuf.Text())
				if err != nil {
					return nil, errors.New("writing failed")
				}
				_, err = comicContentBuffer.WriteRune('\n')
				if err != nil {
					return nil, errors.New("writing failed")
				}
				titleFound = true

			}
			continue
		}

		if !strings.HasPrefix(readText, "Number: ") {
			_, err = comicContentBuffer.WriteString(readText)
			if err != nil {
				return nil, errors.New("writing failed")
			}
			_, err = comicContentBuffer.WriteRune('\n')
			if err != nil {
				return nil, errors.New("writing failed")
			}
		} else {
			return &comicContentBuffer, nil
		}

	}

	return nil, errors.New("no required comic")
}

func write() error {

	err := search()
	if err != nil {
		return err
	}

	request, err := http.NewRequest("GET", getJsonURL(), nil)
	if err != nil {
		return err
	}
	request.Header.Set("Timeout", "10")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}

	var comicBodyUnmarshaled ComicBody
	err = json.NewDecoder(response.Body).Decode(&comicBodyUnmarshaled)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(*fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return errors.New("logger creation fail")
	}

	_, err = file.WriteString(comicBodyUnmarshaled.String())
	if err != nil {
		return errors.New("logger creation fail")
	}

	return nil
}

func search() error {
	file, err := os.OpenFile(*fileName, os.O_CREATE|os.O_RDWR, os.ModeIrregular)
	if err != nil {
		return err
	}
	defer file.Close()

	readFileBuf := bufio.NewScanner(file)
	readFileBuf.Split(bufio.ScanLines)
	var readText string

	for readFileBuf.Scan() {
		readText = readFileBuf.Text()
		if strings.HasPrefix(readText, cfg.ComicNumberToken+strconv.Itoa(*comicNumber)) {
			return errors.New("comic already exists")
		}
	}

	return nil
}

func getJsonURL() string {
	return fmt.Sprintf(cfg.JsonURL, strconv.Itoa(*comicNumber))
}

func getRawURL() string {
	return fmt.Sprintf(cfg.RawURL, strconv.Itoa(*comicNumber))
}
