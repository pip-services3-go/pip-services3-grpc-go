#!/usr/bin/env pwsh

$component = Get-Content -Path "component.json" | ConvertFrom-Json
$buildImage="$($component.registry)/$($component.name):$($component.version)-$($env:TRAVIS_BUILD_NUMBER)-build"
$testImage="$($component.registry)/$($component.name):$($component.version)-$($env:TRAVIS_BUILD_NUMBER)-test"

# Clean up build directories
Get-ChildItem -Path "." -Include "obj" -Recurse | foreach($_) { Remove-Item -Force -Recurse $_.FullName }
Get-ChildItem -Path "." -Include "node_modules" -Recurse | foreach($_) { Remove-Item -Force -Recurse $_.FullName }

# Remove docker images
docker rmi $buildImage --force
docker rmi $testImage --force
docker image prune --force

# Remove existed containers
docker ps -a | Select-String -Pattern "Exit" | foreach($_) { docker rm $_.ToString().Split(" ")[0] }
