#### CanvasItem

``` suneido
( )
```

Abstract base class for other canvas items.

Canvas items normally have a **New** method that accepts the information required to draw the item.  All canvas items must implement a **Paint** method that is passed hdc - a handle to the device context.  If necessary, a **Destroy** method can be defined to free any resources.  CanvasItem defines an empty default Destroy method.

Note: It is not strictly necessary for canvas items to derive from CanvasItem, as long as they have Paint and Destroy methods.