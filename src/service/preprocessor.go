package service

import (
	"encoding/hex"
	"sync"

	"github.com/fitant/xbin-api/src/types"
	"github.com/fitant/xbin-api/src/utils"
)

var encryptionKeysEphemeral chan types.EncryptionStack
var encryptionKeys chan types.EncryptionStack
var mut *sync.Mutex

func populateEncryptionStack(idSize int) {
	mut = new(sync.Mutex)
	go populateChan(idSize, encryptionKeysEphemeral)
	populateChan(idSize+1, encryptionKeys)
}

func populateChan(idSize int, c chan types.EncryptionStack) {
	for {
		mut.Lock()
		x := types.EncryptionStack{
			ID:   utils.GenerateID(idSize),
			Salt: utils.GenSalt(),
		}
		id := []byte(x.ID)
		x.Hash = hex.EncodeToString(utils.HashID(id))
		x.Key = utils.GenKey(id, x.Salt[:])
		mut.Unlock()
		c <- x
	}
}
