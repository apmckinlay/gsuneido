package runtime

// Context is the current Thread
// can't use actual Thread because that causes circular dependency
type Context interface {
	// Call interprets a compiled Suneido function
	Call(fn *SuFunc) Value
}

type Callable interface {
	Params() *ParamSpec
	Call(ctx Context, args ...Value) Value
}

// type CallCooked struct {
// 	fn func (ctx Context, args ...Value) Value
// }

// func (cc CallCooked) Call(ctx Context, as *ArgSpec, args ...Value) Value {

// }
