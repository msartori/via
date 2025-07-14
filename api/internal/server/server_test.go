package server

/*
func TestStartServer_GracefulShutdown(t *testing.T) {

	testutil.InjectNoOpLogger()

	// Execute server in go routine
	go StartServer(cfg)

	// Waits for server start
	time.Sleep(200 * time.Millisecond)

	// Sends a SIGINT signal to simulate Ctrl+C (shutdown)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	// Waits for shutdown completes
	time.Sleep(500 * time.Millisecond)

}*/
