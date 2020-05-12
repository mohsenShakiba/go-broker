package socketserver

//func TestReceiveOpen(t *testing.T) {
//
//	ch := make(chan clientMessage, 10)
//
//	m := mock.NewMockSocket()
//
//	client := socketClient{
//		clientId:        "TEST",
//		isAuthenticated: false,
//		clientType:      0,
//		isClosed:        false,
//		connection:      m,
//		onMessageChan:   ch,
//	}
//
//	msg := "TEST"
//	msgWithPrefix := "TEST"
//	input := []byte(msg)
//
//	_, _ = m.Write(input)
//
//	go client.startReceive()
//
//	res := <-ch
//
//	if string(res.Payload) != msg {
//		t.Fatalf("the input and output of socket mock didn't match input: %s, output: %s", msg, string(res.Payload))
//	}
//
//}
//
//func TestSend(t *testing.T) {
//
//	m := mock.NewMockSocket()
//
//	client := socketClient{
//		clientId:        "TEST",
//		isAuthenticated: false,
//		clientType:      0,
//		isClosed:        false,
//		connection:      m,
//		onMessageChan:   nil,
//	}
//
//	msg := "TEST"
//	input := formatStr(msg)
//	_ = client.send(input)
//
//	b := make([]byte, 8)
//	_, _ = m.Read(b)
//
//	if string(b) != string(input) {
//		t.Fatalf("the input and output of socket mock didn't match input: %s, output: %s", string(b), string(input))
//	}
//
//}
