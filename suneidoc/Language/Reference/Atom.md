<div style="float:right"><span class="toplinks"><a href="Atom/Methods">Methods</a></span></div>

### Atom

An Atom feed generator.

For example:

``` suneido
a = new Atom
a.Feed(title: 'My Feed', link: 'http:/mysite.com/myfeed.html')
a.AddEntry(title: 'Entry One', link: 'http:/mysite.com/myentry1.html', 
    author: 'Andrew', content: 'my text')
a.ToString()
```

would produce:

``` suneido
<?xml version="1.0" encoding="us-ascii"?>
<feed>
    <title>My Feed</title>
    <link rel="alternate" href="http:/mysite.com/myfeed.html" type="text/html" />
    <id>urn:uuid:8525fc2c-a7b2-4ecd-8ce7-be6fa372f180</id>
    <updated>2009-06-03T15:09:17Z</updated>
<entry>
    <title>Entry One</title>
    <link rel="alternate" href="http:/mysite.com/myentry1.html" type="text/html" />
    <author><name>Andrew</name></author>
    <id>urn:uuid:42500b65-3570-4846-8811-6c1992b5ad27</id>
    <updated>2009-06-03T15:09:17Z</updated>
    <content type="text">my text</content>
</entry>
</feed>
```

To use the results you would need to upload it as a file to a web server or access it via Suneido's [HttpServer](<HttpServer.md>)

See also: [Rss2](<Rss2.md>)