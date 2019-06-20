use crate::{get_state, set_state, CosmosMsg, InitParams, SendAmount, SendParams};

use failure::{bail, Error};
use serde::{Deserialize, Serialize};
use serde_json::{from_slice, from_str, to_vec};

/*
If source delegated tokens directly to arbiter, he can send them anywhere, even to himself.
This contract holds a delegation of tokens from source account, which can be revoked
by source at any time.

Arbiter can only choose to move the tokens to a predefined destination account.

Use: Source must grant SendCapability to the contract before initializing it.
*/
#[derive(Serialize, Deserialize)]
struct DelegInitMsg {
    // source delegates tokens to the contract to allow send (can be revoked)
    source: String,
    // arbiter decided if and when to release the delegated tokens
    arbiter: String,
    // destination will receive the payout
    destination: String,
}

#[derive(Serialize, Deserialize)]
struct DelegSendMsg {
    // release defines how many of the tokens it releases 
    release: u64,
}

#[derive(Serialize, Deserialize)]
struct DelegState {
    source: String,
    arbiter: String,
    destination: String,
    stored: u64,
}


pub fn init(params: InitParams) -> Result<Vec<CosmosMsg>, Error> {
    let msg: DelegInitMsg = from_str(params.msg.get())?;

    set_state(to_vec(&DelegState {
        source: msg.source,
        arbiter: msg.arbiter,
        destination: msg.destination,
        stored: params.sent_funds,
    })?);

    Ok(Vec::new())
}

pub fn send(params: SendParams) -> Result<Vec<CosmosMsg>, Error> {
    let msg: DelegSendMsg = from_str(params.msg.get())?;
    let mut state: DelegState = from_slice(&get_state())?;

    if params.sender == state.arbiter {
        let stored = state.stored + params.sent_funds;
        state.stored = stored - msg.release;
        set_state(to_vec(&state)?);

        Ok(vec![CosmosMsg::SendTx {
            from_address: state.source,
            to_address: state.destination,
            amount: vec![SendAmount {
                denom: "earth".into(),
                amount: msg.release.to_string(),
            }],
        }])
    } else {
        bail!("Unauthorized")
    }
}
