## ethvanity

`ethvanity` is a fast and simple command line tool to generate so-called "vanity" Ethereum
addresses.

These are addresses that have a certain pattern - typically, they are specific suffixes or
prefixes to the hex string that represents the address.

For example, say you want your address to start with 3 zeros. Such an address could be
`0x00009cd4D5c81D51c53b61344C31266F8A1C2363`. Of course, there are many such addresses.
The longer your prefix, the longer it'll take to generate such a vanity address, since
the best "algorithm" is to simply generate a random private key, from that get the public
key, and from that get the address - which is the last 20 bytes of the keccak-256 hash of
the public key.

## How To Use

`ethvanity` is simple to use - simply run:

```bash
go install github.com/butcher-of-blaviken/ethvanity
```

Or clone the repository locally on your machine and run:

```bash
# optional, could use go run main.go as well
go build
# omit the './' if you ran 'go install'
./ethvanity -cmd generate -desired-pattern 000 -pattern-position prefix
```

The output is an address printout with the desired pattern and the associated private key in
hexadecimal.

## Important To Note

* This tool does not require the internet. Indeed, it is recommended
you cut off your computer from the internet when executing this tool, just in case :-).

* The longer your prefix/suffix pattern is, the longer it'll take to find your desired address.
A simple rule of thumb is that each extra character will multiply the number of iterations by
a factor of 16, since there are 16 different possibilities for each hexadecimal character.
