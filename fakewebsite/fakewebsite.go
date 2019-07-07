package fakewebsite

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const Sadpanda = `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
</head>
<body>
<img src="sadpanda.jpg">
</body>
</html>
`

func SendSadPanda(w http.ResponseWriter, r *http.Request) {
	/*
		sadpanda_img, err := ioutil.ReadFile("sadpanda.jpg")
		if err != nil {
			panic(err)
		}
		w.Write(sadpanda_img)
	*/
	// https://stackoverflow.com/questions/26744814/serve-image-in-go-that-was-just-created
	fmt.Println("try to send back a image")

	Path := "static/sadpanda.jpg"
	img, err := os.Open(Path)
	if err != nil {
		log.Fatal(err) // perhaps handle this nicer
	}
	defer img.Close()
	w.Header().Set("Content-Type", "image/jpeg") // <-- set the content-type header
	io.Copy(w, img)

}
