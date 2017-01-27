#!/bin/bash
set -euo pipefail

### Normalize path -- all work should be relative to this script's location.
## Set up gopath -- also relative to this dir, so we work in isolation.
cd "$( dirname "${BASH_SOURCE[0]}" )"
export GOPATH="$PWD/.gopath/"

funcs=()
funcs+=("Benchmark_ArrayFlatIntToJson_Xlate")
funcs+=("Benchmark_ArrayFlatIntToJson_Stdlib")
funcs+=("Benchmark_ArrayFlatStrToJson_Xlate")
funcs+=("Benchmark_ArrayFlatStrToJson_Stdlib")

profPath=".gopath/tmp/prof/" ; mkdir -p "$profPath"
go test -i .
echo "${funcs[@]}" | tr " " "\n" | xargs -n1 -I{} \
	go test \
		-run=XXX -bench={} \
		-o "$profPath/bench.bin" \
		-cpuprofile="$profPath/{}.cpu.pprof"
echo "${funcs[@]}" | tr " " "\n" | xargs -n1 -I{} \
	go tool pprof \
		--pdf \
		--output "$profPath/{}.cpu.pdf" \
		"$profPath/bench.bin" "$profPath/{}.cpu.pprof"
ls -lah "$profPath"/*.pdf
