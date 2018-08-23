package middleware

import (
	"fmt"

	"github.com/bluele/gcache"
	"github.com/gin-gonic/gin"
)

// GCache gcache
func GcAche(cache gcache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, err := c.GetRawData()
		if err != nil {
			c.Next()
			return
		}

		// GetRawData buffed
		if len(raw) == 0 {
			raw = c.MustGet("raw").([]byte)
		} else {
			c.Set("raw", raw)
		}

		fmt.Println("++++++++++++: ", raw)

		c.Next()
	}
}
