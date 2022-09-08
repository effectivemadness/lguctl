package execute

import "strings"

type Environ []string

// Unset removes environment variable
func (e *Environ) Unset(key string) {
	for i := range *e {
		if strings.HasPrefix((*e)[i], key+"=") {
			(*e)[i] = (*e)[len(*e)-1]
			*e = (*e)[:len(*e)-1]
			break
		}
	}
}

// Set adds environment variable
func (e *Environ) Set(key, val string) {
	e.Unset(key)
	*e = append(*e, key+"="+val)
}
