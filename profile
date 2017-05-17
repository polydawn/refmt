#!/bin/bash
set -euo pipefail

### Normalize path -- all work should be relative to this script's location.
## Set up gopath -- also relative to this dir, so we work in isolation.
cd "$( dirname "${BASH_SOURCE[0]}" )"
export GOPATH="$PWD/.gopath/"

funcs=()
funcs+=("Benchmark_ArrayFlatIntToJson_Refmt")
funcs+=("Benchmark_ArrayFlatIntToJson_RefmtLegacy")
funcs+=("Benchmark_ArrayFlatIntToJson_Stdlib")
funcs+=("Benchmark_ArrayFlatIntToCbor_Refmt")
funcs+=("Benchmark_ArrayFlatIntToCbor_RefmtLegacy")
funcs+=("Benchmark_ArrayFlatStrToJson_Refmt")
funcs+=("Benchmark_ArrayFlatStrToJson_RefmtLegacy")
funcs+=("Benchmark_ArrayFlatStrToJson_Stdlib")
funcs+=("Benchmark_ArrayFlatStrToCbor_Refmt")
funcs+=("Benchmark_ArrayFlatStrToCbor_RefmtLegacy")

funcs+=("Benchmark_StructToJson_Refmt")
funcs+=("Benchmark_StructToJson_RefmtLegacyFieldRoute")
funcs+=("Benchmark_StructToJson_RefmtLegacyAddrFunc")
funcs+=("Benchmark_StructToJson_Stdlib")
funcs+=("Benchmark_StructToCbor_Refmt")
funcs+=("Benchmark_StructToCbor_RefmtLegacyFieldRoute")
funcs+=("Benchmark_StructToCbor_RefmtLegacyAddrFunc")

profPath=".gopath/tmp/prof/" ; mkdir -p "$profPath"
go test -i .
echo "${funcs[@]}" | tr " " "\n" | xargs -n1 -I{} \
	go test \
		-run=XXX -bench=^{}\$ \
		-o "$profPath/bench.bin" \
		-cpuprofile="$profPath/{}.cpu.pprof" \
		2> /dev/null | grep "^Benchmark_"
echo "${funcs[@]}" | tr " " "\n" | xargs -n1 -I{} \
	go tool pprof \
		--pdf \
		--output "$profPath/{}.cpu.pdf" \
		"$profPath/bench.bin" "$profPath/{}.cpu.pprof"
#ls -lah "$profPath"/*.pdf # 'go tool' already says where it puts it
