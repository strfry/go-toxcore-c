package tox

/*
#include <stdlib.h>
#include <string.h>
#include <tox/tox.h>
#include <tox/toxav.h>


 typedef enum MSIError {
    MSI_E_NONE,
    MSI_E_INVALID_MESSAGE,
    MSI_E_INVALID_PARAM,
    MSI_E_INVALID_STATE,
    MSI_E_STRAY_MESSAGE,
    MSI_E_SYSTEM,
    MSI_E_HANDLE,
    MSI_E_UNDISCLOSED, 
} MSIError;

typedef enum MSICapabilities {
    MSI_CAP_S_AUDIO = 4,  // sending audio
    MSI_CAP_S_VIDEO = 8,  // sending video
    MSI_CAP_R_AUDIO = 16, // receiving audio
    MSI_CAP_R_VIDEO = 32, // receiving video
} MSICapabilities;

typedef enum MSICallState {
    MSI_CALL_INACTIVE, // Default
    MSI_CALL_ACTIVE,
    MSI_CALL_REQUESTING, // when sending call invite
    MSI_CALL_REQUESTED, // when getting call invite
} MSICallState;

typedef enum MSICallbackID {
    MSI_ON_INVITE, // Incoming call
    MSI_ON_START, // Call (RTP transmission) started
    MSI_ON_END, // Call that was active ended
    MSI_ON_ERROR, // On protocol error
    MSI_ON_PEERTIMEOUT, // Peer timed out; stop the call
    MSI_ON_CAPABILITIES, // Peer requested capabilities change
} MSICallbackID;

// **
// * The call struct. Please do not modify outside msi.c
// * /
typedef struct MSICall_s MSICall;
typedef struct MSISession_s MSISession;
typedef struct Messenger Messenger;

typedef int msi_action_cb(void *av, MSICall *call);

MSISession *msi_new(Messenger *m);
void msi_register_callback(MSISession *session, msi_action_cb *callback, MSICallbackID id);

int msi_hangup(MSICall *call);
int msi_answer(MSICall *call, uint8_t capabilities);
int msi_change_capabilities(MSICall *call, uint8_t capabilities);

void callbackMSIActionWrapperForC(void *av, MSICall* call);

*/
import "C"
import (
	//"encoding/hex"
	//"errors"
	"unsafe"
)

type cb_msi_action_ftype func(av interface{}, this *MSICall)

type MSICall struct {
	call *C.MSICall
	session *MSISession
}

func (this *MSICall) Hangup() (error) {
	var cerr = C.msi_hangup(this.call)
	if cerr != 0 {
		return toxerr(cerr)
	}
	return nil
}

func (this *MSICall) Answer(capabilities uint8) (error) {
        var cerr = C.msi_answer(this.call, C.uchar(capabilities))
        if cerr != 0 {
                return toxerr(cerr)
        }
        return nil
}


type MSISession struct {
	tox *Tox
	msi *C.MSISession

	// session datas

	// callbacks
	cb_msi_action                          cb_msi_action_ftype
	cb_msi_action_user_data                interface{}
}


var cbMSISessions = newUserDataMSISession()

//export callbackMSIActionWrapperForC
func callbackMSIActionWrapperForC(av unsafe.Pointer, call *C.struct_MSICall_s) {
	// same hack
        var callptr0 = *(*uintptr)(unsafe.Pointer(call))
        var csession = (*C.MSISession)(unsafe.Pointer(callptr0))

        var this = cbMSISessions.get(csession)

	var msicall = MSICall{}
	msicall.call = call
        msicall.session = this

	if this.cb_msi_action != nil {
		this.cb_msi_action(av, &msicall)
	}
}

func (this *MSISession) RegisterCallback (id C.MSICallbackID, cbfn cb_msi_action_ftype) (error) {
	this.cb_msi_action = cbfn

	var _cbfn = (*C.msi_action_cb)(C.callbackMSIActionWrapperForC)
	C.msi_register_callback(this.msi, _cbfn, id)
	return nil
}

func NewMSISession(tox *Tox) (*MSISession, error) {
	if tox == nil {
		return nil, toxerr("tox can not nil")
	}

	msi := new(MSISession)
	msi.tox = tox

	// NOTE(strfry): The same ugly hack as in toxav_new...
	var toxptr0 = *(*uintptr)(unsafe.Pointer(tox.toxcore))
	var msgptr = (*C.Messenger)(unsafe.Pointer(toxptr0))

	var messenger *C.Messenger = msgptr // *phew*
	msi.msi = C.msi_new(messenger)
	
	cbMSISessions.set(msi.msi, msi)
	return msi, nil
}

func UseUnsafePointer() (unsafe.Pointer) {
	return unsafe.Pointer(nil)
}
