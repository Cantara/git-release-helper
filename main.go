package main

import (
	"bufio"
	"fmt"
	"os"

	log "github.com/cantara/bragi/sbragi"
	"github.com/cantara/buri/version/filter"
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

	var newest *release.Version
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		v, err := release.Parse(scanner.Text())
		if err != nil {
			log.WithError(err).Warning("tag was not a release version")
			continue
		}
		if newest == nil {
			newest = &v
			continue
		}
		if newest.IsStrictlySemanticNewer(filter.AllReleases, v) {
			newest = &v
			continue
		}
	}

	if newest == nil {
		fmt.Println("v0.0.0")
		return
	}
	fmt.Println(newest.String())
}
