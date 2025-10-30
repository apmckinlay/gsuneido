### SelectPrompt

``` suneido
(field) => string
```

Returns the SelectPrompt for a field from its Field_ definition. If it does not exist it looks for non "" Prompt and then Heading. Otherwise it returns the field.

The SelectPrompt is used for e.g. AccessControl Select to differentiate fields that may have the same Prompt and Heading. Normally it is not necessary to define it.

For example:

``` suneido
SelectPrompt("date")
    => "Date"
```


See also:
[Datadict](<Datadict.md>),
[Heading](<Heading.md>),
[Prompt](<Prompt.md>),
[PromptOrHeading](<PromptOrHeading.md>)
