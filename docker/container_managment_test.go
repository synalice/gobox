package docker_test

//var config = docker.ContainerConfig{
//	Image:      "python",
//	LocalImage: false,
//	Cmd:        []string{"python"},
//	TimeLimit:  2 * time.Second,
//	MemoryMB:   64,
//}
//
//func TestNewController(t *testing.T) {
//	_, err := docker.NewController(&config)
//	if err != nil {
//		t.Error(err)
//		t.FailNow()
//	}
//}
//
//func TestRun(t *testing.T) {
//	c, err := docker.NewController(&config)
//	if err != nil {
//		t.Error(err)
//	}
//
//	_, logs, err := c.Run()
//	if err != nil {
//		if err.Error() != "context deadline exceeded" {
//			t.Error(err)
//		}
//		fmt.Println("Container killed due to timeout")
//	}
//
//	fmt.Println(logs)
//}
