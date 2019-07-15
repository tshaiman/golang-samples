package urlshort

import (
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}

}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parseUrls, err := parseYaml(yml)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
		return nil, err
	}
	pathsToUrls := buildMap(parseUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func buildMap(parseUrls []parseUrl) map[string]string {
	pathToUrls := make(map[string]string)
	for _, pu := range parseUrls {
		pathToUrls[pu.Path] = pu.Url
	}
	return pathToUrls
}

func parseYaml(yml []byte) ([]parseUrl, error) {
	var parseUrls []parseUrl
	err := yaml.Unmarshal([]byte(yml), &parseUrls)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}
	return parseUrls, err
}

type parseUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}