//go:build windows

package monitor

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// NewProvider discovers RDP sessions with quser and SSH sessions from sshd processes.
func NewProvider() Provider {
	return windowsProvider{}
}

type windowsProvider struct{}

func (windowsProvider) Snapshot() (Snapshot, error) {
	sessions := make([]Session, 0)
	rdpSessions, err := rdpSessions()
	if err != nil {
		return Snapshot{}, err
	}
	sessions = append(sessions, rdpSessions...)
	sessions = append(sessions, windowsSSHSessions()...)
	sessions = append(sessions, wslSessions()...)
	return Snapshot{Sessions: sessions, At: time.Now()}, nil
}

func rdpSessions() ([]Session, error) {
	output, err := runHidden("quser")
	if err != nil {
		return nil, fmt.Errorf("quser: %w", err)
	}

	now := time.Now()
	var sessions []Session
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		fields := strings.Fields(strings.TrimSpace(strings.TrimPrefix(scanner.Text(), ">")))
		if len(fields) < 4 || strings.EqualFold(fields[0], "USERNAME") || !strings.HasPrefix(strings.ToLower(fields[1]), "rdp-tcp") {
			continue
		}
		startedAt := now
		if len(fields) > 5 {
			startedAt = parseQuserTime(strings.Join(fields[5:], " "), now)
		}
		sessions = append(sessions, Session{
			Kind:      SessionKindRDP,
			User:      fields[0],
			Source:    fields[1],
			ID:        fields[2],
			StartedAt: startedAt,
		})
	}
	return sessions, scanner.Err()
}

func parseQuserTime(value string, fallback time.Time) time.Time {
	for _, layout := range []string{
		"02.01.2006 15:04",
		"1/2/2006 3:04 PM",
		"01/02/2006 3:04 PM",
		"1/2/2006 15:04",
	} {
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return parsed
		}
	}
	return fallback
}

func sshSessions() []Session {
	command := "Get-CimInstance Win32_Process -Filter \"Name='sshd.exe'\" | ForEach-Object { $owner = Invoke-CimMethod -InputObject $_ -MethodName GetOwner; '{0}|{1}' -f $_.ProcessId, $owner.User }"
	output, err := runHidden("powershell.exe", "-NoProfile", "-Command", command)
	if err != nil {
		return nil
	}

	var sessions []Session
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		fields := strings.SplitN(strings.TrimSpace(scanner.Text()), "|", 2)
		if len(fields) != 2 {
			continue
		}
		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		user := strings.TrimSpace(fields[1])
		if !isInteractiveUser(user) {
			continue
		}
		sessions = append(sessions, Session{Kind: SessionKindSSH, User: user, Source: "OpenSSH", ID: strconv.Itoa(pid)})
	}
	return sessions
}

func windowsSSHSessions() []Session {
	command := `$logins = @{}
Get-WinEvent -LogName 'OpenSSH/Operational' -MaxEvents 2000 | Sort-Object TimeCreated | ForEach-Object {
    $message = $_.Message
    if ($message -match 'Accepted .* for (?<user>\S+) from (?<source>\S+) port (?<port>\d+)') {
        $key = "$($matches['source']):$($matches['port'])"
		$logins[$key] = "$($matches['user'])|$($_.TimeCreated.ToUniversalTime().ToString('o'))"
    }
}
$connections = Get-NetTCPConnection -State Established | Where-Object { $_.LocalPort -eq 22 }
foreach ($connection in $connections) {
	$key = "$($connection.RemoteAddress):$($connection.RemotePort)"
	if ($logins.ContainsKey($key)) {
		$parts = $logins[$key] -split '\|', 2
		"$($parts[0])|$key|$($parts[1])"
	}
}`
	output, err := runHidden("powershell.exe", "-NoProfile", "-Command", command)
	if err != nil {
		return nil
	}

	var sessions []Session
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		fields := strings.SplitN(strings.TrimSpace(scanner.Text()), "|", 3)
		if len(fields) != 3 || !isInteractiveUser(fields[0]) {
			continue
		}
		startedAt, err := time.Parse(time.RFC3339Nano, fields[2])
		if err != nil {
			startedAt = time.Now()
		}
		sessions = append(sessions, Session{
			Kind:      SessionKindSSH,
			User:      fields[0],
			Source:    "OpenSSH",
			ID:        "WindowsSSH:" + fields[1],
			StartedAt: startedAt,
		})
	}
	return sessions
}

func isInteractiveUser(user string) bool {
	if user == "" {
		return false
	}
	upper := strings.ToUpper(user)
	return upper != "SYSTEM" && upper != "LOCAL SERVICE" && upper != "NETWORK SERVICE" && !strings.HasPrefix(upper, "DWM-") && !strings.HasPrefix(upper, "UMFD-")
}

func wslSessions() []Session {
	output, err := runHidden("wsl.exe", "-e", "ps", "-eo", "user=,pid=,etimes=,args=")
	if err != nil {
		return nil
	}

	var sessions []Session
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}
		commandLine := strings.Join(fields[3:], " ")
		user := wslSSHUser(commandLine)
		if user == "" {
			continue
		}
		elapsedSeconds, err := strconv.Atoi(fields[2])
		if err != nil {
			continue
		}
		sessions = append(sessions, Session{
			Kind:      SessionKindSSH,
			User:      user,
			Source:    "WSL",
			ID:        "WSL:" + fields[1],
			StartedAt: time.Now().Add(-time.Duration(elapsedSeconds) * time.Second),
		})
	}
	return sessions
}

func wslSSHUser(commandLine string) string {
	const marker = "sshd: "
	start := strings.Index(commandLine, marker)
	if start < 0 {
		return ""
	}
	value := commandLine[start+len(marker):]
	separator := strings.Index(value, "@pts/")
	if separator <= 0 {
		return ""
	}
	user := strings.TrimSpace(value[:separator])
	if !isInteractiveUser(user) {
		return ""
	}
	return user
}

func runHidden(name string, args ...string) ([]byte, error) {
	command := exec.Command(name, args...)
	command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return command.Output()
}
