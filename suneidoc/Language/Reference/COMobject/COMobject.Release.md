<div style="float:right"><span class="builtin">Builtin</span></div>

#### comobject.Release

``` suneido
() => 0
```

Releases the resources used by the comobject.

Once this method is called, any subsequent method call or property access on a comobject will trigger an exception.

The Suneido programmer must ensure this method is called on every comobject or resources will be leaked.

When a comobject is instantiated using a pre-existing IUnknown pointer on an existing underlying COM object [the `COMobject(number)` constructor], the comobject assumes ownership of the pre-existing pointer and will release the IUnknown interface at the appropriate time. Once a comobject "owns" the pre-existing pointer, other code should assume it is no longer valid.

When a cmoobject is instantiated using a progid, it causes COM to instantiate a new underlying COM object and return an IDispatch pointer or, if IDispatch is not available, IUnknown. The comobject "owns" this interface pointer and will release it when the Release member is called.

See also: [Dispatch?](<COMobject.Dispatch?.md>)