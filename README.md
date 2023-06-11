**Trelay** - a set of utilities to simplify work with Terraria Networking.

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

switch p.ID {
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

Or write to `trelay.Packet` to build a packet in memory when length is
not known in advance.
```go
p := &trelay.Packet{ID: 1}
trelay.Fprint(p, "Terraria123") // writes to Packet never* fail

_, err := trelay.Fprint(w, p) // no need to p.Bytes(), Fprint will use p.WriteTo
if err != nil { /* ... */ }
```
> *Assumes Fprint-supported types and no memory issues.

Reuse `trelay.Packet` whenever possible to avoid
unnecessary memory allocations.
```go
_, err := trelay.Fscan(conn, &packet) // Fscan into Packet clears it
if err != nil { /* ... */ }
```
