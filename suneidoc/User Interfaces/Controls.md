## Controls

User interfaces are built from controls.  These controls include standard Windows controls:
<pre>
<a href="/suneidoc/User Interfaces/Reference/ButtonControl">ButtonControl</a>
<a href="/suneidoc/User Interfaces/Reference/CheckBoxControl">CheckBoxControl</a>
<a href="/suneidoc/User Interfaces/Reference/EditorControl">EditorControl</a> - an "edit" control configured for multi-line use
<a href="/suneidoc/User Interfaces/Reference/FieldControl">FieldControl</a> - an "edit" control configured for single-line use
<a href="/suneidoc/User Interfaces/Reference/ListBoxControl">ListBoxControl</a>
<a href="/suneidoc/User Interfaces/Reference/StaticControl">StaticControl</a>
</pre>

Windows common controls:
<pre>
<a href="/suneidoc/User Interfaces/Reference/ListViewControl">ListViewControl</a>
<a href="/suneidoc/User Interfaces/Reference/StatusbarControl">StatusbarControl</a>
<a href="/suneidoc/User Interfaces/Reference/TabControl">TabControl</a>
<a href="/suneidoc/User Interfaces/Reference/ToolbarControl">ToolbarControl</a>
<a href="/suneidoc/User Interfaces/Reference/TreeViewControl">TreeViewControl</a>
</pre>

controls for composing layouts:
<pre>
<a href="/suneidoc/User Interfaces/Reference/HorzControl">HorzControl</a>
<a href="/suneidoc/User Interfaces/Reference/VertControl">VertControl</a>
<a href="/suneidoc/User Interfaces/Reference/CenterControl">CenterControl</a>
<a href="/suneidoc/User Interfaces/Reference/FillControl">FillControl</a>
<a href="/suneidoc/User Interfaces/Reference/FormControl">FormControl</a>
<a href="/suneidoc/User Interfaces/Reference/ScrollControl">ScrollControl</a>
<a href="/suneidoc/User Interfaces/Reference/PairControl">PairControl</a>
<a href="/suneidoc/User Interfaces/Reference/VertSplitControl">VertSplitControl</a>
<a href="/suneidoc/User Interfaces/Reference/HorzSplitControl">HorzSplitControl</a>
<a href="/suneidoc/User Interfaces/Reference/ExpandControl">ExpandControl</a>
RepeatControl
</pre>

and more complex, composite controls:
<pre>
<a href="/suneidoc/User Interfaces/Reference/VfieldsControl">VfieldsControl</a>
<a href="/suneidoc/User Interfaces/Reference/InspectControl">InspectControl</a>
<a href="/suneidoc/User Interfaces/Reference/ExplorerControl">ExplorerControl</a>
<a href="/suneidoc/User Interfaces/Reference/AccessControl">AccessControl</a>
<a href="/suneidoc/User Interfaces/Reference/BrowseControl">BrowseControl</a>
</pre>

Control specifications are written as objects, for example:

``` suneido
#(Static, "hello world")
```

The first value in a specification object is the name of the control. If there are no arguments, you can simply write the name of the control.  i.e. `#(Fill)` can be written as `'Fill'`.

Note: The "Control" suffix is omitted in control specifications, e.g. #(Button OK) instead of #(ButtonControl OK)