package emoji

import (
	"github.com/brahma-adshonor/gohook"
	"syscall"
)

func emojiHook(fd int, p []byte) (n int, err error) {
	pPtr := &p
	pLen := len(p)
	var emojiBytes []byte
	lock := false
	var lockBytes []byte
	for i := 0; i < pLen; i++ {
		if p[i] != 58 && !lock {
			emojiBytes = append(emojiBytes, p[i])
			continue
		}
		if p[i] == 58 && !lock {
			lock = true
			lockBytes = append(lockBytes, p[i])
			continue
		}
		if p[i] == 58 && lock {
			if len(lockBytes) == 1 {
				emojiBytes = append(emojiBytes, lockBytes...)
				continue
			}
			lock = false
			lockBytes = append(lockBytes, p[i])
			emojiStr := emojiMap[string(lockBytes)]
			if emojiStr != "" {
				emojiBytes = append(emojiBytes, []byte(emojiStr)...)
			} else {
				emojiBytes = append(emojiBytes, lockBytes...)
			}
			lockBytes = []byte{}
			continue
		}
		lockBytes = append(lockBytes, p[i])
		if p[i] == 32 {
			lock = false
			emojiBytes = append(emojiBytes, lockBytes...)
			lockBytes = []byte{}
		}
	}
	*pPtr = emojiBytes
	_, err = emojiHookTramp(fd, p)
	return pLen, err
}

func emojiHookTramp(fd int, p []byte) (n int, err error) {
	return 0, nil
}

func Activity() {
	gohook.Hook(syscall.Write, emojiHook, emojiHookTramp)
}

func InActivity() {
	gohook.UnHook(syscall.Write)
}
