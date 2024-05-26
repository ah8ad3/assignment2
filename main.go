package main

func main() {
	var stopC chan bool
	server := newServer(stopC)
	go server.Run()
	<-stopC
}
