package tcp

import (
	"testing"
)

func TestNewPacketHead(t *testing.T) {
	// bodylen := rand.Uint32()
	// p := &PackFrame{}
	// p.set(1)
	// p.SetSubType(PacketTypAsync)
	// p.SetBodyLength(bodylen)
	// if p.Head.GetType() != 1 {
	// 	t.Error("GetType failed")
	// }
	// if p.Head.GetSubType() != PacketTypAsync {
	// 	t.Error("GetSubType failed")
	// }
	// if p.Head.GetBodyLength() != bodylen {
	// 	t.Error("GetBodyLength failed")
	// }
}

func TestPacketHeadReset(t *testing.T) {
	h := PacketHead{}
	hlen := len(h)

	for i := 0; i < hlen; i++ {
		h[i] = 1
	}
	h.Reset()

	for i := 0; i < hlen; i++ {
		if h[i] != 0 {
			t.Error("Reset failed")
		}
	}
}
