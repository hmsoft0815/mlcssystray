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

[Languages]
Name: "en"; MessagesFile: "compiler:Default.isl"
Name: "de"; MessagesFile: "compiler:Languages\German.isl"

[CustomMessages]
en.AutostartTask=Start %1 when I sign in to Windows
de.AutostartTask=%1 bei der Windows-Anmeldung starten

[Tasks]
Name: "autostart"; Description: "{cm:AutostartTask,{#MyAppName}}"; Flags: unchecked

[Files]
Source: "{#SourceExe}"; DestDir: "{app}"; DestName: "mlcsshtrayagent.exe"; Flags: ignoreversion

[Icons]
Name: "{autoprograms}\MLC SSH Tray Agent"; Filename: "{app}\mlcsshtrayagent.exe"

[Registry]
Root: HKCU; Subkey: "Software\Microsoft\Windows\CurrentVersion\Run"; ValueType: string; ValueName: "MLC SSH Tray Agent"; ValueData: """{app}\mlcsshtrayagent.exe"""; Flags: uninsdeletevalue; Tasks: autostart

[Run]
Filename: "{app}\mlcsshtrayagent.exe"; Description: "{cm:LaunchProgram,{#MyAppName}}"; Flags: nowait postinstall skipifsilent
