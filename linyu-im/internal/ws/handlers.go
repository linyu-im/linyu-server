package ws

func HeartbeatHandler(client *Client, req *Request) (interface{}, error) {
	return "heartbeat", nil
}

func AckHandler(client *Client, req *Request) (interface{}, error) {
	Manager.DeleteRetryMessage(req.Device, req.SeqId)
	return "success", nil
}
