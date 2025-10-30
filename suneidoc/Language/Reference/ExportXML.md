### ExportXML

``` suneido
(from_query, to_file = false, header = true)
```

Derived from [Export](<Export.md>)

Class for exporting information as an *eXtensible Markup Language *textfile.

For example:

``` suneido
ExportXML('tables', 'tables.txt')
```

Will take the information out of the tables table and output it to the tables.txt file in a format similar to this:

``` suneido
<?xml version="1.0"?>
<!--  suneido xml export  -->
<!DOCTYPE table [
<!ELEMENT table (record*)>
<!ELEMENT record (table?, tablename?, nextfield?, nrows?, totalsize?)>
<!ELEMENT table (#PCDATA)>
<!ATTLIST table type CDATA #REQUIRED>
<!ELEMENT tablename (#PCDATA)>
<!ATTLIST tablename type CDATA #REQUIRED>
<!ELEMENT nextfield (#PCDATA)>
<!ATTLIST nextfield type CDATA #REQUIRED>
<!ELEMENT nrows (#PCDATA)>
<!ATTLIST nrows type CDATA #REQUIRED>
<!ELEMENT totalsize (#PCDATA)>
<!ATTLIST totalsize type CDATA #REQUIRED>
]>
<table>
<record>
<table type="number">2</table>
<tablename type="string">tables</tablename>
<nextfield type="number">5</nextfield>
<nrows type="number">47</nrows>
<totalsize type="number">4821</totalsize>
</record>
...
</table>
```

Specifying the header: produces an infile Document Type Definition.

See also: [ImportXML](<ImportXML.md>)