package threads

type ErrorMessage struct {
	process   string
	e         error
	receivers []CallBackReceiver[string]
}

type Message[T any] struct {
	process   string
	t         T
	hasResult bool
	receivers []CallBackReceiver[T]
}

func CommonError(process string, e error) ErrorMessage {
	return ErrorMessage{process: process, e: e, receivers: nil}
}

func CommonMessage(process string, result string) Message[string] {
	return Message[string]{process: process, t: result, receivers: nil, hasResult: true}
}

func EmptyMessage(process string) Message[string] {
	return Message[string]{process: process, t: "", receivers: nil, hasResult: false}
}
