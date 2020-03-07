package reg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"mzy-spider/stock"
	"strconv"
)

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", value), 64)
	return value
}
func ServerMiddleware(auth bool) gin.HandlerFunc {
	Action := func(r *gin.Context) {
		runKey := r.Query("run_key")
		if auth {
			rely, err := redis.String(stock.Redis.Do("GET", "run_key"))
			if err != nil || runKey != rely {
				r.String(200, "INVALID PARAMETER VALUE")
				r.Abort()
				return
			}
		}
	}
	return Action
}
