package pdf

import (
	gofpdf "github.com/jung-kurt/gofpdf"
)

func CreateFile() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	pdf.OutputFileAndClose("hello.pdf")
}
