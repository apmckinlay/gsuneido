#### sequence.Instantiate

``` suneido
() => this
```

Forces the sequence to be "instantiated" i.e. converted to a regular object.

Very large or infinite sequences cannot be instantiated.

**Note**: Normally this is not required. Sequences will be automatically instantiated when necessary.

For convenience, Instantiate is also defined for objects, in which case it does nothing.

See also: [sequence.Instantiated?](<sequence.Instantiated?.md>), [Sequence](<../../Basic Data Types/Sequence.md>)