### [Rss2](<../Rss2>) - Methods
`Channel`
: Accepts named arguments of: title, link, description, language, copyright, managingEditor, webMaster, pubDate, lastBuildDate, category, generator, docs, cloud, ttl, image, rating, textinput, skipDays, skipHours. title, link and description are required.

`Image`
: Accepts named arguments of: url, title, link, width, height. url, title, and link are required.

`TextInput`
: Requires named arguments of: title, description, name, and link.

`AddItem`
: Accepts named arguments of: title, link, description, author, category, comments, enclosure, guid, pubDate, source. Either title or description is required.

`ToString() => string`
: Returns the XML as a string.

A feed can only have a single Channel, Image, and TextInput but can have multiple items. Image and TextInput are optional.