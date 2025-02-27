
func deleteExecutable() error {
exePath, err := os.Executable()
if err != nil {
return fmt.Errorf("failed to get executable path: %v", err)
}

cmd := exec.Command("cmd", "/c", "del", exePath)
cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}

err = cmd.Start()
if err != nil {
return fmt.Errorf("failed to start delete command: %v", err)
}

err = cmd.Process.Release()
if err != nil {
return fmt.Errorf("failed to detach process: %v", err)
}

return nil
}