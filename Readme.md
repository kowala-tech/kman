## K-man: naturalsitic documentation parser and presentation generator

The goal is parse help topics and glossary terms from arbitrary markdown (and optionally source code) files, and then generate some useful output from them.

This is a proof of concept that does the parsing. Generation of HTML (or whatever) is still to do.

In order to try it out, run `go run cmd/kman/*.go --go --md`. This will parse all topics from markdown and go files (examples in the `doc` directory and dump them to stdout.
