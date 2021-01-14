package files

import "strings"

//GetPathLevels returns array with the branches name of the path tree
func GetPathLevels(path string) []string {
	branches := []string{}
	if !strings.HasPrefix(path, "/") {
		branches = append(branches, "/")
	}

	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	branches = append(branches, strings.Split(path, "/")...)
	branches[0] = "/"
	return branches
}
