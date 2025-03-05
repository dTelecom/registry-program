https://www.anchor-lang.com/docs/installation

## Build

1. Install [Anchor](https://github.com/coral-xyz/anchor).
1. Navigate to the project root and run `anchor build`. This compiles the project and generates a keypair for the program based on your personal wallet.
1. The following steps only need to be performed on the first build (to replace the public key with one corresponding to your personal wallet).
    1. Run `anchor keys list` to get the new public key for the program.
    1. Replace all occurences of the old public key in the project with the new public key.
    1. Re-compile with `anchor build` (to account for the replacement of the public key in the source code).
1. Run `anchor deploy` to deploy the compiled program on the blockchain.
1. Run the tests with `anchor run test`.

# Build
anchor build

# Test
anchor test

# Deploy to localnet
anchor deploy
