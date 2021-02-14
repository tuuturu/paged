package core

import (
	"fmt"
	"regexp"
)

var dsnRegex = regexp.MustCompile(`(\w+)://(\w+):(.+)@([\w-]+):(\d+)/(\w+)`)

func parseDSN(rawDSN string) *DSN {
	matches := dsnRegex.FindStringSubmatch(rawDSN)

	if len(matches) != 7 {
		return nil
	}

	return &DSN{
		Scheme:       matches[1],
		Username:     matches[2],
		Password:     matches[3],
		URI:          matches[4],
		Port:         matches[5],
		DatabaseName: matches[6],
	}
}

func (d DSN) String() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		d.Scheme,
		d.Username,
		d.Password,
		d.URI,
		d.Port,
		d.DatabaseName,
	)
}
