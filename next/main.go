package main

import (
	"bufio"
	"fmt"
	"os"

	log "github.com/cantara/bragi/sbragi"
	"github.com/cantara/buri/version/release"
	"golang.org/x/exp/slog"
)

func main() {
	n, err := log.NewLogger(slog.HandlerOptions{
		AddSource: true,
		Level:     log.LevelWarning,
	}.NewTextHandler(os.Stderr))
	if err != nil {
		log.WithError(err).Fatal("while creating error logger")
		return
	}
	n.SetDefault()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		v, err := release.Parse(scanner.Text())
		if err != nil {
			log.WithError(err).Warning("tag was not a release version")
			continue
		}

		v.Patch = v.Patch + 1
		fmt.Println(v.String())

		v.Patch = 0
		v.Minor = v.Minor + 1
		fmt.Println(v.String())

		v.Minor = 0
		v.Major = v.Major + 1
		fmt.Println(v.String())
	}
}
