@echo off
rm -rf build/windows
go build -o build/windows/wordo.exe
copy %cd%\config.json %cd%\build\windows