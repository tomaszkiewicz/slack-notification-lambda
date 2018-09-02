@echo off

echo Building...

set GOOS=linux
set GOARCH=amd64
go build -o main ./cmd/lambda

echo Packaging...

%USERPROFILE%\Go\bin\build-lambda-zip.exe -o lambda.zip main

echo Publishing...

rem aws lambda update-function-code --function-name slack-notification-lambda --zip-file fileb://lambda.zip

echo Cleaning...

del main
rem del lambda.zip

echo All done

pause
