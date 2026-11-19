package main

import "fmt"

// 1. Interface Kecil (Single-Method Interfaces)
type Reader interface {
	Read() string
}

type Writer interface {
	Write(data string)
}

// 2. Interface Embedding: Menggabungkan Reader dan Writer
type ReadWriter interface {
	Reader
	Writer
}

// 3. Concrete Struct
type FileSystem struct {
	content string
}

func (fs *FileSystem) Read() string {
	return fs.content
}

func (fs *FileSystem) Write(data string) {
	fs.content = data
}

// 4. "Accept Interface": Fungsi menerima interface yang fleksibel
func CopyContent(w Writer, r Reader) {
	text := r.Read()
	w.Write("Disalin: " + text)
}

func main() {
	source := &FileSystem{content: "Data Rahasia Finansial"}
	destination := &FileSystem{}

	// source bertindak sebagai Reader, destination bertindak sebagai Writer:
	CopyContent(destination, source)

	fmt.Println("Hasil setelah disalin:", destination.Read())
}
