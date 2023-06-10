**Trelay** -- a set of utilities to simplify work with Terraria Networking.

Use `trelay.Fscan` to read data from a reader (network) as sent by Terraria.
```go
var ln uint16
var id byte
_, err := trelay.Fscan(r, &ln, &id)
if err != nil { /* ... */ }
switch id {
case 1:
  var ver string
  _, err := trelay.Fscan(r, &ver)
/* ... */
```

Or use `trelay.Packet` to read an entire packet into a buffer.
```go
var p trelay.Packet
_, err := trelay.Fscan(r, &p) // p.ReadFrom(r) also works!
if err != nil { /* ... */ }

switch p.Id() {
case 1:
  var ver string
  _, err := trelay.Fscan(p, &ver)
/* ... */
```

Use `trelay.Fprint` to write data as expected by Terraria.
```go
// do not forget to cast constants to correct types
_, err := trelay.Fprint(w, uint16(15), byte(1), "Terraria123")
```

Or offload length tracking to `trelay.Builder`.
```go
b := &trelay.Builder{ID: 1}
trelay.Fprint(b, "Terraria123") // writes to Builder never* fail

_, err := trelay.Fprint(w, b) // sending b.Bytes() also works.
if err != nil { /* ... */ }
```
> *Assumes supported types and no memory issues.

Reuse `trelay.Packet` and `trelay.Builder` whenever possible to avoid
unnecessary memory allocations.
```go
_, err := trelay.Fscan(conn, &packet) // Fscan into Packet clears it
if err != nil { /* ... */ }

/* ... */

defer builder.Reset() // reset Builders manually after use
_, err := trelay.Fprint(builder, "Terraria123")
```