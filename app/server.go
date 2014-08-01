package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "strings"
)

var teams []Team
var projects []Project

func main() {
  loadTeams()
  loadProjects()
  // TODO: validate project names are unique
  // TODO: validate every project.Team matches a team.Name
  http.HandleFunc("/", viewHandler)
  http.ListenAndServe(":8080", nil)
}

func loadTeams() () {
  path:= strings.Join([]string{os.Getenv("DATA_DIR"), "/teams.json"}, "")
  data, err := ioutil.ReadFile(path)
  if err == nil && data != nil {
      err = json.Unmarshal(data, &teams)
  }
  if err != nil {
    fmt.Println("Error - Failed to load teams via", path)
  }
}

func loadProjects() () {
  path:= strings.Join([]string{os.Getenv("DATA_DIR"), "/projects.json"}, "")
  data, err := ioutil.ReadFile(path)
  if err == nil && data != nil {
      err = json.Unmarshal(data, &projects)
  }
  if err != nil {
    fmt.Println("Error - Failed to load projects via", path)
  }
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Content-type", "application/json")

  jsonMsg, err := getResponse(r.URL.Path[1:])
  if err != nil {
    http.Error(w, "Oops", http.StatusInternalServerError)
  }
  fmt.Fprintf(w, jsonMsg)
}

func getResponse(query string) (string, error){
  log.Println("Looking up owner of :", query)
  project, err := getProjectByName(query)
  if err != nil {
    return "", err
  }

  team, err := getTeamByName(project.Team)
  if err != nil {
    return "", err
  }

  ownership := Ownership{project.Name, team.Name, team.Email, team.Members, project.Aliases}

  jsonData, err := json.Marshal(ownership)
  if err != nil {
    return "", err
  }

  jsonMsg := string(jsonData[:]) // converting byte array to string
  return jsonMsg, nil
}

func getProjectByName(name string) (*Project, error) {
  lcName := strings.ToLower(name)
  for _, p := range projects {
    if strings.ToLower(p.Name) == lcName {
      return &p, nil
    }
  }
  return getProjectByAlias(lcName)
}

// TODO: support shared aliases for ambigious keywords, e.g. "metrics"
func getProjectByAlias(name string) (*Project, error) {
  lcName := strings.ToLower(name)
  for _, p := range projects {
    for _, a := range p.Aliases {
      if strings.ToLower(a) == lcName {
        return &p, nil
      }
    }
  }
  return nil, fmt.Errorf("Project alias '%s' does not exist.", name)
}


func getTeamByName(name string) (*Team, error) {
  lcName := strings.ToLower(name)
  for _, t := range teams {
    if strings.ToLower(t.Name) == lcName {
      return &t, nil
    }
  }
  return nil, fmt.Errorf("Team '%s' does not exist.", name)
}

type Ownership struct {
  Project string
  Team string
  Email string
  Owners []string
  Aliases []string
}

type Project struct {
  Name string
  Aliases []string
  Team string
}

type Team struct {
  Name string
  Email string
  Members []string
}
