// +build go1.9

package tox

import (
	"fmt"
	"runtime"
	"sync"
)

/*
#include <tox/tox.h>
//#include <tox/toxav.h>

typedef struct MSICall_s MSICall;
typedef struct MSISession_s MSISession;

*/
import "C"

type userDataMSISession struct {
	ud0 map[*C.MSISession]*MSISession
	ud1 *sync.Map
	cc  bool // concurrent?
}

func newUserDataMSISession() *userDataMSISession {
	cc := true
	var ud0 map[*C.MSISession]*MSISession
	var ud1 *sync.Map

	if runtime.GOMAXPROCS(0) == 1 {
		cc = false
		ud0 = make(map[*C.MSISession]*MSISession, 0)
	} else {
		ud1 = new(sync.Map)
	}

	return &userDataMSISession{ud0: ud0, ud1: ud1, cc: cc}
}

func (this *userDataMSISession) set(ctox *C.MSISession, gtox *MSISession) {
	if this.cc {
		key := this.obj2Str(ctox)
		this.ud1.Store(key, gtox)
	} else {
		this.ud0[ctox] = gtox
	}
}

func (this *userDataMSISession) get(ctox *C.MSISession) *MSISession {
	if this.cc {
		key := this.obj2Str(ctox)
		ival, ok := this.ud1.Load(key)
		if !ok {
			return nil
		}
		return ival.(*MSISession)
	} else {
		if _, ok := this.ud0[ctox]; ok {
			return this.ud0[ctox]
		} else {
			return nil
		}
	}
}

func (this *userDataMSISession) del(ctox *C.MSISession) {
	if this.cc {
		key := this.obj2Str(ctox)
		this.ud1.Delete(key)
	} else {
		if _, ok := this.ud0[ctox]; ok {
			delete(this.ud0, ctox)
		}
	}
}

func (this *userDataMSISession) obj2Str(msi *C.MSISession) string {
	return fmt.Sprintf("%p", msi)
}

