mkdir build
go build shex.go
cd compiler
go build compile.go
move ..\shex.exe ..\build
move .\compile ..\build