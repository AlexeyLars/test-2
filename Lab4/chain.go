package main

import "fmt"

type RequestType int

const (
	TypeA RequestType = iota
	TypeB
	TypeC
)

type Request struct {
	requestType RequestType
}

func (r *Request) GetType() RequestType {
	return r.requestType
}

type Handler interface {
	HandleRequest(request *Request)
	SetNextHandler(handler Handler)
}

type BaseHandler struct {
	nextHandler Handler
}

func (h *BaseHandler) SetNextHandler(handler Handler) {
	h.nextHandler = handler
}

type ConcreteHandlerA struct {
	BaseHandler
}

func (h *ConcreteHandlerA) HandleRequest(request *Request) {
	if request.GetType() == TypeA {
		fmt.Println("ConcreteHandlerA handled the request")
	} else if h.nextHandler != nil {
		h.nextHandler.HandleRequest(request)
	} else {
		fmt.Println("Request could not be handled")
	}
}

type ConcreteHandlerB struct {
	BaseHandler
}

func (h *ConcreteHandlerB) HandleRequest(request *Request) {
	if request.GetType() == TypeB {
		fmt.Println("ConcreteHandlerB handled the request")
	} else if h.nextHandler != nil {
		h.nextHandler.HandleRequest(request)
	} else {
		fmt.Println("Request could not be handled")
	}
}

type ConcreteHandlerC struct {
	BaseHandler
}

func (h *ConcreteHandlerC) HandleRequest(request *Request) {
	if request.GetType() == TypeC {
		fmt.Println("ConcreteHandlerC handled the request")
	} else if h.nextHandler != nil {
		h.nextHandler.HandleRequest(request)
	} else {
		fmt.Println("Request could not be handled")
	}
}

func main() {
	fmt.Println("=== Chain of Responsibility Pattern ===")

	handlerA := &ConcreteHandlerA{}
	handlerB := &ConcreteHandlerB{}
	handlerC := &ConcreteHandlerC{}

	handlerA.SetNextHandler(handlerB)
	handlerB.SetNextHandler(handlerC)

	requestA := &Request{requestType: TypeA}
	requestB := &Request{requestType: TypeB}
	requestC := &Request{requestType: TypeC}

	handlerA.HandleRequest(requestA)
	handlerA.HandleRequest(requestB)
	handlerA.HandleRequest(requestC)
	fmt.Println()
}
