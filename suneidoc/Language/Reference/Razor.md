### Razor

``` suneido
(string, context) => string
```

A template processor patterned after Microsoft's ASP.NET MVC Razor.

Razor templates are nested sections of HTML and Suneido code. The '@' character is used to indicate the start of code.

`@name`
: Will be replaced by the value of the variable (HTML encoded). To reference a context member use .name

`@( *expression* )`
: Will replaced by the value of the expression (HTML encoded). To reference a context member use .name

`@{ `*`code`*` }`
: Evaluates code for it's side effects e.g. setting a variable.

`@if (...) { ... }`
: 

`@if (...) { ... } else { ... }`
: 

`@for (...) { ... }`
: 

`@while (...) { ... }`
: 

For example:

``` suneido
s = "<h1>Hello @.name</h1>"
context = #(name: "Fred")
Razor(s, context)
    => "<h1>Hello Fred</h1>"
```

Internally there are two phases, first the template is converted to Suneido source code. Next the source code is evaluated to produce the result.  You can run just the first phase using Razor.Translate(string) which returns the source code.

**Note:** variables and expressions are HTML encoded by default. To avoid encoding, wrap a string with HtmlString. For example, if you had a helper function that returned HTML for an link, it should do something like: `return HtmlString(link)`