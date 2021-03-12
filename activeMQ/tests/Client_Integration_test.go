package tests

import (
	"testing"
	content "twitterApiTest/activeMQ"
	"twitterApiTest/activeMQ/client/impl"
)

func TestSendAndReceive(t *testing.T) {
	client := impl.StompClient{}
	err := client.Connect("localhost:61613")
	checkForError(t,err)
	channel := make(chan []byte)
	err = client.SubscribeToQueue("test",channel)
	checkForError(t,err)
	want := "Integration test :D"
	client.SendMessageToQueue("test",content.TEXT,[]byte(want))
	got := string(<- channel)
	if got != want {
		t.Errorf("Got %s want %s",got, want)
	}
	client.Disconnect()
}


func checkForError(t testing.TB, err error){
	t.Helper()
	if err != nil {
		t.Fatalf("Could not subscribe to queue because of error: %s",err.Error())
	}
}
