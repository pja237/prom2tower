package egress

import (
	"bytes"
	"log"
	"text/template"

	amTemplate "github.com/prometheus/alertmanager/template"
)

type Egress struct {
	TowerHost  string `yaml:"towerHost"`
	Path       string `yaml:"path"`
	Method     string `yaml:"method"`
	Headers    string `yaml:"headers"`    // if external modification is needed at all?
	Body       string `yaml:"body"`       // make this a templateable json
	Options    string `yaml:"options"`    // no idea what i thought with this one?
	TowerToken string `yaml:"towerToken"` // oauth2 token
}

func PrepareEJSON(templ string, amData *amTemplate.Data, l *log.Logger) (*bytes.Buffer, error) {
	var eBody bytes.Buffer

	// todo: this can panic, handle it
	t := template.Must(template.New("TowerJSON").Parse(templ))
	err := t.Execute(&eBody, *amData)
	if err != nil {
		l.Printf("ERROR: Template egress.Body %s\n", err)
		return nil, err
	}
	return &eBody, nil
}
