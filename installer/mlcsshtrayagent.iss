#define MyAppName "MLC SSH Tray Agent"
#define MyAppPublisher "Michael Lechner"

[Setup]
AppId={{A1B9D849-04C8-4FBD-8A8E-9350998D665B}
AppName={#MyAppName}
AppVersion={#AppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL=https://github.com/hmsoft0815/mlcssystray
DefaultDirName={localappdata}\Programs\{#MyAppName}
DefaultGroupName={#MyAppName}
DisableProgramGroupPage=yes
PrivilegesRequired=lowest
ArchitecturesAllowed={#AllowedArchitectures}
ArchitecturesInstallIn64BitMode={#InstallModeArchitecture}
OutputDir={#OutputDir}
OutputBaseFilename=mlcsshtrayagent-setup-{#InstallerArch}
SetupIconFile=..\internal\app\assets\icon.ico
UninstallDisplayIcon={app}\mlcsshtrayagent.exe
LicenseFile=..\LICENSE
Compression=lzma2
SolidCompression=yes
WizardStyle=modern

[Tasks]
Name: "autostart"; Description: "Start MLC SSH Tray Agent when I sign in to Windows"; Flags: unchecked

[Files]
Source: "{#SourceExe}"; DestDir: "{app}"; DestName: "mlcsshtrayagent.exe"; Flags: ignoreversion

[Icons]
Name: "{autoprograms}\MLC SSH Tray Agent"; Filename: "{app}\mlcsshtrayagent.exe"

[Registry]
Root: HKCU; Subkey: "Software\Microsoft\Windows\CurrentVersion\Run"; ValueType: string; ValueName: "MLC SSH Tray Agent"; ValueData: """{app}\mlcsshtrayagent.exe"""; Flags: uninsdeletevalue; Tasks: autostart