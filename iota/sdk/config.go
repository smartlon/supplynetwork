package sdk

const endpoint = "http://202.117.43.212:14265"
const zmqAddress="tcp://202.117.43.212:5556"
// difficulty of the proof of work required to attach a transaction on the tangle
const mwm = 9

// how many milestones back to start the random walk from
const depth = 3

const chainTag = "HYPERLEDGER"

const MamMode = "public" // "restricted" "private"

const MamSideKey = ""

// IOTA Wallet
const DefaultWalletSeed = "RTZKOKTX9WMASJMXG9SGSWNGSAE9TWHACCTQNVLVR9XSDPBMZGVODEUZU9USLLKZAIOZGLSA9UBOTG9LQ"
const DefaultWalletKeyIndex = 4
const DefaultAmount = 100