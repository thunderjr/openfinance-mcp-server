package uri

import (
	"fmt"
	"regexp"
	"strings"
)

type ResourceURI struct {
	Full        string
	Scheme      string
	Path        string
	PathParams  map[string]string
	QueryParams map[string][]string
}

func ParseURI(uri string) (*ResourceURI, error) {
	if uri == "" {
		return nil, fmt.Errorf("empty URI provided")
	}

	parts := strings.SplitN(uri, "://", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid URI format: %s", uri)
	}

	scheme := parts[0]
	remainder := parts[1]

	pathAndQuery := strings.SplitN(remainder, "?", 2)
	path := pathAndQuery[0]

	resourceUri := &ResourceURI{
		Full:        uri,
		Scheme:      scheme,
		Path:        path,
		PathParams:  make(map[string]string),
		QueryParams: make(map[string][]string),
	}

	if len(pathAndQuery) > 1 {
		queryString := pathAndQuery[1]
		if queryString != "" {
			queryParts := strings.Split(queryString, "&")
			for _, part := range queryParts {
				if part == "" {
					continue
				}

				keyValue := strings.SplitN(part, "=", 2)
				key := keyValue[0]

				value := ""
				if len(keyValue) > 1 {
					value = keyValue[1]
				}

				resourceUri.QueryParams[key] = append(resourceUri.QueryParams[key], value)
			}
		}
	}

	return resourceUri, nil
}

func Match(templateURI, concreteURI string) (*ResourceURI, error) {
	tResource, err := ParseURI(templateURI)
	if err != nil {
		return nil, fmt.Errorf("invalid template URI: %w", err)
	}

	cResource, err := ParseURI(concreteURI)
	if err != nil {
		return nil, fmt.Errorf("invalid concrete URI: %w", err)
	}

	if tResource.Scheme != cResource.Scheme {
		return nil, fmt.Errorf("scheme mismatch: expected '%s', got '%s'", tResource.Scheme, cResource.Scheme)
	}

	paramRegex := regexp.MustCompile(`\{([^{}]+)\}`)
	paramMatches := paramRegex.FindAllStringSubmatch(tResource.Path, -1)

	pattern := tResource.Path
	for _, match := range paramMatches {
		placeholder := match[0]
		pattern = strings.Replace(pattern, placeholder, "([^/]+)", 1)
	}

	patternRegex := "^" + pattern + "$"
	matcher := regexp.MustCompile(patternRegex)
	matches := matcher.FindStringSubmatch(cResource.Path)

	if len(matches) == 0 {
		return nil, fmt.Errorf("URI path format mismatch")
	}

	for i, match := range paramMatches {
		paramName := match[1]
		paramValue := matches[i+1]
		cResource.PathParams[paramName] = paramValue
	}

	return cResource, nil
}

func buildURI(template string, pathParams map[string]string, queryParams map[string][]string) (string, error) {
	result := template

	paramRegex := regexp.MustCompile(`\{([^{}]+)\}`)
	paramMatches := paramRegex.FindAllStringSubmatch(template, -1)

	for _, match := range paramMatches {
		placeholder := match[0]
		paramName := match[1]

		value, exists := pathParams[paramName]
		if !exists {
			return "", fmt.Errorf("missing required path parameter '%s'", paramName)
		}

		result = strings.Replace(result, placeholder, value, -1)
	}

	if len(queryParams) > 0 {
		hasQuery := false
		var queryParts []string

		for key, values := range queryParams {
			for _, value := range values {
				queryParts = append(queryParts, fmt.Sprintf("%s=%s", key, value))
				hasQuery = true
			}
		}

		if hasQuery {
			result += "?" + strings.Join(queryParts, "&")
		}
	}

	return result, nil
}
