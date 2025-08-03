package logout

import (
	"net/http"
	"net/url"
	"github.com/gin-gonic/gin"
	"github.com/laurisseau/sportsify-config"
	"github.com/gin-contrib/sessions"
)

// Handler for our logout.
func Handler(ctx *gin.Context) {

    config.LoadSecretsEnv()

	logoutUrl, err := url.Parse("https://" + config.Secrets["AUTH0_DOMAIN"] + "/v2/logout")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + ctx.Request.Host)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", config.Secrets["AUTH0_CLIENT_ID"])
	logoutUrl.RawQuery = parameters.Encode()

	session := sessions.Default(ctx)
	session.Clear()
	session.Save()

	ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}