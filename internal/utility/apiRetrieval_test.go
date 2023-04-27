package utility

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/mock"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCountryFromAPI(t *testing.T) {
	// Changes the working directory to the project directory.
	err := DirChanger(2)
	if err != nil {
		return
	}

	// Testing a get request on local host
	getRequest, _ := http.NewRequest("GET", constants.MOCK_API_URL, nil)
	response := httptest.NewRecorder()
	//Executing the handler
	mock.StubAPI(response, getRequest)

}

func TestMockGetCountryFromAPI(t *testing.T) {
	dirChangeErr := DirChanger(2)
	if dirChangeErr != nil {
		t.Fatal("Error changing directories.")
	}
	/*expectedCountry := structs.Country{
		Name:        map[string]interface{}{"common": "Norway"},
		CountryCode: "NOR",
		Borders:     []string{"SWE", "FIN", "RUS"},
		Cache:       time.Time{},
	}*/

	// Create a mock HTTP server using the StubAPI handler
	mockServer := httptest.NewServer(http.HandlerFunc(mock.StubAPI))
	resp, err := http.Get(mockServer.URL)
	if err != nil {
		t.Fatal("Error")
	} else if resp.StatusCode == http.StatusInternalServerError {
		t.Fatal("Country issues.")
	}
	fmt.Println(resp)
}
