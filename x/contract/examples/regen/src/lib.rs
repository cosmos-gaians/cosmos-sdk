extern crate heapless;
extern crate serde;
extern crate serde_json_core;

use serde::{Deserialize, Serialize};
use serde_json_core::de::from_slice;
use serde_json_core::ser::to_string;
use std::ffi::{CStr, CString};
use std::mem;
use std::os::raw::{c_char, c_void};

use heapless::consts::U1024;
use heapless::String;

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
    sender: &'a str,
    init_msg: RegenInitMsg<'a>,
    init_funds: u64,
}

#[derive(Serialize, Deserialize)]
struct RegenInitMsg<'a> {
    verifier: &'a str,
    beneficiary: &'a str,
}

#[derive(Serialize, Deserialize)]
struct RegenState<'a> {
    verifier: &'a str,
    beneficiary: &'a str,
    payout: u64,
}

#[no_mangle]
pub extern "C" fn init(params_ptr: *mut c_char) -> *mut c_char {
    let params: std::vec::Vec<u8>;
    unsafe {
        params = CStr::from_ptr(params_ptr).to_bytes().to_vec();
    }

    let params: MsgCreateContract =
        from_slice(&params).expect("Could not parse MsgCreateContract json.");
    let state: String<U1024> = to_string(&RegenState {
        verifier: params.init_msg.verifier,
        beneficiary: params.init_msg.beneficiary,
        payout: params.init_funds,
    })
    .unwrap();

    unsafe {
        write(CString::new(state.as_bytes()).unwrap().into_raw());
    }

    CString::new("").unwrap().into_raw()
}

#[derive(Serialize, Deserialize)]
struct MsgSendContract<'a> {
    sender: &'a str,
    msg: &'a str,
    payment: u64,
}

#[derive(Serialize, Deserialize)]
struct RegenSendMsg {}

#[derive(Serialize, Deserialize)]
enum JsonResponse { JsonError, JsonSuccess }

#[derive(Serialize, Deserialize)]
struct JsonError<'a> {
    error: &'a str,
}

#[derive(Serialize, Deserialize)]
struct JsonSuccess<'a> {
    result: &'a str,
}


#[no_mangle]
pub extern "C" fn send(params_ptr: *mut c_char) -> *mut c_char {
    let res = inner_send(params_ptr);
    match res {
        Ok(val) => to_string(&JsonResponse::Success(val)).into_raw(),
        Err(err) => to_string(&JsonResponse::Error(err)).into_raw()
    }
}

fn inner_send(params_ptr: *mut c_char) -> Result<&'str, serde_json_core::de::Error> {
    let params: std::vec::Vec<u8>;
    let state: std::vec::Vec<u8>;

    unsafe {
        params = CStr::from_ptr(params_ptr).to_bytes().to_vec();
        state = CStr::from_ptr(read()).to_bytes().to_vec();
    }

    let params: MsgSendContract = from_slice(&params)?;
    let state: RegenState = from_slice(&state)?;

    if params.sender == state.verifier {
        return Ok(CString::new("Send tx goes here !!!").unwrap())
    } else {
        return Ok(CString::new("").unwrap())
    }
}
