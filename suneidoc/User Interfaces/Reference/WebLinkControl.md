### WebLinkControl

``` suneido
(text, url, prefix = 'http://')
```

Similar to [StaticControl](<StaticControl.md>) but the text is displayed in a link style (blue, underlined) and clicking on it opens your browser for the specified url.

You can change the control to bring up email software by passing in "mailto:" as the prefix.

For example:

``` suneido
Window(#(WebLink, 'Abc', 'www.abc.com'))
```

would display:

>	<span style="font-family: Arial; color: blue"><u>Abc</u></span>