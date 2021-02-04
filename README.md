# BFTWithoutSignatures
BFTWithoutSignatures is a Golang with ZeroMQ implementation of the algorithm:
<div style="text-align: center; font-size: 15px">
From Consensus to Atomic Broadcast: Time-Free Byzantine-Resistant Protocols without Signatures
</div>

<div style="text-align: right; font-size: 13px">
    By Miguel Correia, Nuno Ferreira Neves and Paulo Verissimo
</div>

## Install Golang
If you have not already installed Golang follow the instructions here: https://golang.org/doc/install

## Clone Repository
```bash
cd ~/go/src/
git clone https://github.com/v-petrou/BFTWithoutSignatures.git
```

## Run Program
```bash
go install BFTWithoutSignatures
BFTWithoutSignatures generate_keys <N>
BFTWithoutSignatures <ID> <N> <t> <Clients> <Scenario>
```

## Readings
#### Algorithm
https://www.semanticscholar.org/paper/From-Consensus-to-Atomic-Broadcast%3A-Time-Free-Correia-Neves/59d140a4b70cda42bb4c2bd5eb1908d7c1ba3a87

#### Reliable Broadcast Module
https://www.semanticscholar.org/paper/Asynchronous-Byzantine-Agreement-Protocols-Bracha/510e071390c7fbd166bee9359e79d8a68a273f66

#### Binary Consensus Module
https://www.semanticscholar.org/paper/Signature-free-asynchronous-byzantine-consensus-t-%3C-Most%C3%A9faoui-Hamouma/096c9bd145e9d11e055f65739499ccdac22fb43a
