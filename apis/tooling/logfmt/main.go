package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	var b strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Bytes()
		m := make(map[string]any)
		err := json.Unmarshal(s, &m)
		if err != nil {
			continue
		}
		traceID := "00000000-0000-0000-0000-000000000000"
		if v, ok := m["trace_id"]; ok {
			traceID = fmt.Sprintf("%v", v)
		}
		b.Reset()
		b.WriteString(fmt.Sprintf("%s: %s: %s: %s: %s: %s: ",
			m["service"],
			m["time"],
			m["file"],
			m["level"],
			traceID,
			m["msg"],
		))
		for k, v := range m {
			if k == "service" || k == "time" || k == "file" || k == "level" || k == "trace_id" || k == "msg" {
				continue
			}
			b.WriteString(fmt.Sprintf("%s=[%v]: ", k, v))
		}
		out := b.String()
		fmt.Println(out[:len(out)-2])
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
