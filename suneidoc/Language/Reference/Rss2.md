<div style="float:right"><span class="toplinks"><a href="Rss2/Methods">Methods</a></span></div>

### Rss2

An RSS 2 feed generator.

For example:

``` suneido
r = new Rss2
r.Channel(title: 'My Feed', link: 'mysite.com/feed', description: 'My Blog')
r.AddItem(title: 'Todays Post', description: "Ramblings from today.")
r.ToString()
```

would produce:

``` suneido
<?xml version="1.0" encoding="us-ascii"?>
<rss version="2.0">
<channel>
<title>My Feed</title>
<link>mysite.com/feed</link>
<description>My Blog</description>
<item>
<title>Todays Post</title>
<description>Ramblings from today.</description>
<guid>http://suneido.com/20071019_1202</guid>
</item>
</channel>
</rss>
```

To use the results you would need to upload it as a file to a web server or access it via Suneido's [HttpServer](<HttpServer.md>).

See also: [Atom](<Atom.md>)