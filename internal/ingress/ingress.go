package ingress

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	amTemplate "github.com/prometheus/alertmanager/template"
)

type Ingress struct {
	Path string
}

func GetAmData(req *http.Request, l *log.Logger) (*amTemplate.Data, error) {
	json, err := getBody(req, l)
	if err != nil {
		return nil, err
	}
	amData, err := unmarshalBody(json, l)
	if err != nil {
		return nil, err
	}
	return amData, nil
}

func getBody(req *http.Request, l *log.Logger) ([]byte, error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		l.Printf("ERROR: ReadAll(body) %s\n", err)
		return nil, err
	}
	return data, nil
}

func unmarshalBody(amJSON []byte, l *log.Logger) (*amTemplate.Data, error) {

	var amData amTemplate.Data

	err := json.Unmarshal(amJSON, &amData)
	if err != nil {
		l.Fatalf("ERROR: Unmarshall recv JSON %s\n", err)
		return nil, err
	}
	l.Printf("Unmarshalled json from alertmanager: %#v\n", amData)
	l.Printf("Alert #1: %#v\n", amData.Alerts[0])

	return &amData, nil
}
