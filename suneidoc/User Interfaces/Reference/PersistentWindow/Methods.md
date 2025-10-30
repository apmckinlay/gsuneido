### [PersistentWindow](<../PersistentWindow.md>) - Methods
`DeleteSet(setname)`
: Removes the persistent set definition.  If you want
to redefine a persistent set, the set must be removed 
first.

`Load(set = "IDE")`
: Runs the persistent window set.  This is used in the Init 
function to run the IDE windows. It is best to run user defined 
persistent sets as a command line argument to suneido.exe from 
the command line or from a shortcut.   
For example:  
`suneido IDE`  
from the command line is equivalent to running
just "suneido".