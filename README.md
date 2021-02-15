# BFTWithoutSignatures
<div style="font-size: 15px">
A Golang with ZeroMQ implementation of the algorithm:
</div>
<div style="font-size: 15px">
From Consensus to Atomic Broadcast: Time-Free Byzantine-Resistant Protocols without Signatures
</div>
<div style="font-size: 13px">
    By Miguel Correia, Nuno Ferreira Neves, Paulo Verissimo
</div>

## Modules
### Binary Consensus
A binary consensus protocol performs consensus on a binary value b ∈ {0, 1}. The problem
can be formally defined in terms of three properties:
- Validity: if all correct processes propose the same value b, then any correct process
that decides, decides b.
- Agreement: no two correct processes decide differently.
- Termination: every correct process eventually decides.

### Reliable Broadcast
A reliable broadcast protocol ensures that all correct processes deliver the same messages,
and that messages broadcast by correct processes are delivered. It can be defined in terms
of the following properties:
- Validity: if a correct process broadcasts a message M, then some correct process
eventually delivers M.
- Agreement: if a correct process delivers a message M, then all correct processes
eventually deliver M.
- Integrity: every correct process p delivers at most one message M, and if sender(M) is
correct then M was previously broadcast by sender(M).

#### Multi-Valued Consensus

#### Vector Consensus

#### Atomic Broadcast

## Install Golang
If you have not already installed [Golang](https://golang.org/doc/install) follow the instructions here.

## Clone Repository
```bash
cd ~/go/src/
git clone https://github.com/v-petrou/BFTWithoutSignatures.git
```

## Run
### Manually
To install the program and generate the keys run:
```bash
go install BFTWithoutSignatures
BFTWithoutSignatures generate_keys <N>      // For key generation
```
Open <N> different terminals. In each terminal run:
```bash
BFTWithoutSignatures <ID> <N> <t> <Clients> <Scenario>
```

### Script
Adjust the script (BFTWithoutSignatures/scripts/run.sh) and run:
```bash
bash ~/go/src/BFTWithoutSignatures/scripts/run.sh
```
When you are done and want to kill the processes run:
```bash
bash ~/go/src/BFTWithoutSignatures/scripts/kill.sh
```

## Current project state
- [x] Messenger
- [x] Trusted Dealer

## TODO
- [ ] Threshold Encryption
- [ ] Common-Coin
- [ ] Binary Consensus
- [ ] Reliable Broadcast
- [ ] Multi-Valued Consensus
- [ ] Vector Consensus
- [ ] Atomic Broadcast

## References
- [From Consensus to Atomic Broadcast: Time-Free Byzantine-Resistant Protocols without Signatures](https://www.researchgate.net/publication/220459271_From_Consensus_to_Atomic_Broadcast_Time-Free_Byzantine-Resistant_Protocols_without_Signatures)
- [Asynchronous Byzantine Agreement Protocols](https://www.researchgate.net/publication/220248572_Asynchronous_Byzantine_Agreement_Protocols)
- [Signature-Free Asynchronous Byzantine Consensus](https://www.researchgate.net/publication/266659538_Signature-Free_Asynchronous_Byzantine_Consensus_with_tn3_and_On_Messages)