# Pointer-Only Rows

## Goal

Replace `Row`'s `DbRec` elements with record data only, then use
`util/multi.Ptrs` to optimize the common single-record row. Recover local
database offsets from record pointers inside query execution. Resolve write
identity explicitly at the dbms boundary and retain it in `SuRecord.recoff`.

This separates three concerns: query row representation, local storage address
resolution, and the identity used for subsequent updates and deletes.

## Expected Benefits

- Moving from `DbRec` to `Record` reduces row element size from 24 to 16 bytes
  on 64-bit systems, without changing record contents.
- Moving to `multi.Ptrs[byte]` reduces the row value from a 24-byte slice header
  to a 16-byte value and eliminates the backing-array allocation for one record.
- Multi-record rows use an array of 8-byte pointers rather than string
  descriptors or `DbRec` values.
- Clients are a particularly good fit because server responses normally contain
  one flattened record. The received record allocation remains, but the separate
  one-element row allocation disappears.
- Sorting and buffering preserve record pointers, so offsets do not need to be
  carried separately through each query operator.

## Architectural Boundary

- Local query execution can use its transaction's `Stor` to recover offsets.
  This includes embedded/local execution, not just a network server process.
- Query rows retain the original storage-backed record pointer where database
  identity is needed. Rebuilding an updateable source record loses that identity.
- At the dbms boundary, resolve the updateable source record's offset before
  flattening, squeezing deleted fields, or otherwise copying record bytes.
- Pass that offset explicitly into `SuRecord` construction for local results.
  For remote results, continue transmitting the offset separately and pass the
  received value into `SuRecord` construction.
- Keep `SuRecord.recoff`: an update replaces this offset independently of the
  original row. Do not later recover write identity from that original row.
- Non-updateable results and synthetic record construction need no write offset.
- Clients and `SuRecord` never call the storage reverse lookup. No process-wide
  client/server switch is needed.
- Avoid introducing a `core` dependency on `stor`; `stor` already imports `core`.

## Storage Reverse Lookup

Add a per-`Stor` address index and an operation that resolves a pointer to an
offset, with an explicit indication of whether the pointer belongs to the store.
Offset zero is valid at the storage layer and cannot serve as its not-found
sentinel, even though database data records do not start there.

- Use chunk-size address buckets, keyed by the pointer's high bits.
- Chunk size is a power of two, but chunk addresses are not guaranteed to be
  aligned to that size. Register every bucket intersected by a chunk's range.
- Each bucket holds candidate chunk descriptors, not just one offset. Different
  chunks can share a bucket; small heap chunks must also work.
- Retain each chunk's base address, accessible length, and storage offset.
  Validate containment before calculating the pointer's offset within the chunk.
- Use the actual accessible slice length for partial read-only chunks. A valid
  address range is not necessarily limited to the currently allocated data size;
  define the operation as address translation, not allocation validation.
- Register initial chunks in `NewStor` and new chunks in `extend`. Publish the
  index entry before publishing the chunk to readers or allowing allocations in it.
- Keep index reads safe during extension. A concurrent map with immutable bucket
  candidate lists is one possible implementation; never mutate a list visible to
  readers. Existing extension locking can serialize writers.
- Do not add locks or lookups to the normal `Data` or `Alloc` paths.
- Define behavior for nil, foreign, one-past-end, and closed/unmapped pointers.
  The index must not imply that an unmapped pointer remains usable.
- Do not convert retained numeric addresses back into pointers. Keep the existing
  chunk slices as owners of the mapped or heap storage.

At 40 GiB and 64 MiB chunks, there are about 640 chunks and at most 1,280 bucket
registrations. Index maintenance should be small; the cost of replacing field
loads with lookups still needs workload measurement.

## Removing DbRec.Off

Audit consumers by purpose rather than mechanically replacing every field read.

- Update/delete actions need the original database identity and the new identity
  returned by update. Preserve their protection against processing a newly
  encountered updated version again.
- Bloom filtering should continue using offsets already available from iterators
  wherever possible. Resolve an offset only where the candidate arrives as a row.
- Temporary indexes retain records or rows, so their pointers can survive sorting
  without separate offset metadata. Verify that updateable results retain the
  appropriate original source component.
- Server responses must obtain their write offset before `rowToRecord` changes
  the representation. Keep the existing wire offset field.
- Local dbms APIs must supply explicit write identity when handing results to
  core. Ensure it belongs to the returned result, not to a later lookahead row.
- Update direct write consumers outside the main query API, including LLM tools.
- Remove obsolete constructors, mocks, tests, and assumptions about `DbRec`.

### SameAs

Use pointer equality as a fast path, followed by component-wise full physical
record equality. This avoids offset lookup for comparisons and works for derived
records as well as stored records.

- Preserve component positions and missing components; do not flatten rows for
  this comparison.
- Do not compare only projected fields or grouping columns. Project and semijoin
  use this operation to recognize a particular representative within a group.
- The existing callers rely on stable traversals and complete underlying records,
  including hidden key fields. Key uniqueness makes complete contents sufficient
  to distinguish live table records in those contexts.
- This is not general version identity: identical contents at different historical
  offsets can compare equal. Audit the callers and document the new contract.
- Consider a name or comment change so future callers do not assume offset identity.

### Derived Records and Accounting

`Off == 0` currently also marks derived records. Storage membership is not an exact
replacement: trigger rows deliberately omit identity even when their bytes came
from storage.

- Audit the `Derived` consumers in project, semijoin, and temporary indexes. They
  currently support memory warnings, not result identity or write permissions.
- Choose whether to retain approximate derived-byte accounting using storage
  membership, account for constructed records elsewhere, or revise the warning.
- Do not introduce per-component metadata solely to preserve this warning without
  weighing it against the representation savings.
- Keep updateability controlled by query/API context, not by whether bytes happen
  to belong to a store. Trigger and synthetic records must not acquire write
  identity merely because of their address.

## Pointer-Only Record Access

Use `multi.Ptrs[byte]` to point at serialized record bytes, not at allocated
`Record` descriptors. Provide a small row-access abstraction so field access can
continue using temporary, zero-copy `Record` views where convenient.

- Recover record length from the serialized header. Add pointer-oriented length
  decoding rather than fabricating an arbitrarily large byte slice for `RecLen`.
- Handle the one-byte empty record before reading any additional header bytes;
  then read only the length bytes required by the record format.
- Use nil pointers for missing components. A missing component is different from
  the valid one-byte empty record.
- Preserve the distinction between no result and a row containing missing
  components. Audit existing nil checks, indexing, ranges, slicing, and appends.
- Keep heap record pointers GC-visible. `multi.Ptrs` uses `unsafe.Pointer` and a
  typed pointer backing array; do not replace either with retained `uintptr`s.
- Preserve immutability and ownership when pointing into received or constructed
  record buffers. Buffer reuse must not invalidate retained rows.
- Avoid repeatedly decoding length when processing multiple fields from the same
  component. Reconstruct its `Record` view once where practical.
- Review the existing `util/multi` implementation and tests as a prerequisite,
  including nil elements and retention of heap-backed records.

## Implementation Sequence

1. Record baseline allocation and timing measurements for local scans, client
   result scans, joins, projection, and temporary-index sorting.
2. Add and test the independent `Stor` reverse lookup, including concurrent
   extension. Leave existing row representation unchanged initially.
3. Establish the explicit dbms-to-`SuRecord` write-offset boundary for local and
   remote results, retaining the existing wire format and `recoff` behavior.
4. Replace `DbRec` with `Record` in rows. Update query offset consumers to use
   direct iterator offsets or storage lookup, revise `SameAs`, and settle the
   derived-memory accounting policy.
5. Replace the row slice representation with `multi.Ptrs[byte]`. Introduce the
   pointer-based length/access helpers and update row operations and constructors.
6. Run correctness checks and repeat the baseline measurements. Optimize lookup
   frequency or record-view reconstruction only where measurements justify it.

These stages isolate storage translation, identity plumbing, and representation
changes so failures and performance differences are easier to attribute.

## Validation

- Reverse-lookup round trips for initial and newly allocated chunks, interior
  pointers, unaligned chunk bases, shared buckets, partial chunks, and heap stores.
- Rejection of nil and foreign pointers, plus checks at chunk boundaries and
  races between lookup and extension.
- Local and remote record update/delete, including repeated updates through one
  `SuRecord`, sorted updateable queries, and action iteration after an update.
- Forward/reverse traversal and rewind for project and semijoin; multi-index
  single-record lookup must still detect genuinely nonunique results.
- Derived records, triggers, unions, left joins, missing components, empty records,
  and multi-record joins.
- Length decoding across all header formats and GC retention of heap-built and
  client-received records, including after joining rows.
- Protocol compatibility and correct offset capture before server-side record
  flattening or squeezing.
- Compare allocation counts, retained memory, and throughput. Confirm the client
  single-record row allocation disappears without regressing query hot paths.
- For implementation changes, format Go files and run
  `go test -short -timeout=30s ./...`, relevant race/checkptr tests with the same
  short/timeout flags, `go vet`, and workspace diagnostics. Use `make` if building.

Do not use or overwrite the working `suneido.db` for validation; use isolated
test databases.
