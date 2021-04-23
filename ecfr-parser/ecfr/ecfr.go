package ecfr

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const dateFormat = "2006-01-02"

var (
	ecfrSite         = urlMustParse("https://ecfr.federalregister.gov/api/versioner/v1/")
	ecfrFullXML      = "full/%s/title-%d.xml"
	ecfrVersionsXML  = "versions/title-%d.json"
	ecfrStructureXML = "structure/%s/title-%d.json"
)

func urlMustParse(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	return u
}

func buildQuery(opts []FetchOption) string {
	q := url.Values{}
	for _, opt := range opts {
		v := opt.Values()
		for key, _ := range v {
			q.Set(key, v.Get(key))
		}
	}
	return q.Encode()
}

func FetchFull(date string, title int, opts ...FetchOption) (io.ReadCloser, error) {
	path, err := url.Parse(fmt.Sprintf(ecfrFullXML, date, title))
	if err != nil {
		return nil, err
	}

	path.RawQuery = buildQuery(opts)

	u := ecfrSite.ResolveReference(path)

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%d", resp.StatusCode)
	}

	return resp.Body, nil
}

func FetchStructure(date string, title int, opts ...FetchOption) (io.ReadCloser, error) {
	path, err := url.Parse(fmt.Sprintf(ecfrStructureXML, date, title))
	if err != nil {
		return nil, err
	}

	path.RawQuery = buildQuery(opts)

	u := ecfrSite.ResolveReference(path)

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Println(u.String())
		return nil, fmt.Errorf("%d", resp.StatusCode)
	}

	return resp.Body, nil
}

func FetchVersions(title int, opts ...FetchOption) (io.ReadCloser, error) {
	path, err := url.Parse(fmt.Sprintf(ecfrVersionsXML, title))
	if err != nil {
		return nil, err
	}

	path.RawQuery = buildQuery(opts)

	u := ecfrSite.ResolveReference(path)

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

type FetchOption interface {
	Values() url.Values
}

type partOption struct {
	part string
}

func (p *partOption) Values() url.Values {
	v := url.Values{}
	v.Set("part", p.part)
	return v
}

func Part(part string) FetchOption {
	return &partOption{
		part: part,
	}
}

type subchapterOption struct {
	chapter    string
	subchapter string
}

func (p *subchapterOption) Values() url.Values {
	v := url.Values{}
	v.Set("chapter", p.subchapter)
	v.Set("subchapter", p.chapter)
	return v
}

func Subchapter(chapter string, subchapter string) FetchOption {
	return &subchapterOption{
		chapter:    chapter,
		subchapter: subchapter,
	}
}
