### Heading

``` suneido
(field) => string
```

Returns the Heading for a field from its Field_ definition. If there is no Heading it returns the Prompt. If there is no Heading or Prompt it returns the field.

For example:

``` suneido
Heading("state_prov")
    => "State/
        Province"
```


See also:
[Datadict](<Datadict.md>),
[Prompt](<Prompt.md>),
[PromptOrHeading](<PromptOrHeading.md>),
[SelectPrompt](<SelectPrompt.md>)
