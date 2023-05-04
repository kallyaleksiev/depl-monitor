package controllers

import "fmt"

// Get the derived name of an underlying from the
// name of the MonDepl
func GetUnderlyingName(name string) string {
	return fmt.Sprintf("spicy-%s")
}
