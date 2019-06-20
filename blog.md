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

TODO: more details about the implementation
