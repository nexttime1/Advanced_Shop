package user

import (
	"math/rand"
	"os"
	"runtime"
	"time"
)

func main() {
	source := rand.NewSource(time.Now().UnixNano())
	rand.New(source)
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	srv.NewApp("order-server").Run()
}
