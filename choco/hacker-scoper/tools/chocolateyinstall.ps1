$ErrorActionPreference = 'Stop'; # stop on all errors
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

$packageArgs = @{
  packageName   = $env:ChocolateyPackageName
  destination   = "$toolsDir"
  file          = "$toolsDir\hacker-scoper_2.1.0_windows_386.zip"
  file64        = "$toolsDir\hacker-scoper_2.1.0_windows_amd64.zip"
}

Get-ChocolateyUnzip @packageArgs
Remove-Item -Path $packageArgs.file,$packageArgs.file64