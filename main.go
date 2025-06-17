// package main

// import (
// 	// "github/nawfay/didban/internal/logic"
// 	// "fmt"

// 	"github/nawfay/didban/internal/downloader"
// )


// func main() {

// 	// tmp, _ := logic.DeezerToYtResolver(3135556)
// 	// fmt.Println(tmp)
	
// 	downloader.ExampleClient()
// }

package main

import (
	"context"
	"log"
	"time"

	deezer"github/nawfay/didban/internal/downloader"
)

func main() {
	// Make sure DEEZER_ARL is exported in your shell:
	// export DEEZER_ARL="your_arl_cookie_here"

	// Download in a 5 min context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := deezer.DownloadTrack(ctx, "2864040362", "320", "mytrack.mp3"); err != nil {
		log.Fatal(err)
	}
	log.Println("Download complete!")
}
