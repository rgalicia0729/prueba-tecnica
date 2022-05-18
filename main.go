package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type MyStoredVariable struct {
	Id    string
	Url   string
	Value string
}

var (
	myStoredVariableList []*MyStoredVariable
)

func requestHttp(url string) (*MyStoredVariable, error) {
	clientHttp := &http.Client{}

	peticion, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	peticion.Header.Add("Content-Type", "application/json")
	respuesta, err := clientHttp.Do(peticion)
	if err != nil {
		return nil, err
	}

	defer respuesta.Body.Close()

	cuerpoRespuesta, err := ioutil.ReadAll(respuesta.Body)
	if err != nil {
		return nil, err
	}

	respuestaString := string(cuerpoRespuesta)

	var myStoredVariable MyStoredVariable
	if respuesta.StatusCode == 200 {
		json.Unmarshal([]byte(respuestaString), &myStoredVariable)

	} else {
		return nil, err
	}

	return &myStoredVariable, nil
}

func validResp(id string) bool {
	var resp bool

	for _, value := range myStoredVariableList {
		if value.Id == id {
			resp = true
			break
		}
	}

	return resp
}

func getResults(wg *sync.WaitGroup) {
	defer wg.Done()

	response, err := requestHttp("https://api.chucknorris.io/jokes/random")
	if err != nil {
		log.Print(err)
	}

	if !validResp(response.Id) {
		myStoredVariableList = append(myStoredVariableList, response)
	}
}

func main() {
	e := echo.New()

	e.GET("/api/chucknorris", func(c echo.Context) (err error) {
		myStoredVariableList = nil

		var wg sync.WaitGroup
		for i := 0; i < 15; i++ {
			wg.Add(1)
			go getResults(&wg)
		}
		wg.Wait()

		return c.JSON(http.StatusOK, myStoredVariableList)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
