<div style="float:right"><span class="builtin">Windows Builtin</span></div>

#### COMobject

``` suneido
(progid) => comobject or false
(number) => comobject or false
```

Create a Suneido COMobject wrapping a COM IUnknown or IDispatch interface pointer.

The first form attempts to create either an IDispatch interface instance from the given `progid`. It returns `false` if it fails.

The second form can be used when you have a IUnknown pointer value in a Suneido number - like that returned by AtlAxGetControl.

If COMobject doesn't return `false`, it returns an instance of comobject wrapping the underlying IUnknown or IDispatch pointer. The [Dispatch?](<COMobject.Dispatch?.md>) method indicates whether the object created is an IDispatch or merely an IUnknown.

If the comobject wraps an IDispatch, then properties can be accessed as members, and methods may be called as methods. Otherwise, if it merely wraps IUnknown, then only the [Dispatch?](<COMobject.Dispatch?.md>) and [Release()](<COMobject.Release.md>) built-in methods are available.

For example:

``` suneido
ie = COMobject("InternetExplorer.Application")
ie.Visible = true
ie.Navigate("suneido.com")
ie.Release()
```

Handles the following variant types: boolean, integer (I2, I4, I8, U2, U4, U8), unknown, dispatch, null and empty (converted to zero), and string (BSTR).