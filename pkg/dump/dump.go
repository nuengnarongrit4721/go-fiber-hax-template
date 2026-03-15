package dump

import (
	"encoding/json"
	"fmt"
)

func Print(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Dump Error:", err)
		return
	}
	fmt.Println(string(b))
}
