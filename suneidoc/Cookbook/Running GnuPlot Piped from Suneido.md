## Running GnuPlot Piped from Suneido

**Category:** Interfacing

**Ingredients**

-	creating a text file
-	accessing the system database tables
-	getting information about the current working directory 
-	calling external programs
-	GnuPlot


**Problem**

You want to create charts using GnuPlot to represent data stored in Suneido database.

**Recipe**

Suneido doesn't have much support for charts. GnuPlot is a command line program that can generate different types of charts from text files as well as through  interactive commands given through stdin.

GnuPlot doesn't have much ability to manipulate and transform data. Suneido can be programmed to output the desired data onto text files, which can be used by GnuPlot to generate charts. In addition, Suneido provides the ability to pipe stdin to external programs. This facility can be used to communicate with GnuPlot interactively. Given below is the code for a Control that demonstrate these abilities.

``` suneido
Controller
    {
    Title: "Gnu Plot"
    New()
        {
        ExportTab('tables project tablename, totalsize', "gpdata")
        .gp = RunPiped("gnuplot.exe")
        .gp.Writeline('set loadpath ' $ Display(GetCurrentDirectory()))
        .gp.Writeline('plot "gpdata" using 2:xticlabels(1) notitle with boxes')
        }
    Controls: #(Horz
                    (Field )
                    (Button "Plot")
                )
    On_Plot()
        {
        .gp.Writeline(.Horz.Field.Get())
        }
    Destroy()
        {
        .gp.Close()
        super.Destroy()
        }
    }
```

GnuPlot will bring up a window showing a bar chart (assuming that Gnuplot is installed in your computer and in PATH) showing the size of each table in the database. Simultaneously a window with a text input field and a button named Plot also appears. You may enter any valid GnuPlot command (like plot sin(x)) in the text field and press the button to update the chart that remains in screen.

The first line in the New() method of our Control, the line 
``` suneido
ExportTab('tables project tablename, totalsize', "gpdata")
```
 generates the data file for GnuPlot. For specific needs one probably would want to generate custom files using the File command.

 The next line in the New() method starts the GnuPlot using the RunPiped command and assigns it to the member, gp. There are three commands to run an external progam from Suneido viz. System, Spawn and RunPiped. System and Spawn allow a one time execution of the argument passed. They don't allow continued communication with the external program. RunPiped is the option suited for running external programs from Suneido when continued interaction with it is required. To be able to use an external program via RunPiped, it should support piped commands which GnuPlot does. This line apparently does nothing. However, if you check the Task Manager, you will see GnuPlot as a process. RunPiped may be written in two forms - as shown in my example or as a block. If run as a block, the process won't be available outside the block to allow continued communication with the program. Hence I did not opt for the block form. Another beneficial side effect of using the RunPiped command is that the command window does not appear, which would have appeared if System or Spawn was used. Hence, even in situations where one does not want continued interaction with GnuPlot, RunPiped command may be preferred for running GnuPlot. In such situations the block form may be preferred as one need not remember to close the process.

Next we give GnuPlot our first command using the WriteLine() method of RunPiped. We use WriteLine() and not Write() method because GnuPlot expects each command as a line. We inform GnuPlot the current working directory of Suneido (where Suneido has output the file made using ExportTab() command), so that it will search this directory for the datafile when we specify it in the next command. 

Using the next WriteLine() method we build our first chart. The chart is not pretty, included only to show a chart on startup.

The On_Plot() method is called whenever we press the Plot button in our control. It simply hands over the text in the text field to GnuPlot as a command using, once again, the WriteLine() method. There is no error checking by the control as to the validity of the GnuPlot commands. GnuPlot ignores most errors in commands, executing only valid commands.

As the RunPiped was not used in the block form, we should close the process. We do this in the Destroy() method where we call Close() method of RunPiped. We could call ExitValue() method also, which would return the value returned by the external program. This method may be useful while debugging.

The demo control is only for demonstration of the ability of  Suneido to run an external program and interact with it as long as required, if it supports piped commands. Suneido can not only communicate to the external program using Write() or WriteLine(), but can also receive communications from the external program using Read() and ReadLine() methods of RunPiped. To run GnuPlot, this was not required, hence I have not used it.

**See Also**

[GnuPlot](<http://sourceforge.net/projects/gnuplot>)