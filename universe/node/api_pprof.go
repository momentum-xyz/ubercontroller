package node

import (
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

// Expose the pprof http endpoints in Gin framework.
// See https://pkg.go.dev/runtime/pprof
func registerPProfAPI(g *gin.RouterGroup) {
	g.GET("/pprof", gin.WrapF(pprof.Index))
	g.GET("/cmdline", gin.WrapF(pprof.Cmdline))
	g.GET("/profile", gin.WrapF(pprof.Profile))
	g.GET("/symbol", gin.WrapF(pprof.Symbol))
	g.POST("/symbol", gin.WrapF(pprof.Symbol))
	g.GET("/trace", gin.WrapF(pprof.Trace))

	// https://pkg.go.dev/runtime/pprof#Profile
	for _, s := range []string{"allocs", "block", "goroutine", "heap", "mutex", "threadcreate"} {
		g.GET(s, gin.WrapH(pprof.Handler(s)))
	}
}
