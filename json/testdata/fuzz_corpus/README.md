Fuzz seed corpus for the json package, loaded by `seedJSONCorpus` in
`fuzz_test.go`. Each subdirectory contains seed files used to prime
`FuzzJSONDecode` and related targets.

- `codec_fixtures/` — DAG-JSON positive cases from
  https://github.com/ipld/codec-fixtures (Apache-2.0). One file per
  fixture, named after the fixture's CID.
- `negative/` — Inputs that strict decoders should reject; hex
  literals extracted from
  https://github.com/ipld/codec-fixtures/negative-fixtures/dag-json.
