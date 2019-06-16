extern crate heapless;
extern crate serde;
extern crate serde_json_core;
extern crate numtoa;

use serde::{Deserialize, Serialize};
use serde_json_core::de::from_slice;
use serde_json_core::ser::to_string;
use std::ffi::{CStr, CString};
use std::mem;
use std::os::raw::{c_char, c_void};

use heapless::consts::U1024;
use heapless::String;

use numtoa::NumToA;

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
    contract_address: &'a str,
    sender: &'a str,
    msg: RegenInitMsg<'a>,
    sent_funds: u64,
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

    let pres: serde_json_core::de::Result<MsgCreateContract> = from_slice(&params);
    if pres.is_err() {
        return CString::new(r#"{"error": "Could not parse MsgCreateContract json."}"#)
            .unwrap()
            .into_raw();
    }
    let params = pres.unwrap();

    let state: String<U1024> = to_string(&RegenState {
        verifier: params.msg.verifier,
        beneficiary: params.msg.beneficiary,
        payout: params.sent_funds,
    })
    .unwrap();

    unsafe {
        write(CString::new(state.as_bytes()).unwrap().into_raw());
    }

    CString::new(r#"{"msgs": []}"#).unwrap().into_raw()
}

#[derive(Serialize, Deserialize)]
struct MsgSendContract<'a> {
    contract_address: &'a str,
    sender: &'a str,
    msg: RegenSendMsg,
    sent_funds: u64,
}

#[derive(Serialize, Deserialize)]
struct RegenSendMsg {}

#[no_mangle]
pub extern "C" fn send(params_ptr: *mut c_char) -> *mut c_char {
    let params: std::vec::Vec<u8>;
    let state: std::vec::Vec<u8>;

    unsafe {
        params = CStr::from_ptr(params_ptr).to_bytes().to_vec();
        state = CStr::from_ptr(read()).to_bytes().to_vec();
    }

    let pres: serde_json_core::de::Result<MsgSendContract> = from_slice(&params);
    if pres.is_err() {
        return CString::new(r#"{"error": "Could not parse MsgSendContract json."}"#)
            .unwrap()
            .into_raw();
    }
    let params = pres.unwrap();

    let sres: serde_json_core::de::Result<RegenState> = from_slice(&state);
    if sres.is_err() {
        return CString::new(r#"{"error": "Could not parse RegenState json."}"#)
            .unwrap()
            .into_raw();
    }
    let mut state = sres.unwrap();


    if params.sender == state.verifier {
        // format how much money to produce
        let funds = state.payout + params.sent_funds;

        // TODO: clean this up...
        // I make enough space for numtoa, but need to manually trim off leading spaces
        let mut amount = b"          ".to_vec();
        funds.numtoa_str(10, &mut amount);
        let mut n = 0;
        while amount[n] == b' ' {
            n+=1;
        }
        amount = amount[n..].to_vec();

        state.payout = 0;
        let state_str: String<U1024> = to_string(&state).unwrap();
        unsafe {
            write(CString::new(state_str.as_bytes()).unwrap().into_raw());
        }
        // return CString::new(amount).unwrap().into_raw();

        let mut output = br#"{"msgs":[
            {
                "type":"cosmos-sdk/MsgSend",
                "value":{
                    "from_address":""#.to_vec();
        // params.contract
        let outb = br#"",  
                    "to_address":""#.to_vec();
        // state.beneficiary            
        let outc = br#"",  
                    "amount":[
                        {
                            "denom":"tree",
                            "amount":""#.to_vec();
        // amount
        // let amount = b"500".to_vec();
        let outd = br#""  
                        }
                    ]
                }
            }
        ]}"#.to_vec();
        output.extend(params.contract_address.bytes());
        output.extend(outb);
        output.extend(state.beneficiary.bytes());
        output.extend(outc);
        output.extend(amount);
        output.extend(outd);

        CString::new(output).unwrap().into_raw()
    } else {
        CString::new(r#"{"error": "Unauthorized"}"#)
            .unwrap()
            .into_raw()
    }
}
