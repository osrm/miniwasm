package wasm_hooks

import (
	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	ibchooks "github.com/initia-labs/initia/x/ibc-hooks"
)

var (
	_ ibchooks.OnRecvPacketOverrideHooks            = WasmHooks{}
	_ ibchooks.OnAcknowledgementPacketOverrideHooks = WasmHooks{}
	_ ibchooks.OnTimeoutPacketOverrideHooks         = WasmHooks{}
)

type WasmHooks struct {
	wasmKeeper *wasmkeeper.Keeper
	ac         address.Codec
}

func NewWasmHooks(wasmKeeper *wasmkeeper.Keeper, ac address.Codec) *WasmHooks {
	return &WasmHooks{
		wasmKeeper: wasmKeeper,
		ac:         ac,
	}
}

func (h WasmHooks) OnRecvPacketOverride(im ibchooks.IBCMiddleware, ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	if isIcs20, ics20Data := isIcs20Packet(packet.GetData()); isIcs20 {
		return h.onRecvIcs20Packet(ctx, im, packet, relayer, ics20Data)
	}

	if isIcs721, ics721Data := isIcs721Packet(packet.GetData(), packet.SourcePort); isIcs721 {
		return h.onRecvIcs721Packet(ctx, im, packet, relayer, ics721Data)
	}

	return im.App.OnRecvPacket(ctx, packet, relayer)
}

func (h WasmHooks) OnAcknowledgementPacketOverride(im ibchooks.IBCMiddleware, ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	if isIcs20, ics20Data := isIcs20Packet(packet.GetData()); isIcs20 {
		return h.onAckIcs20Packet(ctx, im, packet, acknowledgement, relayer, ics20Data)
	}

	if isIcs721, ics721Data := isIcs721Packet(packet.GetData(), packet.DestinationPort); isIcs721 {
		return h.onAckIcs721Packet(ctx, im, packet, acknowledgement, relayer, ics721Data)
	}

	return im.App.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

func (h WasmHooks) OnTimeoutPacketOverride(im ibchooks.IBCMiddleware, ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	if isIcs20, ics20Data := isIcs20Packet(packet.GetData()); isIcs20 {
		return h.onTimeoutIcs20Packet(ctx, im, packet, relayer, ics20Data)
	}

	if isIcs721, ics721Data := isIcs721Packet(packet.GetData(), packet.DestinationPort); isIcs721 {
		return h.onTimeoutIcs721Packet(ctx, im, packet, relayer, ics721Data)
	}

	return im.App.OnTimeoutPacket(ctx, packet, relayer)
}

func (h WasmHooks) checkACL(im ibchooks.IBCMiddleware, ctx sdk.Context, addrStr string) (bool, error) {
	addr, err := h.ac.StringToBytes(addrStr)
	if err != nil {
		return false, err
	}

	return im.HooksKeeper.GetAllowed(ctx, addr)
}
