### Datadict

``` suneido
(field, getMembers = false) => object
```

Returns the Field_ definition for the specified field.
If the Field_ definition is not found, and the field starts with "total_", "max_", "min_", or "average_" then it looks for a definition without that prefix. If no definition is found, Field_string is returned.

For example:

``` suneido
Datadict("dollars")
    => Field_dollars

Datadict("total_dollars")
    => Field_dollars

Datadict("unknown")
    => Field_string
```

If the Field_ definition contains members starting with "Control_", "Format_" or "SelectControl_", they will be "injected" into the Control, Format or SelectControl. This is normally used when inheriting so you don't have to redefine the whole control or format. For example:

``` suneido
Field_date_mandatory

Field_date
    {
    Control_mandatory: true
    }
```

This will add mandatory: true to the control definition. The returned value will be an instance of the class, with the corresponding Control, Format or SelectControl overridden.

getMembers is primarily for use by Prompt, PromptOrHeading, Heading, and SelectPrompt. It avoids the overhead of injection if you don't need the Control, Format or SelectControl.

``` suneido
Datadict('date', #(Prompt)) 
    => #(Prompt: 'Date')
```


See also:
[Heading](<Heading.md>),
[Prompt](<Prompt.md>),
[PromptOrHeading](<PromptOrHeading.md>),
[SelectPrompt](<SelectPrompt.md>)
