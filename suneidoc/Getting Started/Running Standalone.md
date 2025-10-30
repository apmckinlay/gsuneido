## Running Standalone

Obviously, you don't want your end users to have to go to the WorkSpace and use Open a Book to run your application. One solution is to create a *persistent set* - a pre-configured set of windows. From the WorkSpace run:

``` suneido
PersistentWindow(#(Book mybook, "My Application") newset: "myset")
```

where `newset:` specifies the name of the persistent set you want to create.  Your IDE windows will close and your book will come up. Closing the book will exit from Suneido.  Now from the command line you can run:

``` suneido
suneido myset
```