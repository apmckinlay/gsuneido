### Plugins

A class implementing a plug-in registry.

Scans libraries for plug-in definitions - records with names starting with "Plugin_". Definitions should look like:
<pre>
#(
ExtensionPoints:
    (
    ('<i>extension_point</i>')
    ...
    )
Contributions:
    (
    ('<i>Plugin</i>', '<i>extension_point</i>', ...)
    ...
    )
)
</pre>

Where *Plugin* is the name of another plug-in definition (without the "Plugin_" prefix).

Contributions are checked that they are for a valid plug-in and extension_point.

The Plugins registry is normally accessed with Plugins(). This calls the Plugins CallClass which returns Suneido.Plugins, creating it if it doesn't exist.

#### Methods:
Plugins.Clear()
: Removes the plugin registry - it will be re-read next time it is accessed. Called by LibView when a library is Use'd or Unuse'd.

Plugins().Reset()
: Re-scans the libraries for plug-in definitions. Reset() should be called after you make changes or additions to definitions. If there are errors they are shown in an alert.

Plugins().Plugins()
: Returns a list of plug-in names.

Plugins().Contributions(plugin, extenpt = false)
: Returns a list of contributions for the specified plug-in (and extension point if specified).

Plugins().ForeachContribution(plugin, extenpt, block)
: Calls the block for each contribution for the specified plug-in and extension point.

Plugins().Errors()
: Returns a list of errors found. 

Plugins().ShowErrors()
: Displays any errors. 


See also:
[Contributions](<Contributions.md>),
[GetContributions](<GetContributions.md>),
[LastContribution](<LastContribution.md>),
[OptContribution](<OptContribution.md>),
[SoleContribution](<SoleContribution.md>)
