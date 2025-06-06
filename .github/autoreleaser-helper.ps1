echo $pwd
$originaldir = (pwd).path
echo $(tree)

echo 'Downloading the releases file...'
Invoke-WebRequest -Uri https://api.github.com/repos/itsignacioportal/hacker-scoper/releases/latest -OutFile $env:TEMP\releases.json

echo 'Installing jq...'
choco install jq

echo 'Parsing download URL from JSON...'
$cmdOutput = type $env:TEMP\releases.json | C:\ProgramData\chocolatey\bin\jq.exe '.assets[11].browser_download_url'

echo 'Downloading the windows_386 file...'
$cmdOutput = $cmdOutput -replace '"',''
Invoke-WebRequest -Uri $cmdOutput -OutFile $env:TEMP\windows_386.tar.gz

echo 'Parsing download URL from JSON...'
$cmdOutput = type $env:TEMP\releases.json | C:\ProgramData\chocolatey\bin\jq.exe '.assets[12].browser_download_url'

echo 'Downloading the windows_amd64 file...'
$cmdOutput = $cmdOutput -replace '"',''
Invoke-WebRequest -Uri $cmdOutput -OutFile $env:TEMP\windows_amd64.tar.gz

echo 'Extracting files...'
cd $env:TEMP
mkdir windows_386
cd windows_386
tar -xvzf ..\windows_386.tar.gz

cd $env:TEMP
mkdir windows_amd64
cd windows_amd64
tar -xvzf ..\windows_amd64.tar.gz

echo 'Parsing latest version tag from JSON...'
$version = type $env:TEMP\releases.json | C:\ProgramData\chocolatey\bin\jq.exe '.tag_name'
$version = $version -replace '"',''

echo 'Preparing Chocolatey file...'
cd $originaldir
echo $pwd
tree
Copy-Item choco\chocolateyinstall_template.ps1 choco\hacker-scoper\tools\chocolateyinstall.ps1
$filePath = "choco\hacker-scoper\tools\chocolateyinstall.ps1"
(Get-Content $filePath).Replace("VERSIONHERE",$version) | Set-Content $filePath

Copy-Item choco\hacker-scoper_template.nuspec choco\hacker-scoper\hacker-scoper.nuspec
$filePath = "choco\hacker-scoper\hacker-scoper.nuspec"
$version = $version -replace 'v',''
(Get-Content $filePath).Replace("VERSIONHERE",$version) | Set-Content $filePath

echo 'Compressing files...'
Compress-Archive $env:TEMP\windows_386\hacker-scoper.exe -DestinationPath choco\hacker-scoper\tools\hacker-scoper_$($version)_windows_386.zip
Compress-Archive $env:TEMP\windows_amd64\hacker-scoper.exe -DestinationPath choco\hacker-scoper\tools\hacker-scoper_$($version)_windows_amd64.zip

cd choco\hacker-scoper
echo $pwd