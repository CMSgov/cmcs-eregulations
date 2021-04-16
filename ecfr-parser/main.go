package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/cmsgov/cmcs-eregulations/ecfr-parser/ecfr"
)

type Part struct {
	XMLName   xml.Name     `xml:"DIV5"`
	Name      string       `xml:"N,attr"`
	Type      string       `xml:"TYPE,attr"`
	Header    string       `xml:"HEAD"`
	Authority string       `xml:"AUTH>PSPACE"`
	Source    string       `xml:"SOURCE>PSPACE"`
	Children  PartChildren `xml:",any"`
}

type Subpart struct {
	Header   string          `xml:"HEAD"`
	Type     string          `xml:"TYPE,attr"`
	Children SubpartChildren `xml:",any"`
}

type SubjectGroup struct {
	Header string `xml:"HEAD"`
	Type   string          `xml:"TYPE,attr"`
	Children SubjectGroupChildren `xml:",any"`
}

type Section struct {
	Type string `xml:"TYPE,attr"`
	Header   string          `xml:"HEAD"`
	Children SectionChildren `xml:",any"`
}

type Paragraph struct {
	Type string
	Content string `xml:",innerxml"`
}

type Extract struct {
	Type string
	Content string `xml:",innerxml"`
}

type Citation struct {
	Type string
	Content string `xml:",innerxml"`
}

type PartChildren []interface{}
type SubpartChildren []interface{}
type SectionChildren []interface{}
type SubjectGroupChildren []interface{}

func (c *PartChildren) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	switch start.Name.Local {
	case "DIV6":
		child := &Subpart{}
		if err := d.DecodeElement(child, &start); err != nil {
			return err
		}
		*c = append(*c, child)
	case "DIV8":
		child := &Section{}
		if err := d.DecodeElement(child, &start); err != nil {
			return err
		}
		*c = append(*c, child)
	default:
		//return fmt.Errorf("Unknown XML type in Part: %+v", start)
		fmt.Printf("Unknown XML type in Part: %+v\n", start)
		d.Skip()
	}

	return nil
}

func (c *SubpartChildren) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	switch start.Name.Local {
	case "DIV8":
		child := &Section{}
		if err := d.DecodeElement(child, &start); err != nil {
			return err
		}
		*c = append(*c, child)
	case "DIV7":
		child := &SubjectGroup{}
		if err := d.DecodeElement(child, &start); err != nil {
			return err
		}
		*c = append(*c, child)
	default:
		//return fmt.Errorf("Unknown XML type in Subpart: %+v", start)
		fmt.Printf("Unknown XML type in Subpart: %+v\n", start)
		d.Skip()
	}
	return nil
}

func (c *SubjectGroupChildren) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	switch start.Name.Local {
	case "DIV8":
		child := &Section{}
		if err := d.DecodeElement(child, &start); err != nil {
			return err
		}
		*c = append(*c, child)
	default:
		//return fmt.Errorf("Unknown XML type in Subpart: %+v", start)
		fmt.Printf("Unknown XML type in Subject Group: %+v\n", start)
		d.Skip()
	}
	return nil
}

func (c *SectionChildren) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	switch start.Name.Local {
	case "P":
		child := &Paragraph{ Type: "Paragraph"}
		if err := d.DecodeElement(child, &start); err != nil {
			return err
		}
		*c = append(*c, child)
	case "EXTRACT":
		child := &Extract{ Type: "Extract" }
		if err := d.DecodeElement(child, &start); err != nil {
			return err
		}
		*c = append(*c, child)
	case "CITA":
		child := &Citation{ Type: "Citation" }
		if err := d.DecodeElement(child, &start); err != nil {
			return err
		}
		*c = append(*c, child)
	default:
		//return fmt.Errorf("Unknown XML type in Section: %+v", start)
		fmt.Printf("Unknown XML type in Section: %+v\n", start)
		d.Skip()
	}

	return nil
}

func ParsePart(b io.ReadCloser) (*Part, error) {
	d := xml.NewDecoder(b)

	p := &Part{}
	log.Println("[DEBUG] Decoding part")
	if err := d.Decode(p); err != nil {
		return p, err
	}
	return p, nil
}

func main() {
	today := time.Now()

	start := time.Now()
	log.Println("[DEBUG] fetching part")
	body, err := ecfr.FetchPart(today, 42, ecfr.Part(433))
	if err != nil {
		log.Fatal(err)
	}
	defer body.Close()

	p, err := ParsePart(body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[DEBUG] Marshaling JSON of Part")
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(p); err != nil {
		log.Fatal(err)
	}

	log.Println("Run time:", time.Since(start))
}