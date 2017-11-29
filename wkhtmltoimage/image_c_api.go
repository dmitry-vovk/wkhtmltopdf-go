package wkhtmltoimage

//#cgo CFLAGS: -I/usr/local/include
//#cgo LDFLAGS: -L/usr/local/lib -lwkhtmltox -Wall -ansi -pedantic -ggdb
//#include <stdbool.h>
//#include <stdio.h>
//#include <string.h>
//#include <stdlib.h>
//#include <wkhtmltox/image.h>
//extern void finishedCallback(void*, const int);
//extern void progressChangeCallback(void*, const int);
//extern void errorCallback(void*, char *msg);
//extern void warningCallback(void*, char *msg);
//extern void phaseChangeCallback(void*);
//static void setup_callbacks(wkhtmltoimage_converter * c) {
//  wkhtmltoimage_set_finished_callback(c, (wkhtmltoimage_int_callback) finishedCallback);
//  wkhtmltoimage_set_progress_changed_callback(c, (wkhtmltoimage_int_callback) progressChangeCallback);
//  wkhtmltoimage_set_error_callback(c, (wkhtmltoimage_str_callback) errorCallback);
//  wkhtmltoimage_set_warning_callback(c, (wkhtmltoimage_str_callback) warningCallback);
//  wkhtmltoimage_set_phase_changed_callback(c, (wkhtmltoimage_void_callback) phaseChangeCallback);
//}
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

type GlobalSettings struct {
	s *C.wkhtmltoimage_global_settings
}

type Converter struct {
	c               *C.wkhtmltoimage_converter
	Finished        func(*Converter, int)
	ProgressChanged func(*Converter, int)
	Error           func(*Converter, string)
	Warning         func(*Converter, string)
	Phase           func(*Converter)
}

var converterMap = make(map[unsafe.Pointer]*Converter)

func init() {
	C.wkhtmltoimage_init(C.false)
}

func NewGlobalSettings() *GlobalSettings {
	return &GlobalSettings{s: C.wkhtmltoimage_create_global_settings()}
}

func (g *GlobalSettings) Set(name, value string) *GlobalSettings {
	cName := C.CString(name)
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cValue))
	C.wkhtmltoimage_set_global_setting(g.s, cName, cValue)
	return g
}

func (g *GlobalSettings) NewConverter() *Converter {
	c := &Converter{c: C.wkhtmltoimage_create_converter(g.s, nil)}
	C.setup_callbacks(c.c)
	return c
}

//export finishedCallback
func finishedCallback(c unsafe.Pointer, s C.int) {
	cm := converterMap[c]
	if cm.Finished != nil {
		cm.Finished(cm, int(s))
	}
}

//export progressChangeCallback
func progressChangeCallback(c unsafe.Pointer, p C.int) {
	cm := converterMap[c]
	if cm.ProgressChanged != nil {
		cm.ProgressChanged(cm, int(p))
	}
}

//export errorCallback
func errorCallback(c unsafe.Pointer, msg *C.char) {
	cm := converterMap[c]
	if cm.Error != nil {
		cm.Error(cm, C.GoString(msg))
	}
}

//export warningCallback
func warningCallback(c unsafe.Pointer, msg *C.char) {
	cm := converterMap[c]
	if cm.Warning != nil {
		cm.Warning(cm, C.GoString(msg))
	}
}

//export phaseChangeCallback
func phaseChangeCallback(c unsafe.Pointer) {
	cm := converterMap[c]
	if cm.Phase != nil {
		cm.Phase(cm)
	}
}

func (c *Converter) Convert() error {
	// To route callbacks right, we need to save a reference
	// to the converter object, base on the pointer.
	converterMap[unsafe.Pointer(c.c)] = c
	status := C.wkhtmltoimage_convert(c.c)
	delete(converterMap, unsafe.Pointer(c.c))
	if status != C.int(1) {
		return fmt.Errorf("convert failed (%d)", status)
	}
	return nil
}

func (c *Converter) ErrorCode() int {
	return int(C.wkhtmltoimage_http_error_code(c.c))
}

func (c *Converter) Destroy() {
	C.wkhtmltoimage_destroy_converter(c.c)
}

func (c *Converter) Payload() ([]byte, int) {
	var (
		payloadPt *C.uchar
		payload   []byte
	)
	size := int(C.wkhtmltoimage_get_output(c.c, &payloadPt))
	header := (*reflect.SliceHeader)(unsafe.Pointer(&payload))
	header.Len = size
	header.Cap = size
	header.Data = uintptr(unsafe.Pointer(payloadPt))
	return payload, size
}

func (c *Converter) PhaseDescription(i int) string {
	return C.GoString(C.wkhtmltoimage_phase_description(c.c, C.int(i)))
}

func (c *Converter) CurrentPhase() int {
	return int(C.wkhtmltoimage_current_phase(c.c))
}
