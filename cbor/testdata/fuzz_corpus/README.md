Fuzz seed corpus for the cbor package, loaded by `seedCborCorpus` in
`fuzz_test.go`. Each subdirectory contains binary seed files used to
prime `FuzzCborDecode` and related targets.

- `codec_fixtures/` — DAG-CBOR positive cases from
  https://github.com/ipld/codec-fixtures (Apache-2.0). One file per
  fixture, named after the fixture's CID.
- `appendix_a/` — CBOR test vectors from RFC 7049 Appendix A, as
  vendored by https://github.com/rvagg/cborg (Apache-2.0).
- `negative/` — Inputs that DAG-CBOR-strict decoders should reject;
  hex literals extracted from
  https://github.com/ipld/codec-fixtures/negative-fixtures/dag-cbor.
