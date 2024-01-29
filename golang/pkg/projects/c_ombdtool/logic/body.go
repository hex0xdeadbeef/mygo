package issuetool

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	cfg "golang/pkg/projects/c_ombdtool/config"
	lgr "golang/pkg/projects/c_ombdtool/logger"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	Logger                   *log.Logger
	key, method, methodValue *string
	methods                  map[string]Emtpy
)

func init() {
	// Logger initialization
	{
		Logger = lgr.LoggerInit()
	}

	// Required params flags initialization
	{
		key = flag.String("key", "", "your ombd api auth key")
		method = flag.String("method", "", "method (id, title)")
		methodValue = flag.String("methodvalue", "", "e.g. tt1285016 or god")
	}

	// Validations
	methods = map[string]Emtpy{
		"id":    {},
		"title": {},
	}
}

/* Controller */
func Execute() {
	Logger.Println("Session created")
	defer Logger.Println("Session closed")

	flag.Parse()

	// Required params initialization
	if len(os.Args[1:]) > cfg.ValidParamsCount {
		Logger.Println(errors.New("invalid parameters count").Error())
	}
	Logger.Printf("USER DATA: \nkey: %s\nmethod: %s\nmethod value: %s\n", *key, *method, *methodValue)

	// Validations
	err := paramsValidation()
	if err != nil {
		Logger.Println(err.Error())
		return
	}
	err = methodValidation()
	if err != nil {
		Logger.Println(err.Error())
		return
	}

	err = getFilmInfo()
	if err != nil {
		fmt.Println(err)
		Logger.Println(err)
	}
}

/* Main methods */

func getFilmInfo() error {
	// Request creation
	request, err := http.NewRequest("GET", getURL(), nil)
	if err != nil {
		return err
	}

	// Do a request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}

	// Write response data
	var unmarshalledResponse FilmInfo
	if err := json.NewDecoder(response.Body).Decode(&unmarshalledResponse); err != nil {
		return err
	}

	fmt.Printf("Film info:\nTitle: %s\nYear: %s\nReleased: %s\nRuntime: %s\nGenre:%s\n\n",
		unmarshalledResponse.Title,
		unmarshalledResponse.Year,
		unmarshalledResponse.Released,
		unmarshalledResponse.Runtime,
		unmarshalledResponse.Genre,
	)
	Logger.Println("info is written")

	err = savePoster(unmarshalledResponse.PosterURL)
	if err != nil {
		return err
	}
	Logger.Println("poster is saved")
	fmt.Println("Poster successfully saved!")
	return nil
}

func savePoster(imageURL string) error {
	// Request creation
	request, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return err
	}

	// Do a request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}

	file, err := os.Create("poster.jpg")
	if err != nil {
		return err
	}
	defer file.Close()

	io.Copy(file, response.Body)

	return nil
}

func methodValidation() error {
	for existingMethod := range methods {
		if existingMethod == *method {
			return nil
		}
	}

	return errors.New("invalid method parameter")
}

func paramsValidation() error {
	if *key == "" || *method == "" || *methodValue == "" {

		return errors.New("required parameter empty")
	}
	return nil
}

func getURL() string {
	if *method == "id" {
		return fmt.Sprintf(cfg.APIURL, *key) + "i=" + *methodValue
	}
	return fmt.Sprintf(cfg.APIURL, *key) + "t=" + *methodValue
}
