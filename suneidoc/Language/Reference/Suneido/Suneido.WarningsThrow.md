<div style="float:right"><span class="builtin">Builtin</span></div>

#### Suneido.WarningsThrow

``` suneido
(arg = true)
```

Controls whether certain built-in warnings throw exceptions. By default, these warnings just log. But to find them it can be helpful to make them throw.

If the argument is:
`true`
: all the warnings will throw exceptions

`false`
: none of the warnings will throw exception (just log)

`regular expression`
: warnings that match will throw exceptions

The warnings include: (as of 2024-01-11)

-	`Sort! functions should return true or false`
-	`object named size > ...`
-	`object list size > ...`
-	`project-map large > ...`
-	`project-map derived large > ...`
-	`summarize-map large > ...`
-	`temp index large > ...`
-	`temp index derived large > ...`
-	`query where explode large > ...`