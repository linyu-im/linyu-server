package ws

func HeartbeatHandler(client *Client, req *Request) (interface{}, error) {
	return "heartbeat", nil
}
