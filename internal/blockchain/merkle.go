package blockchain

import "github.com/Mesit-Rathnayake/nexus/internal/crypto"

// MerkleRoot calculates the Merkle root of a collection of transaction hashes.
func MerkleRoot(hashes [][32]byte) [32]byte {
	if len(hashes) == 0 {
		return [32]byte{}
	}

	current := make([][32]byte, len(hashes))
	copy(current, hashes)

	for len(current) > 1 {
		next := make([][32]byte, 0, (len(current)+1)/2)

		for i := 0; i < len(current); i += 2 {
			left := current[i]

			right := left

			if i+1 < len(current) {
				right = current[i+1]
			}

			next = append(next, hashPair(left, right))
		}

		current = next
	}

	return current[0]
}

func hashPair(left, right [32]byte) [32]byte {
	data := make([]byte, 0, 64)

	data = append(data, left[:]...)
	data = append(data, right[:]...)

	return crypto.Hash(data)
}
