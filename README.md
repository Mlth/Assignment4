To run the program:
- Open five terminals
- Type 'go run peer.go x y', where x is the portnumber minus 5500, and y is the amount of peers you want to run.
The x-values have to start from 0 and be consecutive. Here is an example:
''' 
go run peer.go 0 5
go run peer.go 1 5
go run peer.go 2 5
go run peer.go 3 5
go run peer.go 4 5
'''
- You can now type either yes or no, depending on whether you want a peer to enter the critical state
- You can then type either yes or no, depending on whether you want a peer to levae the critical state

All operations can be found in the OutputLog.log file.

If you are using Windows or linux, you can also use one of the script files to start the program with 5 peers:
WINDOWS:
    double click or run the run.bat
LINUX
    in a terminal type "sh run.sh"