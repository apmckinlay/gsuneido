<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/TreeViewControl/Methods">Methods</a></span></div>

### TreeViewControl

``` suneido
(readonly = false, style = 0)
```

A wrapper for the Windows Common Controls "treeview".
`readonly`
: If true `TVS.DISABLEDRAGDROP` is added, otherwise `TVS.EDITLABELS`

`style`
: The standard windows styles are: `WS.CHILD, WS.VISIBLE, WS.BORDER, TVS.HASLINES, TVS.LINESATROOT, TVS.HASBUTTONS, TVS.SHOWSELALWAYS`

Methods:
`AddItem(parent, name, image, container?, param)`
: To add an item to the tree

`GetName(item)`
: Get the item name

`GetParam(item)`
: Get the item number

`GetChildren(item)`
: Get the list of children in an item

`GetParent(item)`
: Get the item parent

`GetImage(item)`
: Get the image code of item

`SetName(item, name)`
: Set the item name

`SetImageList(images)`
: Sets the imagelist for the treeview

`SetImage(item, image)`
: Set the item image code

`SelectItem(item)`
: Select the item name

`EditLabel(item)`
: Permit editing a name

`ExpandItem(item, collapse= false)`
: Expand an item

`Expanded?(item)`
: Return if the item is expanded

`DeleteItem(item)`
: Delete an item

`Children?(item)`
: Return if the item is a children

`HasChildren?(item)`
: Return if the item have childrens

Messages:
`SelectTreeItem(olditem, newitem)`
: 

`Expanding(item)`
: 

`Collapsed(item)`
: 

`Rename(item, newname)`
: 

`Move(dragging, target)`
: 

`Container?(item)`
: 

`Children?(item)`
: 

Example to build a TreeView directly (without a TreeModel):

``` suneido
Controller
    {
    Title: "TreeView example"
    New()
        {
        .tree = .TreeView
        .folder = 0; .openfolder = 1; .document = 2
        images = CreateImageList(16, 16, IDI.FOLDER, IDI.OPENFOLDER,
                                     IDI.DOCUMENT)
        .tree.SetImageList(images)
        rootnode1 = .tree.AddItem(TVI.ROOT, "Names_1", .folder, true)
        rootnode2 = .tree.AddItem(TVI.ROOT, "Names_2", .folder, true)

        .tree.AddItem(rootnode1, 'one', .document)
        rootnode1a=.tree.AddItem(rootnode1, 'two', .folder, true)
        .tree.AddItem(rootnode1, 'three', .document)
        .tree.AddItem(rootnode1, 'four', .document)
        .tree.AddItem(rootnode1a, 'two 1', .document)
        .tree.AddItem(rootnode1a, 'two 2', .document)
        //
        rootnode2a=.tree.AddItem(rootnode2, 'five', .folder, true)
        rootnode2b=.tree.AddItem(rootnode2a, 'six', .folder, true)
        .tree.AddItem(rootnode2b, 'seven', .document, false)
        }
    Controls: (TreeView readonly:)
    Children?(item)
        { return .tree.HasChildren?(item) }
    Expanding(item)
        { .tree.SetImage(item, .openfolder) }
    Collapsed(item)
        { .tree.SetImage(item, .folder) }
    SelectTreeItem(olditem, newitem)
        {
        if (not .tree.HasChildren?(newitem))
            Print(.tree.GetName(newitem))
        }
    }
```

See also: [TreeModel](<TreeModel.md>)