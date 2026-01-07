package routes

import (
	"asphalt9db/utils"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	utils.RenderHtml(w, "home.html", nil)
}
