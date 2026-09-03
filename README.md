# go-molding

`go-molding` generates allocation-light reference collectors and appliers from
Go struct tags. Generated code performs no runtime reflection and supports
arbitrary reference kinds without generator changes.

## Usage

Declare each reference source and destination on the same struct:

```go
type Item struct {
	UserID     string `json:"userId" ref:"user:UserName"`
	UserName   string `json:"userName"`
	TenantId   string `json:"tenantId" ref:"tenant:TenantName"`
	TenantName string `json:"tenantName"`
}

//go:generate go run github.com/ninenhan/go-molding -type Page -output page_molding.gen.go
type Page struct {
	List []Item `json:"list" ref:"inline"`
}
```

Then run:

```sh
go generate ./...
```

The generated type exposes:

```go
RefIDs(kind string) []string
ApplyRefs(kind string, references map[string]string)
```

For example, `RefIDs("user")` and `RefIDs("tenant")` return independent,
deduplicated ID sets. `ApplyRefs` writes values only for IDs present in the
supplied map.
