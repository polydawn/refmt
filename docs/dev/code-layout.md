code layout
===========

- **TokenSourceDriver** *struct*

  Top-level control for a system that emits tokens.

  Both decoding (deserializing binary formats)
  and marshalling (visiting an in-memory structure and preparing it for serialization)
  are token sources.

  View this like a program stack; `TokenSourceMachine`s are the functions in the stack.

  - **TokenSourceMachine** *interface*

	A single state machine that emits tokens implements this.
	This interface makes the individual machines composable to `TokenSourceDriver`.

	In automata theory, `TokenSourceMachine` is generally a DFA FSM -- that is,
	it can operate with a finite amount of memory -- and we use the `TokenSourceDriver`
	to gather together `TokenSourceMachine`s into stacks, thus giving us a
	*pushdown automata*, which has sufficient power to handle recursive structure.

	In practical terms, restricting `TokenSourceMachine`s to a DFA means they
	can operate without allocating memory on the heap.  This confers a large
	bonus to overall performance.

	If a `TokenSourceDriver` is like a stack, and the `TokenSourceMachine` like a function,
	then each time you call the `TokenSourceMachine` is like stepping through a function one
	line (or one instruction) at a time.

    - *Implementations:*
      - **WildcardEncodeMachine** -- turns any `interface{}` into tokens (works by looking up a more specific encode machine, then yielding to it).
      - **WildcardMapEncodeMachine** -- turns a `map[K]V` into tokens (works for any type, using reflection).
      - **StringkeyedMapEncodeMachine** -- turns a `map[string]interface{}` into tokens (faster than the reflection used in the wildcard system).
      - **LiteralEncodeMachine** -- turns primitives like `int` and `string` into tokens (hardly even a DFA; only ever takes one step).
      - **AtlasStructEncodeMachine** -- uses an `Atlas` to visit and emit tokens covering an arbitrary struct type.
      - **PolymorphicUnionEncodeMachine** -- uses a `PolymorphAtlas` to look at an arbitrary value, emit a single-entry map, and trigger a more specific encode machine (the single key is presumably consumed by a `PolymorphicUnionDecoderMachine` and used look up the matching decoder machine).
      - **PolymorphicEnvelopEncodeMachine** -- like `PolymorphicUnionEncodeMachine`, but emits tokens for a map styled like `{kind:typeAbc, msg:{...}}`.

- **TokenSinkDriver** *struct*

  Top-level control for a system that consumes tokens.  The inverse of `TokenSourceDriver`.

  Both encoding (serializing data into a binary format)
  and unmarshalling (visiting an in-memory structure and populating its fields)
  are token sinks.

  View this like a program stack; `TokenSinkMachine`s are the functions in the stack.

  - **TokenSinkMachine** *interface*

    A single state machines that consumes tokens implements this.
	This interface makes the individual machines composable to `TokenSinkDriver`.

	See the documentation of `TokenSourceMachine` for more information;
	sinks and sources follow the same model (e.g. these are DFAs, etc).

    - *Implementations:*
      - **WildcardDecodeMachine** -- populates an `interface{}` (usually with `map[string]interface{}` or `[]interface{}`, for lack of more specific type info).
      - **WildcardMapDecodeMachine** -- populates a `map[interface{}]interface{}`.
      - **StringkeyedMapDecodeMachine** -- populates a `map[string]interface{}` (this will
        yield errors if non-string types are found as map keys).
      - **LiteralDecodeMachine** -- populates `string`, `int`, etc.
      - **AtlasStructDecodeMachine** -- uses an `Atlas` to visit fields (presumably all in one structure, but the sky's the limit really since `Atlas` can suggest arbitrary memory locations).
      - **PolymorphicUnionDecodeMachine** -- expects to consume a single-entry map, and will attempt to shell out to a more specific decoder machine based on the key string.
      - **PolymorphicEnvelopDecodeMachine** -- like `PolymorphicUnionDecodeMachine`, but expects to consume two entries: the value of one will be used as a type hint, and the other as the content to send to the specific decoder machine (note this may be significantly less efficient to decode, since contrasted with the union machine, it may require buffering if the type hint entry doesn't come first).
