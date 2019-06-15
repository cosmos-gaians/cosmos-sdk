extern crate serde;
extern crate serde_json_core;
extern crate heapless;

use std::ffi::{CString, CStr};
use std::mem;
use std::os::raw::{c_char, c_void};
use serde::{Deserialize, Serialize};
use serde_json_core::de::from_slice;
use serde_json_core::ser::{to_string};

use heapless::String;
use heapless::consts::U1024;

extern "C" {
    fn read() -> *mut c_char;
    fn write(string: *mut c_char);
}

#[no_mangle]
pub extern "C" fn allocate(size: usize) -> *mut c_void {
    let mut buffer = Vec::with_capacity(size);
    let pointer = buffer.as_mut_ptr();
    mem::forget(buffer);

    pointer as *mut c_void
}

#[no_mangle]
pub extern "C" fn deallocate(pointer: *mut c_void, capacity: usize) {
    unsafe {
        let _ = Vec::from_raw_parts(pointer, 0, capacity);
    }
}

#[derive(Serialize, Deserialize)]
struct MsgCreateContract<'a> {
	sender:    &'a str,
	init_msg:   RegenInitMsg<'a>,
	init_funds: u64
}

#[derive(Serialize, Deserialize)]
struct RegenInitMsg<'a> {
    verifier:    &'a str,
	beneficiary:    &'a str,
}

#[derive(Serialize, Deserialize)]
struct RegenState<'a> {
    verifier:    &'a str,
	beneficiary:    &'a str,
    payout: u64
}

#[no_mangle]
pub extern "C" fn init(params_ptr: *mut c_char) -> *mut c_char {
    let params: MsgCreateContract = from_slice(&CStr::from_ptr(params_ptr).to_bytes()
        .to_vec()).expect("Could not parse MsgCreateContract json.");

    let state: String<U1024> = to_string(&RegenState {
        verifier: params.init_msg.verifier,
        beneficiary: params.init_msg.beneficiary,
        payout: params.init_funds
    }).unwrap();

    unsafe { 
        write(CString::new(state.as_bytes()).unwrap().into_raw());
    }

    CString::new("").unwrap().into_raw()
}

#[derive(Serialize, Deserialize)]
struct MsgSendContract<'a> {
	sender:   &'a str,
	msg:      &'a str,
	payment:  u64
}

#[derive(Serialize, Deserialize)]
struct RegenSendMsg {}

#[no_mangle]
pub extern "C" fn send(params_ptr: *mut c_char) -> *mut c_char {
    let params: MsgCreateContract = from_slice(&CStr::from_ptr(params_ptr).to_bytes()
        .to_vec()).expect("Could not parse MsgSendContract json.");

    let state: RegenState = from_slice(&CStr::from_ptr(read()).to_bytes()
        .to_vec()).expect("Could not parse RegenState json.");

    if params.sender == state.verifier {
        CString::new("Send tx goes here").unwrap().into_raw()
    } else {
        CString::new("").unwrap().into_raw()
    }

}
