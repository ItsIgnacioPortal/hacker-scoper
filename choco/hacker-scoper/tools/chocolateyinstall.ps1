$ErrorActionPreference = 'Stop'; # stop on all errors
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

$packageArgs = @{
  packageName   = $env:ChocolateyPackageName
  destination   = "$toolsDir"
  file          = "$toolsDir\hacker-scoper_2.1.0_windows_386.zip"
  file64        = "$toolsDir\hacker-scoper_2.1.0_windows_amd64.zip"

  # Checksums are now required as of 0.10.0.
  checksum      = '88B76712B7857C0B1AC652F3CB5B768080FA3E35C41636516C822AD2E1141B5F'
  checksumType  = 'sha256' #default is md5, can also be sha1, sha256 or sha512
  checksum64    = 'A870587BD25C4E837535312B972108ED32BB1A6685E1C72771B5BF7F7EC15CF0'
  checksumType64= 'sha256' #default is checksumType
}

Get-ChocolateyUnzip @packageArgs
Remove-Item -Path $packageArgs.file,$packageArgs.file64