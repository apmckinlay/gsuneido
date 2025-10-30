### Ftsearch

Full text index and search. Uses a Snowball stemmer and BM25 scoring.
`Ftsearch.Tokens(string) => object`
: Returns a list of processed tokens from the string. This is for debugging.

``` suneido
Ftsearch.Tokens("business selling is 343-8887 401k foo123")
=> #("busi", "sell", "343-8887", "343", "8887", "401k", "401", "foo123", "foo", "123")
```

`Ftsearch.Create() => builder`
: Creates a "builder" that is used to build a search index.

`builder.Add(id, title, text)`
: The id should be a small integer. Gaps in the numbering will increase memory usage.
: The title is indexed along with the text, but the words in it are treated as more important.

`builder.Pack() or index.Pack() => data`
: Returns the index as a single binary string.

`FtSearch.Load(data) => index`
: Loads the index so it can be searched.

`index.Search(query, scores = false) => list`
: Searches the index for the words in the query. It returns a list of the top 20 matching document ids, highest scoring first. If scores is true, the results are objects containing id and score.

`index.Update(id, oldTitle, oldText, newTitle, newText)`
: If oldTitle and oldText aren't empty it removes them from the index.
Note: They should exactly match what was previously added to the index.
If newTitle and newText aren't empty it adds them to the index. (Like build.Add)
Note: index.Update is less efficient than builder.

For example:

``` suneido
builder = Ftsearch.Create()
builder.Add(0, "First", "Now is the time")
builder.Add(1, "Second", "This is not the first document")
data = builder.Pack()
index = Ftsearch.Load(data)
Print(index.Search("nonexistent"))
Print(index.Search("TIME"))
Print(index.Search("first", scores:))
Print(index.Search("first time", scores:))
=>	#()
	#(0)
	#(#(id: 0, score: .2865053035333572), #(id: 1, score: .16044296997868))
	#(#(id: 0, score: .9796524840933026), #(id: 1, score: .16044296997868))
```

#### Limits

id's must be in the range 0 to 64k

The maximum number of unique terms (stemmed words) is 64k

The highest count of a particular term in a document is 255.
After that the count just sticks at 255.

Single letters or digits are not indexed.

Terms longer than 32 characters are not indexed.