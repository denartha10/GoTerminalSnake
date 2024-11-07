### Things I learned

In UNIX like operating systems including linux and macosx, there are two modes a terminal can be in there is Raw input mode and Cooked input mode, also known as character and line mode

The terminal driver is by default a by default a line-based system: characters are buffered internally until a carriage return (Enter or Return) before it is passed to the program

terminal can be placed into "raw" mode where the characters are not processed by the terminal driver, but are sent straight through (it can be set that INTR and QUIT characters are still processed). This allows programs like emacs and vi

### Carriage Return Line Feed

In cooked mode the terminal driver will correctly convert \n to Carriage Return + Line Feed. This will set the cursor to the left and at the bottom. In Raw mode however we must manuall include the \r character
