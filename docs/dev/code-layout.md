code layout
===========

-----------
User-facing
-----------

- example code:
  ```
  xlate.NewJsonEncoder(stdout).Marshal(123)
  ```
  - Creates a `TokenSourceDriver` as the `TokenSource` for walking the object (`123`).
  - Creates a `JsonSerializer` as the `TokenSink` outputting to `stdout`.
  - Both are placed into a `TokenPump`, which powers the rest of the transaction.
  - The `TokenSourceDriver` sees the type of object it was given and immediately delegates to a `LiteralEncodeMachine`.
  - The `LiteralEncodeMachine` steps once, yields a token (it's an `int`), and reports done.  There are no other stacked machines, so the `TokenSourceDriver` overall reports done.
  - The `JsonSerializer`, invoked and given a token by the `TokenPump`, writes the number to `stdout`.  The `JsonSerializer` knows that it just wrote a literal type, and since it's not deep in an object tree that's the end of a json entity, so it reports done.
  - The `TokenPump` sees both sides finished in unison, and returns done, with no error!

-----------------------
Token stream interfaces
-----------------------

Xlate handles every translation by converting things into a token stream,
then processing the token stream into the desired result format.

`TokenSource` and `TokenSink` describe how to produce and process token streams, respectively.
Listing their implementations effectively lists every format that xlate can convert to and from.

- **TokenSource** *interface*

  Anything that emits a stream of tokens.

  Call a `TokenSource` repeatedly; each call will yield a token into the memory provided
  (and typically cause the concrete implementation to advance the `TokenSource`'s internal state by single step, as appropriate).
  Errors or end-of-stream are indicated by the return codes.

  - *Implementations*:
    - **JsonDeserializer** -- constructed with an `io.Reader`, from which (hopefully-)json-formatted bytes will be consumed and converted into tokens.
    - **CborDeserializer** -- constructed with an `io.Reader`, from which (hopefully-)cbor-formatted bytes will be consumed and converted into tokens.
    - **TokenSourceDriver** -- constructed with a reference to any object, which will be visited and all fields emitted one by one as tokens.

- **TokenSink** *interface*

  Anything that consumes a stream of tokens.

  Call a `TokenSink` repeatedly; each call will copy the token content into its new home
  (whatever that may be, whether a serial byte stream or a complex in-memory layout).
  Errors and expected end-of-stream are indicated by the return codes.

    - **JsonSerializer** -- constructed with an `io.Writer`, to which json-formatted bytes are flushed as each token is received.
    - **CborSerializer** -- constructed with an `io.Writer`, to which cbor-formatted bytes are flushed as each token is received.
    - **TokenSinkDriver** -- constructed with a reference to any object (or empty `interface{}`), which will be populated based on tokens received.

- **TokenPump** *struct*

  Joins a `TokenSource` and `TokenSink`, advancing them in lock-step, and looping until they're complete.

  The `TokenPump` will return one of:
    - an error, if any step function of either the source or sink returns an error;
    - an error if the sink expects to be done while the source is continuing;
    - or a done bool with nil error, if both the source and sink become done during the same step.

---------------------------
Object/Token morphism tools
---------------------------

These interfaces compose to make `TokenSource`s and `TokenSink`s that operate on in-memory data.
working smoothly with the golang type system.

Users will rarely interact with any of these directly.

- **TokenSourceDriver** *struct*

  Top-level control for a system that walks in-memory structures and emits tokens
  describing the object it's observing.

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

  Top-level control for a system that walks in-memory structures and consumes tokens,
  populating the in-memory fields.  The inverse of `TokenSourceDriver`.

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
