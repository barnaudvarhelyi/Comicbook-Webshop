@echo off

set ExeName=test
echo Building %ExeName%...

set GOOS=windows

if not exist .\build (
	mkdir .\build
)

if exist .\build\%ExeName%.exe (
	del .\build\%ExeName%.exe
)

@echo on
go build -v -o .\build\%ExeName%.exe

@echo off
if exist .\build\%ExeName%.exe (
	pushd .\build
	.\%ExeName%.exe
	popd
)