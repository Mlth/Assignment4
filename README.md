# To run the program:
1. Open five terminals
2. Type ```'go run peer.go x y'```, 
   
   Where ```x``` is the portnumber minus 5500, and ```y``` is the amount of peers you want to run.

#### The ```x-values``` have to start from 0 and be consecutive. Here is an example:
```go
    go run peer.go 0 5
    go run peer.go 1 5
    go run peer.go 2 5
    go run peer.go 3 5
    go run peer.go 4 5
```
- You can now type either **yes** or **no**, depending on whether you want a peer to enter the critical state or not.

-  You can then type either **yes** or **no**, depending on whether you want a peer to leave the critical state or not.

All operations can be found in the corresponding x_output.log file where **x** is the corresponding port number.

#### NOTE: Optionally if you are using Windows or linux, you may also be able to use one of the script files to start the program for you, with a default of 5 peers.
- WINDOWS: double click and run the run.bat
   
- LINUX: in a terminal type "sh run.sh"