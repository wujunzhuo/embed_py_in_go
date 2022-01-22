package main

import (
	"log"

	py "github.com/go-python/cpy3"
)

func RunPyHandler(data []byte) (byte, []byte) {
	path, err := py.Py_GetPath()
	if err != nil {
		log.Printf("Error: py.Py_GetPath")
		return 0, nil
	}

	err = py.Py_SetPath(".:" + path)
	if err != nil {
		log.Printf("Error: py.Py_SetPath")
		return 0, nil
	}

	py.Py_Initialize()
	if !py.Py_IsInitialized() {
		log.Printf("Error: py.Py_IsInitialized")
		return 0, nil
	}
	defer py.Py_Finalize()

	module := py.PyImport_ImportModule("app")
	if module == nil {
		log.Printf("Error: py.PyImport_ImportModule")
		return 0, nil
	}
	defer module.DecRef()

	handler := module.GetAttrString("handler")
	if handler == nil {
		log.Printf("Error: module.GetAttrString")
		return 0, nil
	}
	defer handler.DecRef()

	req := py.PyBytes_FromString(string(data))
	if req == nil {
		log.Printf("py.PyBytes_FromString")
		return 0, nil
	}
	defer req.DecRef()

	reqTuple := py.PyTuple_New(1)
	if reqTuple == nil {
		log.Printf("py.PyTuple_New")
		return 0, nil
	}
	defer reqTuple.DecRef()

	if py.PyTuple_SetItem(reqTuple, 0, req) != 0 {
		log.Printf("py.PyTuple_SetItem")
		return 0, nil
	}

	resTuple := handler.CallFunctionObjArgs(req)
	if resTuple == nil {
		log.Printf("handler.CallFunctionObjArgs")
		return 0, nil
	}
	defer resTuple.DecRef()

	if !py.PyTuple_Check(resTuple) {
		log.Printf("py.PyTuple_Check")
		return 0, nil
	}

	if py.PyTuple_Size(resTuple) != 2 {
		log.Printf("py.PyTuple_Size")
		return 0, nil
	}

	tag := py.PyLong_AsLong(py.PyTuple_GetItem(resTuple, 0))
	res := py.PyBytes_AsString(py.PyTuple_GetItem(resTuple, 1))

	return byte(tag), []byte(res)
}

func main() {
	req := []byte{0x16, 0x23, 0x31}
	tag, res := RunPyHandler(req)
	log.Printf("tag: %# x", tag)
	log.Printf("res: %# x", res)
}
