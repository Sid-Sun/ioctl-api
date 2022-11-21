package service

import (
	"encoding/hex"
	"sync"

	"github.com/sid-sun/ioctl-api/src/types"
	"github.com/sid-sun/ioctl-api/src/utils"
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
	ptr := -1
	stacksArr := make([]types.EncryptionStack, 25)
	stacksLen := len(stacksArr)
	genetating := false
	for {
		if ptr >= 0 {
			c <- stacksArr[ptr]
			ptr++
		}
		if !genetating && (ptr >= stacksLen-5 || ptr == -1) {
			genetating = true
			go func() {
				for i := 0; i < stacksLen; i++ {
					mut.Lock()
					x := types.EncryptionStack{
						ID:   utils.GenerateID(idSize),
						Salt: utils.GenSalt(),
					}
					id := []byte(x.ID)
					x.Hash = hex.EncodeToString(utils.HashID(id))
					x.Key = utils.GenKey(id, x.Salt[:])
					stacksArr[i] = x
					mut.Unlock()
				}
				ptr = 0
				genetating = false
			}()
		}
	}
}
