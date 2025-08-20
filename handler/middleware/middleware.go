package middleware

import (
	"fmt"
	"net/http"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// IsAuthenticated is a middleware that checks if
// the user has already been authenticated previously.
/*
func IsAuthenticated(ctx *gin.Context) {
	fmt.Println("auth method")
	if sessions.Default(ctx).Get("profile") == nil {
		ctx.Redirect(http.StatusSeeOther, "/")
	} else {
		ctx.Next()
	}
}
*/

	func IsAuthenticated(ctx *gin.Context) {
    fmt.Println("auth method")
    session := sessions.Default(ctx)

    if session.Get("profile") == nil {
        // For APIs: send 401
        ctx.AbortWithStatus(http.StatusUnauthorized)
        return

        // For browser routes: redirect + abort
        // ctx.Redirect(http.StatusSeeOther, "/")
        // ctx.Abort()
        // return
    }

    ctx.Next()
}
