### Boxes and Stretch

This is a simple method of arranging rectangular items relative to each other in a way that makes it easy to add, resize, or remove items.  Although it might appear more awkward to specify a layout this way, it is much easier to modify than if, for example, you had used a visual tool that simply recorded the exact locations of each item.

Items have minimum x and y dimensions as well as x and y stretch factors.  If a stretch factor is 0 then the item has a fixed size in that dimension.  If a stretch factor is greater than 0, then the item may be larger depending on the size of the container it is placed in.

The two most common containers in this pattern are ones that place their items one after another either horizontally or vertically.

For example:

``` suneido
(Horz a b c d)
```

would simply place items a, b, c, d side by side.

A "Fill" item with a minimum size of 0 and a stretch of >0 can be used to center or right/bottom justify, for example:

``` suneido
(Horz a Fill b)     // left justify a, right justify b
(Horz Fill a Fill)  // center a
```

The stretch factors can be used to divide space proportionally:

``` suneido
(Horz (Pane xstretch: 1) (Pane xstretch: 3)
```

would make the right hand pane three times as wide as the left hand one.

Of course containers can be combined and nested:

``` suneido
(Vert (Horz a b) c)
```

would place a and b side by side, with c below them.

The boxes and stretch pattern is used for user interface and report layout.

See TEX by Knuth. Also see Composite in Design Patterns by Gamma, Helm, Johnson, Vlissides.