package db

import (
	"fmt"
)

const URLEndPoint = "https://ik.imagekit.io/apeironarchives"

func BuildRawFromId(id string, isportrait bool) string {
	return fmt.Sprintf("%s/raw/%s.jpg", URLEndPoint, id)
}

func BuildPreviewFromId(id string, isportrait bool) string {
	var tr string
	if isportrait {
		tr = "tr:h-640"
	} else {
		tr = "tr:w-640"
	}
	return fmt.Sprintf("%s/%s/raw/%s.jpg", URLEndPoint, tr, id)
}
