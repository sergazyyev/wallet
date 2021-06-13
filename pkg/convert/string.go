package convert

import "fmt"

// String converting any value to string
func String(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
