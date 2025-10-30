### [XmlNode](<../XmlNode>) - Methods
`Name() => string`
: Returns the name (tag) of the node (or "" for text nodes).

`Attributes() => object`
: Returns an object with named members containing the attribute values.

`Children() => list`
: Returns a list of child nodes.

`Text() => string`
: Returns the concatenation of .Text() from all the child nodes, recursively.

`[i] => xmlNode`
: Returns the i'th child node.

`._name => string`
: Returns the value of an attribute.

`.name => xmlNode_list`
: Returns an XmlNode_list of the children with that name. An XmlNode_list behaves similar to an XmlNode except that .name will collect the results of applying .name to each member of the list.

`ToString()`
: Returns an XML string.