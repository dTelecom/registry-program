# Solana Registry Client

A Go-based CLI client for interacting with the Solana Registry smart contract. The registry supports two types of entries:
- **Client entries**: For client accounts with expiration time and request limits
- **Node entries**: For node accounts with associated IP addresses

## Prerequisites

- Go 1.16 or later
- Solana CLI 1.14.0 or later
- Access to a Solana network (local, devnet, testnet, or mainnet)
- Private key with sufficient SOL for transactions
- Running Solana validator with RPC and WebSocket endpoints

## Compatibility

This client is compatible with:
- Solana-go v1.8.4
- Solana validator 1.14.0 or later
- Local validator: requires both RPC (default: 8899) and WebSocket (default: 8900) ports to be available

## Installation

1. Clone the repository
2. Install dependencies:
```bash
go mod tidy
```

## Configuration

Create a `.env` file in the root directory with the following contents (you can copy from `.env.example`):

```env
# Solana RPC URL (local/dev/test/mainnet)
SOLANA_RPC_URL=http://127.0.0.1:8899

# Solana WebSocket URL for PubSub
SOLANA_WS_URL=ws://127.0.0.1:8900

# Your private key in base58 format
WALLET_PRIVATE_KEY=your_private_key_here

# Registry program ID
PROGRAM_ID=E2FcHsC9STeB6FEtxBKGAwMTX7cbfYMyjSHKs4QbBAmh
```

### Network Configuration

#### Local Validator
```env
SOLANA_RPC_URL=http://127.0.0.1:8899
SOLANA_WS_URL=ws://127.0.0.1:8900
```

#### Devnet
```env
SOLANA_RPC_URL=https://api.devnet.solana.com
SOLANA_WS_URL=wss://api.devnet.solana.com
```

#### Testnet
```env
SOLANA_RPC_URL=https://api.testnet.solana.com
SOLANA_WS_URL=wss://api.testnet.solana.com
```

#### Mainnet-Beta
```env
SOLANA_RPC_URL=https://api.mainnet-beta.solana.com
SOLANA_WS_URL=wss://api.mainnet-beta.solana.com
```

## Building

```bash
go build -o registry-client
```

## Usage

### Registry Management

#### Create a new registry:
```bash
./registry-client create <registry_name>
```

### Client Operations

#### Add a client to the registry:
```bash
./registry-client add-client <registry_name> <account_to_add> <valid_days> <limit>
```
- `registry_name`: Name of the registry
- `account_to_add`: The public key of the client account to add
- `valid_days`: Number of days the client registration should remain valid
- `limit`: Request limit for the client

#### Get client information:
```bash
./registry-client get-client <registry_name> <account_to_check>
```
Displays:
- Parent registry
- Registered account
- Valid until timestamp
- Request limit

#### Delete a client:
```bash
./registry-client delete-client <registry_name> <account_to_delete>
```

### Node Operations

#### Add a node to the registry:
```bash
./registry-client add-node <registry_name> <account_to_add> <domain>
```
- `registry_name`: Name of the registry
- `account_to_add`: The public key of the node account to add
- `domain`: Fully qualified domain name (FQDN) for the node
  - Must be 253 characters or less (compliant with DNS specification)
  - Can include subdomains (e.g., `node1.example.com`)
  - Should follow standard DNS naming conventions

#### Get node information:
```bash
./registry-client get-node <registry_name> <account_to_check>
```
Displays:
- Parent registry
- Registered account
- Domain name
- Online status
- Active status

#### Update node active status:
```bash
./registry-client update-node-active <registry_name> <authority> <account_to_update> <active>
```
- `registry_name`: Name of the registry
- `authority`: The public key of the authority node making the update
- `account_to_update`: The public key of the node to update
- `active`: Boolean value (true/false) to set the node's active status

Node active status rules:
- Any node in a registry can update any other node's active status in the same registry
- Nodes from different registries cannot update each other's active status
- The authority must be a registered node in the registry
- Both the target node and authority node must be in the same registry

#### Delete a node:
```bash
./registry-client delete-node <registry_name> <account_to_delete>
```

### Node Status Management

The registry supports two types of node status:

1. **Online Status**
   - A numeric value indicating the node's online status
   - Must be non-negative
   - Can be updated by the node itself
   - Useful for health checks and monitoring

2. **Active Status**
   - A boolean flag indicating if the node is active in the registry
   - Can be updated by any node in the same registry
   - Useful for peer-to-peer node management
   - Enables nodes to mark other nodes as active/inactive based on their observations

Example workflow:
```bash
# Register a node
./registry-client add-node my-registry Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr node1.example.com

# Update node's online status (by the node itself)
./registry-client update-node-online my-registry Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr 100

# Another node updates this node's active status
./registry-client update-node-active my-registry 5ZWj7a1f8tWkjBESHKgrLmXshuXxqeY9pW6qo8JoSXd5 Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr true

# Check node's status
./registry-client get-node my-registry Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr
```

### Node Status Best Practices

1. **Online Status**
   - Use for real-time health monitoring
   - Update frequently based on node's operational metrics
   - Consider implementing an auto-update mechanism
   - Use thresholds to determine node health

2. **Active Status**
   - Use for peer consensus on node availability
   - Implement voting or consensus mechanisms
   - Consider timeouts for inactive nodes
   - Use for load balancing and failover

3. **Status Coordination**
   - Combine both statuses for comprehensive node management
   - Use online status for immediate health
   - Use active status for long-term availability
   - Implement recovery mechanisms for inactive nodes

4. **Security Considerations**
   - Only allow authorized nodes to update statuses
   - Monitor for suspicious status changes
   - Implement rate limiting for status updates
   - Log all status changes for audit purposes

### Utility Commands

#### Check wallet balance:
```bash
./registry-client balance
```
Displays the current SOL balance of your wallet in both SOL and lamports.

#### Request SOL airdrop (devnet/testnet only):
```bash
./registry-client airdrop [amount_in_sol]
```
- `amount_in_sol`: Optional. Amount of SOL to request (default: 1 SOL)
- Note: Airdrop is only available on devnet and testnet networks

## Examples

### Client Registry Example
```bash
# Create a new registry
./registry-client create "clients"
# Output: Registry created. Transaction signature: 4Qv1...

# Add a client (valid for 30 days with limit of 1000)
./registry-client add-client clients Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr 30 1000
# Output: Client account added to registry. Transaction signature: 3tWE...

# Check client registration
./registry-client get-client clients Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr
# Output: Client registry entry information...

# Delete client registration
./registry-client delete-client clients Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr
# Output: Client account deleted from registry. Transaction signature: 2vNE...
```

### Node Registry Example
```bash
# Create a new registry
./registry-client create "nodes"
# Output: Registry created. Transaction signature: 4Qv1...

# Add a node
./registry-client add-node nodes Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr example.com
# Output: Node account added to registry. Transaction signature: 3tWE...

# Check node registration
./registry-client get-node nodes Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr
# Output: Node registry entry information...

# Delete node registration
./registry-client delete-node nodes Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr
# Output: Node account deleted from registry. Transaction signature: 2vNE...
```

## Troubleshooting

### Common Issues

1. WebSocket Connection Error
   - Ensure your validator is running and the WebSocket port (8900 for local) is accessible
   - Check that you're using the correct WebSocket URL for your network
   - For local validator, make sure both RPC (8899) and WebSocket (8900) ports are open

2. Invalid Program ID
   - Verify that the PROGRAM_ID in your .env file matches your deployed program
   - For local development, the program ID will change each time you redeploy unless you specify a keypair

3. Insufficient Funds
   - Use `balance` command to check your current SOL balance
   - Use `airdrop` command on devnet/testnet to get test SOL
   - For mainnet, you'll need to transfer real SOL to your wallet

4. Transaction Errors
   - Check that your program ID is correct
   - Verify that the account you're trying to modify exists
   - Ensure you have the necessary permissions (correct authority key)
   - Make sure you have enough SOL to pay for transaction fees

5. Domain Name Errors
   - Domain names must be 253 characters or less (compliant with DNS specification)
   - Must be a valid fully qualified domain name (FQDN)
   - Common issues:
     - Exceeding the 253-character limit
     - Using invalid characters in the domain name
     - Missing proper domain structure (e.g., missing TLD)
   - Make sure you have permission to use the domain
   - Consider using a subdomain for better organization (e.g., `node1.yourservice.com`)

### Domain Name Guidelines

The registry supports fully qualified domain names (FQDNs) with the following specifications:
- Maximum length: 253 characters (compliant with DNS standards)
- Supported formats:
  - Top-level domains (e.g., `example.com`)
  - Subdomains (e.g., `node1.example.com`)
  - Multiple levels (e.g., `node1.region1.service.example.com`)
- Best practices:
  - Use descriptive domain names for easy identification
  - Consider using subdomains for better organization
  - Follow DNS naming conventions
  - Ensure you have proper DNS configuration for the domain
  - Verify domain ownership before registration 