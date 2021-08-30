package main

// // Server transport. Object which provides client transports.
// type TServerTransport interface {
// 	Listen() error
// 	Accept() (TTransport, error)
// 	Close() error

// 	// Optional method implementation. This signals to the server transport
// 	// that it should break out of any accept() or listen() that it is currently
// 	// blocked on. This method, if implemented, MUST be thread safe, as it may
// 	// be called from a different thread context than the other TServerTransport
// 	// methods.
// 	Interrupt() error
// }

// type StreamServerTransport struct {
// }

// func (a *StreamServerTransport) Listen() error {
// 	// for {
// 	// 	time.Sleep(1 * time.Second)
// 	// }
// 	return nil
// }
// func (a *StreamServerTransport) Accept() (thrift.TTransport, error) {
// 	return thrift.NewStreamTransport(os.Stdin, os.Stdout), nil
// }

// func (a *StreamServerTransport) Close() error {
// 	return nil
// }

// func (a *StreamServerTransport) Interrupt() error {
// 	return nil
// }
