package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
    "strings"
    "encoding/json"
    "fmt"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
    "github.com/google/uuid"
    "github.com/typegaro/HamstersTunnel/pkg/models/service"
)

func TestHandlerNewService(t *testing.T){
    e := echo.New()
    s := NewServiceManager()

	
    saveValue := "false"
    url := fmt.Sprintf("/service?save=%s", saveValue)
    reqBody := models.NewServiceReq{
        Name:          "TestService",
        TCP:           true,
        UDP:           false,
        HTTP:          true,
        PortBlackList: []string{"22", "80", "443"},
    }

    jsonBody,err := json.Marshal(reqBody);
    assert.NoError(t, err, "Errore nella conversione JSON")

    req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(string(jsonBody)))
    req.Header.Set("Content-Type", "application/json")
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    if assert.NoError(t, s.HandlerNewService(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)

        var response map[string]string
        err := json.Unmarshal(rec.Body.Bytes(), &response)
        assert.NoError(t, err)

        serviceID, exists := response["service_id"]
        assert.True(t, exists, "service_id non presente nella risposta")

        _, err = uuid.Parse(serviceID)
        assert.NoError(t, err, "service_id non è un UUID valido")
    }
}
