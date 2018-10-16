package runtime

type Callable interface {
	Params() *ParamSpec
	Call(t *Thread, as *ArgSpec) Value // raw args on stack
}
