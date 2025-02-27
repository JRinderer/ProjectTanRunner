@echo off
net session >nul 2>&1
if %errorlevel%==0 (
  goto Setup
) else (
  goto AdminErr
)
:Setup
@echo on
"%~dp0SetupClient.exe" /S
@echo off
if %errorlevel% ==0 (
  reg add "HKLM\SOFTWARE\WOW6432Node\Tanium\Tanium Client\Sensor Data\Tags" /v NonDomainJoined /t REG_SZ /d "Added via Custom Installer NDJ-NC-PAC-1"
  reg add "HKLM\SOFTWARE\WOW6432Node\Tanium\Tanium Client\Sensor Data\Tags" /v NonCorporate /t REG_SZ /d "Added via Custom Installer NDJ-NC-PAC-1"
  reg add "HKLM\SOFTWARE\WOW6432Node\Tanium\Tanium Client\Sensor Data\Tags" /v PacApts /t REG_SZ /d "Added via Custom Installer NDJ-NC-PAC-1"
  echo Setup complete.
) else (
  echo Setup failed.
)
goto Exit
:AdminErr
echo Install script must be run as administrator.
:Exit
pause