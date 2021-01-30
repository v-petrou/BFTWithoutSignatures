# BFTWithoutSignatures
A Golang with ZeroMQ implementation of the algorithm:
From Consensus to Atomic Broadcast: Time-Free Byzantine-Resistant Protocols without Signatures.

To clone the repository:
    1. Open Terminal.
    2. Change the current working directory to the location where you want the cloned directory.
    3. Run
        $ git clone https://github.com/v-petrou/BFTWithoutSignatures.git

To run the program do:
    $ go install BFTWithoutSignatures
    $ BFTWithoutSignatures generate_keys <N>
    $ BFTWithoutSignatures <ID> <N> <t> <Clients> <Scenario>
