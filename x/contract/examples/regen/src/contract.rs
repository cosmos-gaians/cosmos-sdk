use crate::{get_state, set_state, CosmosMsg, InitParams, SendAmount, SendParams};

use failure::{bail, Error};
use serde::{Deserialize, Serialize};
use serde_json::{from_slice, to_vec};

#[derive(Serialize, Deserialize)]
struct RegenInitMsg {
    verifier: Vec<u8>,
    beneficiary: Vec<u8>,
}

#[derive(Serialize, Deserialize)]
struct RegenState {
    verifier: Vec<u8>,
    beneficiary: Vec<u8>,
    payout: u64,
}

#[derive(Serialize, Deserialize)]
struct RegenSendMsg {}

pub fn init(params: InitParams) -> Result<Vec<CosmosMsg>, Error> {
    let msg: RegenInitMsg = from_slice(&params.msg)?;

    set_state(to_vec(&RegenState {
        verifier: msg.verifier,
        beneficiary: msg.beneficiary,
        payout: params.sent_funds,
    })?);

    Ok(Vec::new())
}

pub fn send(params: SendParams) -> Result<Vec<CosmosMsg>, Error> {
    let state: RegenState = from_slice(&get_state())?;

    if params.sender == state.verifier {
        Ok(vec![CosmosMsg::SendTx {
            from_address: params.contract_address,
            to_address: state.beneficiary,
            amount: vec![SendAmount {
                denom: "tree".into(),
                amount: state.payout.to_string(),
            }],
        }])
    } else {
        bail!("Unauthorized")
    }
}
