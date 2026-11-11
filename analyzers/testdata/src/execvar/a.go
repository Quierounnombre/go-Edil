package execvar

import "os/exec"

func bad(name string, script string) {
	exec.Command(name, "x")          // want `command name must be a constant`
	exec.Command("sh", "-c", script) // want `non-constant script`
}

func good(branch string) {
	exec.Command("git", "checkout", branch)
	exec.Command("sh", "-c", "echo hi")
}
