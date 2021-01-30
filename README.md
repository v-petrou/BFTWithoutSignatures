# BFTWithoutSignatures
### BFTWithoutSignatures is a Golang with ZeroMQ implementation of the algorithm:
<div style="text-align: center; font-size: 15px">
From Consensus to Atomic Broadcast: Time-Free Byzantine-Resistant Protocols without Signatures
</div>
<div style="text-align: right; font-size: 13px">
    By Miguel Correia, Nuno Ferreira Neves and Paulo Verissimo
</div>

## Install Golang
1. If you have not already installed Golang follow the instructions here: https://golang.org/doc/install

## Clone the repository
1. Open Terminal.
2. Run
```bash
cd ~/go/src/
git clone https://github.com/v-petrou/BFTWithoutSignatures.git
```

## Run the program
```bash
go install BFTWithoutSignatures
BFTWithoutSignatures generate_keys <N>
BFTWithoutSignatures <ID> <N> <t> <Clients> <Scenario>
```