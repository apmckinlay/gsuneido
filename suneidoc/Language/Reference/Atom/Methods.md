### [Atom](<../Atom>) - Methods
`Feed`
: Specify the main feed element. Accepts title, link, author, id, updated, subtitle, rights, contributor, generator, post, category, icon, and logo attributes. title and link are required.

`AddEntry`
: Add an entry. Accepts title, link, author, id, updated, summary, published, contributor, content, edit, category, rights, and source attributes. title and link are required.

`ToString() => string`
: Returns the XML as a string.

author must be specified either on the feed or on each entry

If content starts with "\<div>" then type is xhtml, else if content starts with "\<" then type is html, else type is text

Currently <u>not</u> supported:

-	entry without a feed
-	multiple author, category, contributor, link
-	author attributes other than name
-	link types other than rel: 'alternate', type: 'text/html'