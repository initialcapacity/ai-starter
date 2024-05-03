package testsupport

import "fmt"

func Stream(json string) []byte {
	return []byte(fmt.Sprintf("data:%s\ndata:[DONE]", json))
}
