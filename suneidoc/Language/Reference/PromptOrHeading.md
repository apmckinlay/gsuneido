### PromptOrHeading

``` suneido
( field ) => string
```

Returns the Prompt for a field, if it has one,
or else the Heading
from its Field_ definition.

For example:

``` suneido
PromptOrHeading("firstname")
    => "First Name"
```


See also:
[Datadict](<Datadict.md>),
[Heading](<Heading.md>),
[Prompt](<Prompt.md>),
[SelectPrompt](<SelectPrompt.md>)
