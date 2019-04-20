package handlers

import "net/http"

//Index - перенаправление на нужную страницу
func Index(w http.ResponseWriter, r *http.Request) {
	if !Ready {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/stream", http.StatusSeeOther)
	}
}
