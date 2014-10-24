#!/bin/bash

printf "** Building linux/386\n"
go-linux-386 build -a -o bin/linux-386/check_pa_cpu github.com/zerklabs/check_pa_cpu

printf "** Building linux/amd64\n"
go-linux-amd64 build -a -o bin/linux-amd64/check_pa_cpu github.com/zerklabs/check_pa_cpu

printf "** Building windows/386\n"
go-windows-386 build -a -o bin/windows-386/check_pa_cpu.exe github.com/zerklabs/check_pa_cpu

printf "** Building windows/amd64\n"
go-windows-amd64 build -a -o bin/windows-amd64/check_pa_cpu.exe github.com/zerklabs/check_pa_cpu
