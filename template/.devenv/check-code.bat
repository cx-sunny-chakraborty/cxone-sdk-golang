@echo off

:: --------------------------
:: LINT THE CODE
:: --------------------------
if %1.==lint. (
	if %2.==fix. (
		echo ####################################################################
		echo ## Check and fix code format - lint
		echo ####################################################################
		pushd "%~dp0"
		cd ..
        goimports -w .
        golangci-lint run ./... --fix
		popd
	) else (
		echo ####################################################################
		echo ## Check code format - lint
		echo ####################################################################
		pushd "%~dp0"
		cd ..
		golangci-lint run ./... 
		popd
	)
) else (
	echo No valid command passed
	echo Use "lint", "lint fix"
)