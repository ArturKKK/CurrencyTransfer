package parser

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	CharCode  string `xml:"CharCode"`
	VunitRate string `xml:"VunitRate"`
}
