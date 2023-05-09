package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/cantara/bragi/sbragi"
	"github.com/cantara/buri/version/filter"
	"github.com/cantara/buri/version/release"
	"github.com/cantara/buri/version/snapshot"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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

	r, err := git.PlainOpen(".")
	if err != nil {
		log.WithError(err).Fatal("could not find git repo on disk")
		return
	}
	tags, err := r.Tags()
	if err != nil {
		log.WithError(err).Fatal("while getting tags")
	}
	defer tags.Close()
	var newest *release.Version
	var hash plumbing.Hash
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		if ref == nil {
			log.Debug("tag refference was nil")
			return nil
		}
		if !ref.Name().IsTag() {
			log.Debug("refference was not tag")
			return nil
		}
		v, err := release.Parse(ref.Name().Short())
		if err != nil {
			log.WithError(err).Warning("tag was not release version", "tag", ref.Name().Short())
			return nil
		}
		if newest == nil {
			newest = &v
			hash = ref.Hash()
			return nil
		}
		if newest.IsStrictlySemanticNewer(filter.AllReleases, v) {
			newest = &v
			hash = ref.Hash()
			return nil
		}
		return nil
	})
	if err != nil {
		log.WithError(err).Fatal("during tag iteration")
		return
	}
	if newest == nil {
		fmt.Print("v0.0.0")
		return
	}
	ref, err := r.Head()
	if err != nil {
		log.WithError(err).Fatal("while getting git head")
		return
	}
	hco, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.WithError(err).Fatal("while getting head commit object")
		return
	}
	tco, err := r.CommitObject(hash)
	if err != nil {
		log.WithError(err).Fatal("while getting head commit object")
		return
	}
	sv := snapshot.Version{
		Version:   *newest,
		TimeStamp: time.Now(),
		Iteration: len(hco.ParentHashes) - len(tco.ParentHashes) + 1,
	}
	fmt.Println(sv.String())
}
