# UBI Benchmark

A benchmark tool for testing computing provider resources, specifically designed for Filecoin sealing operations and ZK proof generation. This tool generates benchmark tasks that can be submitted to the UBI (Universal Benchmark Infrastructure) network for testing computing provider capabilities.

## Overview

UBI Benchmark is a utility tool that helps evaluate the performance of computing providers by:

- **Benchmarking Filecoin Sealing Operations**: Tests AddPiece, PreCommit1, and PreCommit2 phases
- **Generating Commit1 (C1) Proofs**: Creates C1 outputs for ZK proof computation tasks
- **Generating Commit2 (C2) Proofs**: Computes final proofs for verification
- **Automated Task Generation**: Runs as a daemon to automatically generate and upload benchmark tasks
- **Storage Integration**: Uploads task files to MCS (Multi-Chain Storage) or Titan storage services

This tool is designed to work in conjunction with the [go-computing-provider](https://github.com/swanchain/go-computing-provider) project, which executes the generated tasks on the computing provider network.

## Ecosystem Architecture

The UBI (Universal Benchmark Infrastructure) system consists of three main components:

1.  **ubi-benchmark** (This Repository):
    *   **Role**: Task Generator & Verifier.
    *   **Function**: Generates cryptographic proof tasks (e.g., Filecoin C2), uploads inputs to storage, submits task metadata to the Hub, and verifies completed proofs.
    *   **Note**: It does *not* handle marketplace logic like matching, payments, or disputes.

2.  **ubi-engine** (The Hub):
    *   **Role**: Marketplace & Engine.
    *   **Function**: Acts as the central management service. It receives tasks from `ubi-benchmark`, manages the "Deal Market" (matching tasks to providers), handles task lifecycle (Created → Sent → Done → Verified → Rewarded), and manages payments and disputes.

3.  **go-computing-provider** (CP):
    *   **Role**: Worker / Executor.
    *   **Function**: Runs on provider hardware. It pulls tasks from `ubi-engine` (the Hub), executes the heavy computation (e.g., generating ZK proofs), and submits the results back to the Hub.

## Features

- **Sealing Benchmark**: Measure sealing performance with configurable sector sizes (512MiB, 32GiB)
- **Parallel Processing**: Support for parallel execution of sealing operations
- **GPU Support**: Optional GPU acceleration for proof generation
- **Task Generation**: Automatically generate C1 outputs for different chain heights
- **Storage Upload**: Upload generated tasks to MCS or Titan storage
- **Task Submission**: Submit tasks to UBI Hub for computing provider execution
- **Proof Verification**: Verify generated proofs for correctness

## Prerequisites

- **Go**: Version 1.22.0 or higher
- **Filecoin Parameters**: V28 parameters for Filecoin proofs (512MiB and/or 32GiB)
- **Storage**: At least 200GB free space for parameters and sector data
- **MCS Account** (optional): For uploading tasks to Multi-Chain Storage
- **Titan Account** (optional): Alternative storage backend

## Installation

### Build from Source

1. Clone the repository:
```bash
git clone https://github.com/swanchain/ubi-benchmark.git
cd ubi-benchmark
```

2. Initialize submodules:
```bash
git submodule update --init --recursive
```

3. Build the binary:
```bash
make build
```

This will create the `ubi-bench` executable in the project root.

### Download Filecoin Parameters

Before running benchmarks, you need to download the Filecoin v28 parameters:

```bash
# Set the path where parameters will be stored (at least 200GB needed)
export PARENT_PATH="<V28_PARAMS_PATH>"

# Download 512MiB parameters
curl -fsSL https://raw.githubusercontent.com/swanchain/go-computing-provider/releases/ubi/fetch-param-512.sh | bash

# Download 32GiB parameters (optional, for 32GiB sector benchmarks)
curl -fsSL https://raw.githubusercontent.com/swanchain/go-computing-provider/releases/ubi/fetch-param-32.sh | bash
```

## Configuration

Create a `config.toml` file in the project root:

```toml
[MCS]
ApiKey = "MCS_xxxxx"                    # Get from https://www.multichain.storage -> Settings -> Create API Key
BucketName = "YOUR-BUCKET-NAME"         # Get from https://www.multichain.storage -> Bucket -> Add Bucket
Network = "polygon.mainnet"             # polygon.mainnet for mainnet, polygon.mumbai for testnet

[HUB]
HUB_URL = "UBI-TASK-BASE-URL"           # UBI Hub API endpoint for task submission
TASK_URL = "TASK-STATS-URL"            # URL to check task statistics
CHECK_INTERVAL = 1                      # Check interval in minutes
BATCH_NUM = 1                           # Number of batches to generate per check
ENABLE_TITAN = 1                        # Enable Titan storage (0 = disabled, 1 = enabled)
TITAN_KEY = ""                          # Titan API key (if using Titan)
TITAN_FOLDER_512 = 607                  # Titan folder ID for 512MiB tasks
TITAN_FOLDER_32 = 608                   # Titan folder ID for 32GiB tasks
```

## Usage

### Basic Sealing Benchmark

Run a basic sealing benchmark:

```bash
./ubi-bench sealing \
  --storage-dir ~/.ubi-bench \
  --sector-size 512MiB \
  --num-sectors 1 \
  --parallel 1 \
  --miner-addr t01000
```

Options:
- `--storage-dir`: Directory for storing sectors (default: `~/.ubi-bench`)
- `--sector-size`: Sector size (e.g., `512MiB`, `32GiB`)
- `--num-sectors`: Number of sectors to seal
- `--parallel`: Number of parallel PreCommit1 operations
- `--miner-addr`: Miner address (default: `t01000`)
- `--no-gpu`: Disable GPU usage
- `--ticket-preimage`: Custom ticket random value

### Generate Commit1 (C1) Output

Generate a C1 output from a C1 input file:

```bash
./ubi-bench c1 \
  --storage-dir ~/.ubi-bench \
  --height <CHAIN_HEIGHT> \
  c1in-<miner>-<sector>.json
```

This generates a `c1out-*.json` file containing the C1 output and C2 input parameters.

### Execute Commit2 (C2) Proof Computation

Compute the final C2 proof:

```bash
./ubi-bench c2 \
  --storage-dir /var/tmp \
  c1out-<miner>-<sector>-<epoch>.json
```

Options:
- `--no-gpu`: Disable GPU for proof computation
- `--storage-dir`: Storage directory (default: `/var/tmp`)

### Verify Proof

Verify a generated proof:

```bash
./ubi-bench verify \
  --height <CHAIN_HEIGHT> \
  c2-<miner>-<sector>-<epoch>.json
```

### Batch C1 Generation

Generate multiple C1 outputs in batch:

```bash
./ubi-bench batch \
  --storage-dir ~/.ubi-bench \
  --start <START_HEIGHT> \
  --num <NUM_BATCHES> \
  c1in-<miner>-<sector>.json
```

### Upload C1 Results to MCS

Upload generated C1 outputs to MCS storage:

```bash
./ubi-bench upload \
  --c1-dir <C1_OUTPUT_DIRECTORY> \
  --type 512
```

Options:
- `--c1-dir`: Directory containing C1 output files
- `--type`: Sector type (`512` for 512MiB, `32` for 32GiB)

### Daemon Mode

Run as a daemon to automatically generate and upload tasks:

```bash
./ubi-bench daemon \
  --storage-dir ~/.ubi-bench \
  --last-height <LAST_HEIGHT> \
  --sector-type 512 \
  c1in-<miner>-<sector>.json
```

The daemon will:
1. Check task statistics from the HUB
2. Generate C1 outputs when task count is below threshold
3. Upload files to MCS or Titan storage
4. Submit tasks to the UBI Hub

Options:
- `--storage-dir`: Storage directory for sectors
- `--last-height`: Last processed chain height
- `--sector-type`: Sector type (`512` or `32`)

## Environment Variables

For GPU acceleration, set these environment variables:

```bash
export FIL_PROOFS_PARAMETER_CACHE=$PARENT_PATH
export RUST_GPU_TOOLS_CUSTOM_GPU="GeForce RTX 4090:16384"  # Adjust for your GPU
export FIL_PROOFS_USE_GPU_COLUMN_BUILDER=1
export FIL_PROOFS_USE_GPU_TREE_BUILDER=1
export FIL_PROOFS_USE_MULTICORE_SDR=1
```

For CPU-only mode:

```bash
export BELLMAN_NO_GPU=1
```

## Resource Types

The tool supports different resource types for task submission:

- **CPU512** (ID: 1): CPU for 512MiB sectors
- **CPU32G** (ID: 2): CPU for 32GiB sectors
- **GPU512** (ID: 3): GPU for 512MiB sectors
- **GPU32G** (ID: 4): GPU for 32GiB sectors

## Task Types

The ecosystem supports various task types, though `ubi-benchmark` currently focuses on Filecoin proofs.

### Implemented in ubi-benchmark
*   **Type 1: Fil-C2-512M**: Filecoin Commit2 proof for 512MiB sectors.
*   **Type 4: Fil-C2-32G**: Filecoin Commit2 proof for 32GiB sectors.

### Supported by ubi-engine (The Hub)
The engine is designed to support additional types which may be generated by other tools or future versions of this benchmark:
*   **Type 2: Aleo**: Aleo ZK proofs (defined in engine).
*   **Type 3: AI**: AI model training/inference tasks (defined in engine).

## Output Format

### Benchmark Results

When running sealing benchmarks, the output includes:

```
environment variable list:
BELLMAN_NO_GPU=1
FIL_PROOFS_USE_GPU_COLUMN_BUILDER=1
----
results (v28) SectorSize:(536870912), SectorNumber:(1)
seal: addPiece: 2.5s (214.7 MiB/s)
seal: preCommit phase 1: 45.2s (11.9 MiB/s)
seal: preCommit phase 2: 12.3s (43.6 MiB/s)
```

### JSON Output

Use `--json-out` flag for JSON formatted results:

```bash
./ubi-bench sealing --json-out --num-sectors 1
```

## Integration with Computing Provider

This tool generates benchmark tasks that are consumed by the [go-computing-provider](https://github.com/swanchain/go-computing-provider) project through the `ubi-engine` Hub:

1.  **Task Generation**: `ubi-benchmark` generates C1 outputs (which are inputs for C2 tasks) and uploads them to MCS/Titan storage.
2.  **Task Submission**: `ubi-benchmark` submits the task metadata (including storage URLs) to the `ubi-engine` (Hub) via HTTP POST.
3.  **Task Distribution**: `ubi-engine` matches the task with a suitable `go-computing-provider` based on resource requirements.
4.  **Task Execution**: The `go-computing-provider` pulls the task from the Hub, downloads inputs, executes the C2 proof computation, and submits the proof back to the Hub.
5.  **Proof Verification**: `ubi-engine` (or `ubi-benchmark` acting as a verifier) verifies the submitted proof.
6.  **Settlement**: `ubi-engine` handles the reward settlement on-chain.

## Troubleshooting

### Common Issues

1. **Missing Parameters**: Ensure Filecoin v28 parameters are downloaded and `FIL_PROOFS_PARAMETER_CACHE` is set correctly
2. **Storage Space**: Ensure sufficient disk space for sectors and parameters
3. **GPU Issues**: If GPU errors occur, try `--no-gpu` flag or check GPU driver installation
4. **MCS Upload Failures**: Verify MCS API key and bucket name in `config.toml`
5. **Network Issues**: Check HUB_URL and TASK_URL configuration

### Logs

Enable debug logging:

```bash
export GOLOG_LOG_LEVEL=debug
./ubi-bench <command>
```

## Related Projects

- [go-computing-provider](https://github.com/swanchain/go-computing-provider): Computing provider implementation for executing ZK proof tasks
- [Swan Chain](https://swanchain.io): The blockchain network for decentralized computing

## License

Apache 2.0

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Support

For questions or issues:
- Open an issue on GitHub
- Join the [Swan Chain Discord](https://discord.gg/Jd2BFSVCKw)
- Check the [Swan Chain Documentation](https://docs.swanchain.io)
