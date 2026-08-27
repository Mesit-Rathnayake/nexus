package blockchain

func (bc *Blockchain) IsValid() bool {
	if bc == nil || len(bc.Blocks) == 0 {
		return false
	}

	for i, block := range bc.Blocks {
		if block == nil {
			return false
		}

		// Verify the block's stored hash.
		calculatedHash := block.CalculateHash()

		if calculatedHash != block.Hash {
			return false
		}

		// Genesis block has no previous block.
		if i == 0 {
			if block.Header.PreviousHash != [32]byte{} {
				return false
			}

			continue
		}

		previousBlock := bc.Blocks[i-1]

		// Verify the cryptographic link between blocks.
		if block.Header.PreviousHash != previousBlock.Hash {
			return false
		}

		// Verify the Merkle root.
		transactionHashes := make([][32]byte, 0, len(block.Transactions))

		for _, tx := range block.Transactions {
			if tx == nil {
				return false
			}

			if tx.ID != tx.CalculateID() {
				return false
			}

			transactionHashes = append(transactionHashes, tx.ID)
		}

		expectedRoot := MerkleRoot(transactionHashes)

		if block.Header.MerkleRoot != expectedRoot {
			return false
		}
	}

	return true
}
