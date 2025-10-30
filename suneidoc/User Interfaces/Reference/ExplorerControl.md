<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/ExplorerControl/Methods">Methods</a></span></div>

### ExplorerControl

``` suneido
( treemodel, view )
```

Creates a HorzSplitControl with a TreeViewControl in the left hand pane and the view control in the right hand pane.

The treemodel must have the following methods:
`Children( num ) => list`
: The object must be a list of objects containing name, num, and group.

`Children?( num ) => true or false`
: 

`Container?( num ) => true or false`
: 

`DeleteItem( num )`
: 

`Get( num ) => object`
: 

`NewItem( object )`
: 

`Nextnum() => number`
: 

`Static?( num ) => true or false`
: 

`Update( object )`
: 

The view must have the following methods:
`Get() => object`
: 

`Set(object)`
: 

`GetState()`
: For persistence.

`SetState(stateobject)`
: For persistence.

`Dirty?() => true or false`
: 

`SetReadOnly(readonly = true)`
: 

See also:
[ExplorerAdapter](<ExplorerAdapterControl.md>)