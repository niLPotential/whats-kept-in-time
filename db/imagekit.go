package db

import (
	"fmt"
)

const URLEndPoint = "https://ik.imagekit.io/apeironarchives"

func BuildURLFromId(id string) string {
	return fmt.Sprintf("%s/raw/%s.jpg", URLEndPoint, id)
}
