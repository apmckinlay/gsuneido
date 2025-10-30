### TranslateLanguage

``` suneido
(string [ , string ... ])
```

Uses the translatelanguage database table and the Suneido.Language setting
to translate strings to another language.

If Suneido.Language is not set, or if it is set to "english",
then TranslateLanguage simply returns the string passed to it.
If Suneido.Language is set to another language, for example:

``` suneido
Suneido.Language = #(name: "spanish", charset: "DEFAULT")
```

then TranslateLanguage will look for the string to be translated
in the trlang_from field of the translatelanguage table.
If it does not find it, the string is returned untranslated.
If it does find it, then the value of the column 
corresponding to the current language is examined.
If that column is empty, the string is returned untranslated,
otherwise the value of the column is the translation.
For example, if the language was set to "spanish" as above,
the column used would be trlang_spanish.
For example:

``` suneido
TranslateLanguage("Replace")
    => "Reemplazar"
```

To reduce the number of translatelanguage table entries,
some extra steps are done during translation.

-	If the string to be translated ends with "..."
		this is stripped off before looking it up
		and then added back on to the translation.
		trlang_from values should *not* end with "...".
		This allows the same entry to be used for menu items with "..."
		and the corresponding dialog title without "...".
-	If the string to be translated contains "&"
		(i.e. menu shortcut keys)
		these are removed before looking it up.
		trlang_from value should *not* contain "&".
-	If the translated string contains "&"
		but the original string did not,
		then any "&" are removed from the translation.


For example:

``` suneido
TranslateLanguage("Replace...")
    => "Reemplazar..."

TranslateLanguage("&Replace")
    => "R&eemplazar"
```

**Parameters**

Strings to be translated can also contain *parameters*,
either numbered (e.g. %1) or named (e.g. %file).
These are replaced by the corresponding argument
passed to TranslateLanguage.
For example:

``` suneido
TranslateLanguage("Can't open %1 for %2", "output.txt", "append")
    => "Can't open output.txt for append"
```

or:

``` suneido
TranslateLanguage("Can't open %file for %mode", 
    mode: "append", file: "output.txt")
    => "Can't open output.txt for append"
```

Parameter substitution is always done, 
even if no language translation is done.
Normally, translations should include the same parameters
as the English version,
although not necessarily in the same order.