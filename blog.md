Regen Network, and our compatriots from TruStory, IOV, Althea and Wallet Connect combined forces to become Team Gaians for the Berlin HackAtom.  This blog will review what we created and what it can be used for.

Our team focused on hacking towards a clear use case: user friendly participation in smart contracts about ecological health. In order to create usable functionality for this the team had to overcome some key challenges: create a smart contracting framework for the Cosmos SDK and improved key management.

And that is exactly what was accomplished. In 36 hours, the team pulled together an incredibly powerful and flexible smart contracting module using WebAssembly (Wasm) and key management modules that interact to create platform that is more powerful and robust than Ethereum, and built in the Cosmos SDK, making interoperability with other chains easy.

Let’s walk through the functionality that was built using the use-case we considered - a simple three-party contract for Regen Network involving a farmer, funder, and verifier. The funder and land steward create a contract that sends money to the farmer if the approved third-party verifier certifies that the land steward successfully implemented regenerative practices. In a real-world scenario our contract would be a bit more complex, possibly involving decentralized verification using satellite data, but for now this is a minimum viable contract.

# Making our blockchain work for non-technical users

Needing to manage keys and tokens is one of biggest barriers preventing non-technical users from adopting blockchain based applications.

We have heard time and time again that we cannot expect farmers and agricultural field agents to be comfortable managing private keys and filling up a wallet with tokens just to pay for gas to verify contracts.

We added several new features to the Cosmos SDK which we believe help solve this problem: key groups, delegated fees, and delegated actions.

## Delegated Fees
One of the most basic usability enhancements we wanted to add was delegated fees.

Our UX team has confronted us several times with the fact that the current design of blockchains would require users like agricultural field agents to have a live wallet filled up with tokens, likely on their smartphone, just to send verification reports to contracts. This is simply untenable for non-technical users.

So we looked at where fees are being handled in the `AnteHandler` of the Cosmos SDK and made a few changes. In our fork, the Cosmos `StdTx` type now includes a field for `FeeAccount`. To check whether fees have been delegated to the account that actually signed the transaction we created a `delegation` module which allows any account to delegate a fee allowance to any other account. Fee allowances can then be defined by implementing the `FeeAllowance` interface:

```go
// FeeAllowance defines a permission for one account to use another account's balance
// to pay fees
type FeeAllowance interface {
	// Accept checks whether this allowance allows the provided fees to be spent,
	// and optionally updates the allowance or deletes it entirely
	Accept(fee sdk.Coins, block abci.Header) (allow bool, updated FeeAllowance, delete bool)
}
```

This interface allows for quite a bit of flexibility as to how fees can be delegates. A user could specify a daily or monthly limit for example.

For our use case, the verifier of the contract could be delegated a certain daily fee allowance by their employer and there would be no need for these field agents to worry about tokens at all. Only the management of the organization would need to worry about keeping the company wallet filled and delegating appropriate allowances to their employees.

## Delegated Actions

This solves part of the problem that our field verifier might have. This user can have a key on their phone that never touches tokens, but when we're writing our contract who is the verifier that actually signs the claim? Do we actually know which particular field agent will actually sign the claim? Does everything need to get signed by the main corporate key? Do all field agents need to walk around with a phone that has the master key that's referenced in our contract?

A better solution would be to allow an organization to delegate some permissions - such as signing a verification claim - to their employees. To solve this we came up with a generalized solution for delegated actions.

Delegated actions allow any user to delegate a permission to perform any action on their behalf to another user. We can do this with the `Capability` interface which allows for any level of permission granularity desired:

```go
type Capability interface {
	// MsgType returns the type of Msg's that this capability can accept
	MsgType() sdk.Msg
	// Accept determines whether this grant allows the provided action, and if
	// so provides an upgraded capability grant
	Accept(msg sdk.Msg, block abci.Header) (allow bool, updated Capability, delete bool)
}
```

This interface is pretty flexible and lets a Cosmos SDK developer define a capability for any
`sdk.Msg`. Because capabilities can send back an updated state we can write capabilities that set a total spend limit and then are deleted when that amount of coins has been set. By default, every capability grant can also be set with an expiration time.

## Key Groups

Another part of our generalized solution to the organizational and individual key management problem is key groups.

Key groups are basically multi-signature wallets which allow keys to be added or removed. Every member of a group can be assigned a weight and the group as a whole has a decision threshold which is the number of weighted votes that need to be acquired to take an action. Key groups can include other groups or even contracts as their group members.

Groups can make proposals on any action that can be taken on the blockchain and send them back to the main app router via the `delegation` module which allows groups to execute actions that have been delegated to them.

So to simplify the management of our verification organization, they could have a top-level key group at for the owners of the company. That top-level group would delegate fees and the ability to sign specific contracts to some of their employees. Each user in turn could have a key group which includes the keys for all of their devices. Key groups could even include a trusted third-party provider that signs using normal username/password methods to create a multi-factor authentication scenario that requires a device plus the third-party to sign.

We believe these functionalities together form a solid basis for greatly simplifying key management for individuals and organizations and we hope they will enable wide adoption of Cosmos-based Dapps by non-technical users.

# A WASM Virtual Machine on top of Cosmos 

The other part of our user scenario involves a simple multi-party smart contract. Our hackathon team was pretty ambitious so while one part of the team was implementing the key management solutions described above, another part was integrating a WebAssembly (Wasm) virtual machine into the Cosmos SDK.

We chose to implement our smart contracting framework using a Wasm virtual machine because there is already a good ecosystem of tooling around Wasm and this would allow us to use a sophisticated language like Rust.

Wwe originally explored [Perlin Networks Life](https://github.com/perlin-network/life) as our Wasm VM because it has built-in gas metering code, but we ran into some issues - possibly just due to lack of examples - and switched to [wasmer](https://github.com/wasmerio/wasmer) which was straightforward to work with.

To integrate the Wasm VM with the rest of our Cosmos blockchain we came up with an interesting model where a contract can act like any other account in the blockchain send messages back to the main message router as if the contract were a user that had signed a transaction. The messages that a contract sends are sent as the return value of the contract's `main` function and then run after the contract returns so there is no issue of re-entrancy.

In order to know whether a contract has permission to execute a specific action we needed another layer on top of Cosmos's `BaseApp` router. This functionality is handled by the `delegation` module described above which checks if the action can be executed directly by the contract of if this action was delegated by some other account.

We defined a relatively simple API for smart contracts, that covered our famer/funder/verifier use case. This could be expanded as needed in the future, but was kept to the minimal amount of complexity needed to perform the task at hand.

In order to run a contract, there are three steps:

1. [Upload the smart contract wasm code](https://github.com/cosmos-gaians/cosmos-sdk/blob/hackatom/x/contract/handler.go#L59). This is a large chunk (dozens of KB for now, but it can be optimized) and is stored one time, allowing us to instantiate multiple contract instances for a single piece of code without needing to upload that code multiple times.
2. [Create a contract](https://github.com/cosmos-gaians/cosmos-sdk/blob/hackatom/x/contract/handler.go#L74), which is one instance of the code. This creates a new account for this instance. It then moves any sent tokens to the contract and calls the `init` method exported by the wasm code.
3. Once there is a contract instance, you can [call the contract](https://github.com/cosmos-gaians/cosmos-sdk/blob/hackatom/x/contract/handler.go#L83) any number of times, in order to execute the contract logic implemented in its `send` method exported from the Wasm code.

Both `MsgCreateContract` and `MsgSendContract` take a set of tokens to transfer to the contract, as well as a raw `[]byte` with a contract-specific message. This message is known to the client and the Wasm code, but doesn't need to be known to the SDK code, just like Tendermint is agnostic to tx bytes. For simplicity, we chose to settle on JSON as the serialization format. We were all quite happy with this, but other methods could be used in the future.

### Looking at an example contract

When we pass information to the contract, we create a JSON message containing the `contract_address`, `sender`, and `sent_funds` from the SDK, along with the raw app-dependent bytes sent to the contract. This information is passed to both [SendParams](https://github.com/cosmos-gaians/cosmos-sdk/blob/hackatom/x/contract/examples/regen/src/lib.rs#L17-L25) and [InitParams](https://github.com/cosmos-gaians/cosmos-sdk/blob/hackatom/x/contract/examples/regen/src/lib.rs#L91-L98):


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

Each contract can then define it's own format for parsing the custom message, such as [RegenInitMsg](https://github.com/cosmos-gaians/cosmos-sdk/blob/2020d2d11834e1cb2fe1971fcc83201a5aac2a8b/x/contract/examples/regen/src/contract.rs#L7-L11). Note that we pass `sender` and `contract_address` as standard Cosmos bech32 address strings, eg. `cosmos1q....`.

```rust
#[derive(Serialize, Deserialize)]
struct RegenInitMsg {
    verifier: String,
    beneficiary: String,
}
```

When the contract wants to react to these incoming messages, it needs some state to make its decisions. We export functions `get_state()` and `set_state()` to the wasm contract, which allow it to set and get arbitrary data in a dedicated key in the contract substore. We could extend this API in the future, but it sufficed for now. We also define the contract's state as a JSON structure, which is initialized in the `init` call, and can be used to control execution of the `send` call.

```rust
#[derive(Serialize, Deserialize)]
struct RegenState {
    verifier: String,
    beneficiary: String,
    payout: u64,
    funder: String,
}
```

Now that we see the functionality the framework exposes to a contract, we can see how easy it is to write some logic. [init](https://github.com/cosmos-gaians/cosmos-sdk/blob/hackatom/x/contract/examples/regen/src/contract.rs#L24-L35) just stores the information on who can execute the contract. Note that it uses a mix of information passed by the SDK such as `params.sent_funds` and user-specified content such as `msg.verifier` and `msg.beneficiary`:

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

[send](https://github.com/cosmos-gaians/cosmos-sdk/blob/hackatom/x/contract/examples/regen/src/contract.rs#L37-L55) can now compare the stored state with the SDK-verified `params.sender` to control whether to release the funds or not. Note that the return value of both `init` and `send` is `Result<Vec<CosmosMsg>, Error>`, which means it can return an Error code, or a (possibly empty) list of cosmos messages to dispatch on success. In `send`, we use both of these features:

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

If the caller of `MsgSendContract` is the same as the verifier set in `MsgCreateContract`, then we will dispatch a message moving the funds inside the contract to the beneficiary specified in `init`. If it is called by any other user, then it will return an unauthorized error. (The observant reader may notice we hardcoded the token denomination to be `earth`, which we did to simplify parsing... with a bit of polish outside of the Hackatom, we can use proper `sdk.Coins` types).

Of course, this example is quite simple, but you can quickly see that a large variety of custom escrows, authorization, and key management solutions can be built with this framework, using only a few lines of Rust. Making that Rust easy to work with involved quite a bit of work wrangling Wasm and cgo bindings, and figuring how to pass arbitrary json through the wasm function call interface, which is defined in go-speak as `f(args ...int32) int32`. Details of that and how to interface Go, Rust, and Wasm through low-level C-bindings is quite interesting for anyone working to integrate wasm into their go applications, but a topic for another blog post.

# Future Directions

Since building these features at the Hackathon there has been interest in the Cosmos community beyond just Regen Network in bringing these features to production level. As a follow-up to some discussions online, we have created two community groups to work on key management and WASM with the goal of producing PR's to be merged into the Cosmos SDK:

- The Cosmos Community Group on Key Management: https://github.com/cosmos-cg-key-management
- CosmWasm (Cosmos Community Group on WebAssembly (WASM) integration): https://github.com/cosmwasm

Up until now most development on Cosmos SDK has happened directly from the Cosmos team, albeit with quite a large number of contributions from different community groups. The idea of a community group or working group is to create a space where members of different organizations can collaborate on common problems without the ownership of the process falling within either organization. For these community groups, we have setup a simple structure involving a dedicated GitHub organization and a telegram channel for live discussion. Possibly this model can be replicated for other similar community efforts.

We welcome community members interested in getting these features to production-level to join these groups and collaborate, and we hope these contributions bring long-lasting benefits to the community.