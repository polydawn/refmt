#!/bin/bash
set -euo pipefail

### Normalize path -- all work should be relative to this script's location.
## Set up gopath -- also relative to this dir, so we work in isolation.
cd "$( dirname "${BASH_SOURCE[0]}" )"
export GOPATH="$PWD/.gopath/"

funcs=()
funcs+=("Benchmark_ArrayFlatIntToJson_Xlate")
funcs+=("Benchmark_ArrayFlatIntToCbor_Xlate")
funcs+=("Benchmark_ArrayFlatIntToJson_Stdlib")
funcs+=("Benchmark_ArrayFlatStrToJson_Xlate")
funcs+=("Benchmark_ArrayFlatStrToCbor_Xlate")
funcs+=("Benchmark_ArrayFlatStrToJson_Stdlib")

funcs+=("Benchmark_StructToJson_XlateFieldRoute")
funcs+=("Benchmark_StructToCbor_XlateFieldRoute")
funcs+=("Benchmark_StructToJson_XlateAddrFunc")
funcs+=("Benchmark_StructToCbor_XlateAddrFunc")
funcs+=("Benchmark_StructToJson_Stdlib")

profPath=".gopath/tmp/prof/" ; mkdir -p "$profPath"
go test -i .
echo "${funcs[@]}" | tr " " "\n" | xargs -n1 -I{} \
	go test \
		-run=XXX -bench={} \
		-o "$profPath/bench.bin" \
		-cpuprofile="$profPath/{}.cpu.pprof" \
		2> /dev/null | grep "^Benchmark_"
echo "${funcs[@]}" | tr " " "\n" | xargs -n1 -I{} \
	go tool pprof \
		--pdf \
		--output "$profPath/{}.cpu.pdf" \
		"$profPath/bench.bin" "$profPath/{}.cpu.pprof"
#ls -lah "$profPath"/*.pdf # 'go tool' already says where it puts it
