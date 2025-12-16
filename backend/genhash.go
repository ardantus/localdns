package main
import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)
func main() {
    hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 14) // Password: admin123
    fmt.Println(string(hash))
}
