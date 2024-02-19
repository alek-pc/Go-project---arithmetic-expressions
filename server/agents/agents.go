package agents

import (
	"fmt"
	"html/template"
	"net/http"
	set "arifm_operations/server/settings"
)
type agents struct{
	Busy bool
	Busy_agents int
}
type response struct{
	Agents []agents
}
func AgentsHandler(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("./templates/agents_page.html")  // шаблонизатор страницы
	if err != nil{
		http.Error(w, "parsing page", 500)
		fmt.Println(err)
		return
	}
	resp := response{Agents: make([]agents, 0)}
	for _, agent := range set.Agents.Agents{
		resp.Agents = append(resp.Agents, agents{Busy: agent.Busy, Busy_agents: agent.MaxNum - agent.CurNum})
	}

	tmpl.Execute(w, resp)
}