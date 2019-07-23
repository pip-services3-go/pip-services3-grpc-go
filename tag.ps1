#!/usr/bin/env pwsh

Set-StrictMode -Version latest
$ErrorActionPreference = "Stop"

# Get component data and set necessary variables
$component = Get-Content -Path "component.json" | ConvertFrom-Json
$tag="v$($component.version)-$($component.build)"

# Define group name
$pos = $component.registry.IndexOf("/")
$groupName = ""
if ($pos -gt 0) {
    $groupName = $component.registry.Substring($pos + 1, $component.registry.Length - $pos - 1)
}

# Change git remote so git use ssh on push
git remote set-url origin git@gitlab.com:$groupName/$($component.name).git

git add ./obj/*
git add ./component.json
git commit -m "project build by GitLab CI [skip ci]"

# Set git tag
git tag $tag -a -m "Generated tag from GitLabCI for build #$($component.build) [ci skip]"
git push --tags origin HEAD:master 