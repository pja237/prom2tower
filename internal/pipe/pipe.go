package pipe

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/pja237/prom2tower/internal/egress"
	"github.com/pja237/prom2tower/internal/ingress"
)

type Pipe struct {
	Name    string
	Ingress ingress.Ingress
	Egress  egress.Egress
}

type PipeWorker struct {
	Id   int
	Err  error
	pipe *Pipe
	log  *log.Logger
}

func returnString(v interface{}) *string {
	if s, ok := v.(string); ok {
		return &s
	} else {
		return nil
	}
}

func NewPipeWorker(id int, g map[string]interface{}, p Pipe, l *log.Logger) (*PipeWorker, error) {
	// todo: here we check the correctness of config
	// todo: also send "globals" config here so that config like "towerHost" can be overriden
	//       or set on global level only and propagated to pipe here
	var newp = PipeWorker{
		Id:   id,
		pipe: &p,
		log:  l,
	}

	if newp.pipe.Egress.TowerHost == "" {
		l.Printf("%s towerHost not defined\n", newp.pipe.Name)
		if v, ok := g["towerHost"]; !ok {
			l.Printf("%s towerHost not defined globally, exit!\n", newp.pipe.Name)
			return nil, errors.New("towerHost not defined")
		} else {
			l.Printf("%s using global towerHost %s\n", newp.pipe.Name, g["towerHost"])
			newp.pipe.Egress.TowerHost = *returnString(v)
		}
	}
	if newp.pipe.Egress.TowerToken == "" {
		l.Printf("%s towerToken not defined\n", newp.pipe.Name)
		if v, ok := g["towerToken"]; !ok {
			l.Printf("%s towerToken not defined globally, exit!\n", newp.pipe.Name)
			return nil, errors.New("towerToken not defined")
		} else {
			l.Printf("%s using global towerToken %s\n", newp.pipe.Name, g["towerToken"])
			newp.pipe.Egress.TowerToken = *returnString(v)
		}
	}

	switch {
	case newp.pipe.Egress.TowerHost == "":
		return nil, errors.New("towerHost not defined")
	case newp.pipe.Egress.TowerToken == "":
		return nil, errors.New("towerToken not defined")
	case newp.pipe.Egress.Path == "":
		return nil, errors.New("path not defined")
	case newp.pipe.Egress.Body == "":
		return nil, errors.New("body not defined")
	case newp.pipe.Egress.Method == "":
		return nil, errors.New("method not defined")
	}

	return &newp, nil
}

func (pw PipeWorker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	pw.log.Printf("Serving request from pipe %d %s\n", pw.Id, pw.pipe.Name)
	pw.log.Printf("Incoming request: %s\n", req.Method)
	pw.log.Printf("Incoming request header: %v\n", req.Header)

	// Get unmarshalled JSON body from alertmanager into ingressData structure (from alertmanager/template)
	iData, err := ingress.GetAmData(req, pw.log)
	if err != nil {
		pw.log.Printf("ERROR: GetAmData(req) %s\n", err)
	}

	// Template egress json body using incoming iData
	eJSON, err := egress.PrepareEJSON(pw.pipe.Egress.Body, iData, pw.log)
	if err != nil {
		pw.log.Fatalf("ERROR: PrepareEgressJSON %s\n", err)
	}

	// prepare egress request and send
	sendBody := strings.NewReader(eJSON.String())
	client := http.Client{}
	reqCli, err := http.NewRequest(pw.pipe.Egress.Method, pw.pipe.Egress.TowerHost+pw.pipe.Egress.Path, sendBody)
	if err != nil {
		pw.log.Printf("NewRequest() failed with %s\n", err)
	} else {
		// add headers and send
		// todo: customizable headers from config
		reqCli.Header.Add("Authorization", "Bearer "+pw.pipe.Egress.TowerToken)
		reqCli.Header.Add("Content-Type", "application/json")
		if resp, err := client.Do(reqCli); err != nil {
			pw.log.Printf("client.Do(req) failed with %s\n", err)
		} else {
			pw.log.Printf("Response: %#v\n", resp.Status)
		}
	}

}

func (pw *PipeWorker) SpinUp(wg *sync.WaitGroup) error {
	defer wg.Done()

	pw.log.Printf("Registering pipeWorker %d: %s on %s\n", pw.Id, pw.pipe.Name, pw.pipe.Ingress.Path)
	http.Handle(pw.pipe.Ingress.Path, pw)

	return nil
}
