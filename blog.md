Regen Network, and our compatriots from TruStory, IOV, Althea and Wallet Connect combined forces to become Team Gaians for the Berlin HackAtom.  This blog will review what we created and what it can be used for.

Out team focused on hacking towards a clear use case: user friendly participation in smart contracts about ecological health.  In order to create this usability the team had to overcome some key challenges: 
Create a smart contracting framework for the Cosmos SDK using WebAssembly
Improve key management 

And that is exactly what was accomplished.  In 36 hours, the team pulled together an incredibly powerful and flexible smart contracting module and permissions management module that interact to create a VM that is more powerful and robust than Ethereum, and built in the Cosmos SDK, making interoperability with other chains easy.

Let’s walk through the functionality that was built using the use-case we were considering - a simple three-party contract for Regen Network involving a farmer, funder, and verifier. The funder and land steward create a contract that sends money to the farmer if the approved third-party verify certifies that the land steward successfully implemented the regenerative practices. In a real-world scenario our contract would be more complex, possibly involving some decentralized verification using satellite data as well, but for now this is a minimum viable contract.

To make this scenario possible on the blockchain, we wanted to implement some of the missing technology pieces.
In particular, some sort of smart contracting framework and better key management for the non-technical users of our system.

# Making our blockchain work for non-technical users

Needing to manage keys and tokens is one of the most challenging problems of getting blockchain adopted by non-technical users.
We have heard time and time again that we cannot expect farmers and field agents to be comfortable managing private keys and filling
up a wallet with tokens just to pay for gas to verify contracts.

We added several new features to the Cosmos SDK which we believe will help solve this problem: key groups, delegated fees,
and delegated actions.

## Delegated Fees

One of the base layer things that we intuited was needed is delegated fees. Our UX team has confronted us several times
with the fact that the current design of blockchains would require users like field verifiers to have a live wallet
filled up with tokens just to send verification reports to contracts. This is simply an untenable for non-technical
users.

So we looked at where fees are being handled in the `AnteHandler` of the Cosmos SDK and made a few tweaks. In our
fork, the Cosmos `StdTx` type now includes a field for `FeeAccount`. To check whether fees have been delegated
to the account that actually signed the transaction we created a `delegation` module which allows any account
to delegate a fee allowance to any other account.

Fee allowances are defined by the `FeeAllowance` interface:

```go
// FeeAllowance defines a permission for one account to use another account's balance
// to pay fees
type FeeAllowance interface {
	// Accept checks whether this allowance allows the provided fees to be spent,
	// and optionally updates the allowance or deletes it entirely
	Accept(fee sdk.Coins, block abci.Header) (allow bool, updated FeeAllowance, delete bool)
}
```

This interface allows for quite a bit of flexibility as to how fees can be delegates. A user could specify a daily or 
montly limit for example.

For our use case, the verifier of the contract could be delegated a certain daily fee allowance by their employer and
there would be no need for these verifiers to worry about tokens at all. Only the management of the organization would
to worry about keeping company wallet filled and delegating appropriate allowances for their employees.

## Delegated Actions

So now we've solved part of the problem that our field verifier might have. This user can have a key on their phone
that never touches tokens, but when we're writing our contract who is the verifier? Do we actually know which verifier
will actually sign the contract? Does everything need to get signed by the main corporate key? Do all verifiers need
to walk around with a phone that has the same key as every other verifier so that the contract only needs to reference
one key?

Part of our generalized solution is delegated actions.

Delegated actions allow any user to delegate a permission to any action on the blockchain with any level of granularity
desired. We can do this with the `Capability` interface:

```go
type Capability interface {
	// MsgType returns the type of Msg's that this capability can accept
	MsgType() sdk.Msg
	// Accept determines whether this grant allows the provided action, and if
	// so provides an upgraded capability grant
	Accept(msg sdk.Msg, block abci.Header) (allow bool, updated Capability, delete bool)
}
```

This interface is pretty flexible and lets a Cosmos SDK developer define any type of capability desired for any
`sdk.Msg`. Because capabilities can send back an updated state we can have capabilities like a spend limit that 
spends down to zero or a daily spend limit like for fees. By default, every capability grant can also be set with
an expiration time.

## Key Groups

Another part of our generalized solution to the organizational and individual key management problem is key groups.

Key groups are basically multisig wallets which allow keys to be added or removed. Every member of a group can
be assigned a weight and the group as a whole has a decision threshold which is the number of weighted votes that
need to be acquired to take an action.

Groups can make proposals on any action that can be taken on the blockchain and send them back to the main
app router via the `delegation` module which allows groups to execute actions that have been delegated to them.

So to simplify the management of our verification organization they could have a key group at the top level for
the owners of the company. That top-level group could delegate fees and the ability to sign specific contracts
to some of their employees. Each user in turn could have a key group for themselves which has keys for each of their
devices and could even specify the key of a trusted key-recovery provider if they desired.

Hopefully these functionalities together will form a solid basis for greatly simplifying key management for individuals
and organizations and enable wide adoption of Cosmos-based Dapps by non-technical users.

# A WASM Virtual Machine on top of Cosmos 

We chose to implement our smart contracting framework using a WebAssembly (WASM) virtual machine. There is already a
good ecosystem of tooling around WASM and this would allow us to use a sophisticated language like Rust.

As our VM we originally explored [Perlin Networks Life](https://github.com/perlin-network/life) because it has built-in gas
metering code, but we ran into some issues - possibly just due to lack of examples - and switched to
[wasmer](https://github.com/wasmerio/wasmer) which was straightforward to work with.

To integrate the WASM VM with the rest of our Cosmos blockchain we came up with an interesting model where a contract
can act like any other account in the blockchain send messages back to the main message router as if the contract
were a user that had signed the transaction. The messages that a contract sends are sent as the return value of the
contract's `main` function and then run after the contract returns so there is no issue of re-entrancy.

In order to know whether a contract has permission to execute a specific action we needed another layer on top of
Cosmos's `BaseApp` router. This functionality is handled by the `delegation` module which checks if the action
can be executed directly by the contract of if this action was delegated by some other account.

We defined a relatively simple API for smart contracts, that managed to cover all our use cases. This should be
expanded as needed in the future, but we tend to minimal complexity needed to perform the task at hand.
In order to run a contract, there are three steps:

1. [Upload the smart contract wasm code](https://github.com/cosmos-gaians/cosmos-sdk/blob/hackatom/x/contract/handler.go#L59). 
   This is a large chunk (dozens of KB now, to be optimized) and is stored one
   time, allowing us to instantiate multiple instances of a contract without uploading more code.
2. [Create a contract](https://github.com/cosmos-gaians/cosmos-sdk/blob/hackatom/x/contract/handler.go#L74), 
   which is one instance of the code. This creates a new account for an instance of this code.
   It then moves any sent tokens to the contract and calls the `init` method exported by the wasm code.
3. Once there is a contract instance, you can [call the contract](https://github.com/cosmos-gaians/cosmos-sdk/blob/hackatom/x/contract/handler.go#L83) 
   any number of times, in order to execute the logic. This will call the `send` method exported by the wasm code.

Both the `MsgCreateContract` and `MsgSendContract` take a set of tokens to transfer to the contract, as well as 
raw `[]byte` with the contract-specific message. This messages is known to the client and the wasm code, but doesn't
need to be known to the sdk code, just like tendermint is agnostic to the tx bytes. For simplicity, we chose to
settle on JSON as the serialization format, and we were all quite happy with this. 

### Looking at an example contract

When we pass the information
to the contract, we create a json message containing verified information from the sdk, along with some raw 
app-dependent bytes that can be interpreted by an individual contract. We provide the verified information of 
`contract_address`, `sender`, and `sent_funds` in both 
[SendParams](https://github.com/cosmos-gaians/cosmos-sdk/blob/hackatom/x/contract/examples/regen/src/lib.rs#L17-L25) and 
[InitParams](https://github.com/cosmos-gaians/cosmos-sdk/blob/hackatom/x/contract/examples/regen/src/lib.rs#L91-L98):


```rust
#[derive(Serialize, Deserialize)]
pub struct InitParams<'a> {
    contract_address: String,
    sender: String,
    #[serde(borrow)]
    msg: &'a RawValue,
    sent_funds: u64,
}

#[derive(Serialize, Deserialize)]
pub struct SendParams<'a> {
    contract_address: String,
    sender: String,
    #[serde(borrow)]
    msg: &'a RawValue,
    sent_funds: u64,
}
```

Each contract can then define it's own format for parsing the msg, such as [RegenInitMsg](https://github.com/cosmos-gaians/cosmos-sdk/blob/2020d2d11834e1cb2fe1971fcc83201a5aac2a8b/x/contract/examples/regen/src/contract.rs#L7-L11). Note we define all these items in standard cosmos-sdk
json format, which is eg `cosmos1q....` for the addresses below, as well as the `sender` and `contract_address` passes in the params

```rust
#[derive(Serialize, Deserialize)]
struct RegenInitMsg {
    verifier: String,
    beneficiary: String,
}
```

When the contract wants to react to these incoming messages, it needs some state to make it's decisions. 
We export functions `get_state()` and `set_state()` to the wasm contract, which allow it to set and get
arbitrary data in a dedicated key in the contract substore. We also define this as a json structure,
which is initialized in the `init` call, and can be used to control execution of the `send` call.

```rust
#[derive(Serialize, Deserialize)]
struct RegenState {
    verifier: String,
    beneficiary: String,
    payout: u64,
    funder: String,
}
```

Now we see the functionality the framework exposes to a contract, we can see how easy it is to write some logic.
[init](https://github.com/cosmos-gaians/cosmos-sdk/blob/hackatom/x/contract/examples/regen/src/contract.rs#L24-L35) 
will just store the information on who can execute the contract. Note that it uses a mix of sdk-verified
information `params.sent_funds`, with user-specified content `msg.verifier`:

```rust
pub fn init(params: InitParams) -> Result<Vec<CosmosMsg>, Error> {
    let msg: RegenInitMsg = from_str(params.msg.get())?;

    set_state(to_vec(&RegenState {
        verifier: msg.verifier,
        beneficiary: msg.beneficiary,
        payout: params.sent_funds,
        funder: params.sender
    })?);

    Ok(Vec::new())
}
```

[send](https://github.com/cosmos-gaians/cosmos-sdk/blob/hackatom/x/contract/examples/regen/src/contract.rs#L37-L55) can now
compare the stored state with the sdk-verified `params.sender` to control whether to release the funds or not.
Note that the return value of both `init` and `send` is `Result<Vec<CosmosMsg>, Error>`, which means it can
return an Error code, or a (possibly empty) list of cosmos messages to dispatch on success. In send, we use both
of these features:

```rust
pub fn send(params: SendParams) -> Result<Vec<CosmosMsg>, Error> {
    let mut state: RegenState = from_slice(&get_state())?;
    let funds = state.payout + params.sent_funds;
    state.payout = 0;
    set_state(to_vec(&state)?);

    if params.sender == state.verifier {
        Ok(vec![CosmosMsg::SendTx {
            from_address: params.contract_address,
            to_address: state.beneficiary,
            amount: vec![SendAmount {
                denom: "earth".into(),
                amount: funds.to_string(),
            }],
        }])
    } else {
        bail!("Unauthorized")
    }
}
```

If the caller of `MsgSendContract` is the same as the verifier set in `MsgCreateContract`, then we will dispatch
a message moving the funds inside the contract to the beneficiary specified in `init`. If it is called by any other
user, then it will return an unauthorized error. (The observant reader may notice we hardcoded the token denomination
to be "earth", which we did to simplify parsing... with a bit of polish outside of Hackatom, we can use proper sdk Coin types).

Of course, this example is quite simple, but you can quickly see that a large variety of custom escrows, authorization, and 
key management solutions can be built in this framework, with only a few lines of rust. In order to make that rust easy to work
with, there was quite some work wrangling wasm and cgo bindings, and figuring how to pass arbitrary json through the wasm
function call interface, which is defined in go-speak as `f(args ...int32) int32`. Details of that and how to interface
go/rust/wasm through low-level c-bindings is quite interesting for anyone working to integrate wasm into their
go applications, but a topic for another blog post.

