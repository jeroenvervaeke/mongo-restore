package main

import (
	"os"

	"github.com/mongodb/mongo-tools/common/log"
	"github.com/mongodb/mongo-tools/common/signals"
	"github.com/mongodb/mongo-tools/common/util"
	"github.com/mongodb/mongo-tools/mongorestore"
)

var (
	VersionStr = "0.0.0-POC"
	GitCommit  = "build-without-git-commit"
)

// MongoRestore main, copied from: https://github.com/mongodb/mongo-tools/blob/0fe8aa034b5152a9d03c74568e1add529d47a52c/mongorestore/main/mongorestore.go
func main() {
	// Only modification needed pass os.Args[2:] instead of os.Args[1:], since this is a subcommand now
	opts, err := mongorestore.ParseOptions(os.Args[1:], VersionStr, GitCommit)

	if err != nil {
		log.Logvf(log.Always, "error parsing command line options: %s", err.Error())
		log.Logvf(log.Always, util.ShortUsage("mongorestore"))
		os.Exit(util.ExitFailure)
	}

	// print help or version info, if specified
	if opts.PrintHelp(false) {
		return
	}

	if opts.PrintVersion() {
		return
	}

	restore, err := mongorestore.New(opts)
	if err != nil {
		log.Logvf(log.Always, err.Error())
		os.Exit(util.ExitFailure)
	}
	defer restore.Close()

	finishedChan := signals.HandleWithInterrupt(restore.HandleInterrupt)
	defer close(finishedChan)

	result := restore.Restore()
	if result.Err != nil {
		log.Logvf(log.Always, "Failed: %v", result.Err)
	}

	if restore.ToolOptions.WriteConcern.Acknowledged() {
		log.Logvf(
			log.Always,
			"%v document(s) restored successfully. %v document(s) failed to restore.",
			result.Successes,
			result.Failures,
		)
	} else {
		log.Logvf(log.Always, "done")
	}

	if result.Err != nil {
		os.Exit(util.ExitFailure)
	}
	os.Exit(util.ExitSuccess)
}
