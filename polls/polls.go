package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"gopkg.in/yaml.v2"
)

const questionsYaml string = "questions.yaml"
const locationAnsweres string = "answeres"

var templates = template.Must(template.ParseFiles("html/main.html", "html/view.html", "html/results.html", "html/header.html", "html/footer.html"))

func renderTemplate(w http.ResponseWriter, template string, params map[interface{}]interface{}) {
	err := templates.ExecuteTemplate(w, template+".html", params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loadPollQuestions() map[interface{}]interface{} {
	body, err := ioutil.ReadFile(questionsYaml)
	if err != nil {
		fmt.Println("Could not read questions for the poll", err)
		os.Exit(1)
	}
	poll := make(map[interface{}]interface{})
	err = yaml.Unmarshal(body, &poll)
	if err != nil {
		fmt.Println("Could not load ymal file", err)
		os.Exit(1)
	}
	return poll
}

func saveResults(r *http.Request) error {
	answeres := make(map[string]string)

	err := r.ParseForm()
	if err != nil {
		return err
	}

	for key, value := range r.Form {
		answeres[key] = value[0]
	}

	answeresYaml, err := yaml.Marshal(answeres)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%d.json", time.Now().Unix())
	fmt.Println("Writing poll results to", filename)
	err = ioutil.WriteFile(filepath.Join(locationAnsweres, filename), answeresYaml, 0755)
	return err
}

func handlerShowPollResults(w http.ResponseWriter, r *http.Request) {
	type userPollResult map[string]string
	var allPollAnsweres []userPollResult
	files, err := ioutil.ReadDir(locationAnsweres)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	for _, file := range files {
		data, err := ioutil.ReadFile(filepath.Join(locationAnsweres, file.Name()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		userAnsweres := make(userPollResult)
		err = yaml.Unmarshal(data, &userAnsweres)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		allPollAnsweres = append(allPollAnsweres, userAnsweres)
	}

	templateVars := make(map[interface{}]interface{})
	templateVars["q"] = loadPollQuestions()
	templateVars["a"] = allPollAnsweres
	err = templates.ExecuteTemplate(w, "results.html", templateVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func handlerFillPoll(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		err := saveResults(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte("Thank You"))

	} else {
		renderTemplate(w, "view", loadPollQuestions())
	}
}

func handlerMain(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "main", nil)
}

func main() {
	http.HandleFunc("/poll/", handlerMain)
	http.HandleFunc("/poll/fill", handlerFillPoll)
	http.HandleFunc("/poll/list", handlerShowPollResults)
	fmt.Println("Let's listen")
	log.Fatal(http.ListenAndServe("127.0.0.1:5000", nil))
}
