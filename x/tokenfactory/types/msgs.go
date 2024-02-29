package types

import (
	"cosmossdk.io/core/address"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// constants
const (
	TypeMsgCreateDenom       = "create_denom"
	TypeMsgMint              = "tf_mint"
	TypeMsgBurn              = "tf_burn"
	TypeMsgForceTransfer     = "force_transfer"
	TypeMsgChangeAdmin       = "change_admin"
	TypeMsgSetDenomMetadata  = "set_denom_metadata"
	TypeMsgSetBeforeSendHook = "set_before_send_hook"
)

var (
	_ sdk.Msg = &MsgCreateDenom{}
	_ sdk.Msg = &MsgMint{}
	_ sdk.Msg = &MsgBurn{}
	_ sdk.Msg = &MsgForceTransfer{}
	_ sdk.Msg = &MsgChangeAdmin{}
	_ sdk.Msg = &MsgSetDenomMetadata{}
	_ sdk.Msg = &MsgSetBeforeSendHook{}
)

// NewMsgCreateDenom creates a msg to create a new denom
func NewMsgCreateDenom(sender, subdenom string) *MsgCreateDenom {
	return &MsgCreateDenom{
		Sender:   sender,
		Subdenom: subdenom,
	}
}

func (m MsgCreateDenom) Validate(accAddrCodec address.Codec) error {
	if addr, err := accAddrCodec.StringToBytes(m.Sender); err != nil {
		return err
	} else if len(addr) == 0 {
		return ErrEmptySender
	}

	if _, err := GetTokenDenom(m.Sender, m.Subdenom); err != nil {
		return ErrInvalidDenom
	}

	return nil
}

// NewMsgMint creates a message to mint tokens
func NewMsgMint(sender string, amount sdk.Coin) *MsgMint {
	return &MsgMint{
		Sender: sender,
		Amount: amount,
	}
}

func NewMsgMintTo(sender string, amount sdk.Coin, mintToAddress string) *MsgMint {
	return &MsgMint{
		Sender:        sender,
		Amount:        amount,
		MintToAddress: mintToAddress,
	}
}

func (m MsgMint) Validate(accAddrCodec address.Codec) error {
	if addr, err := accAddrCodec.StringToBytes(m.Sender); err != nil {
		return err
	} else if len(addr) == 0 {
		return ErrEmptySender
	}

	if addr, err := accAddrCodec.StringToBytes(m.MintToAddress); err != nil {
		return err
	} else if len(addr) == 0 {
		return ErrEmptyMintToAddress
	}

	if !m.Amount.IsValid() || m.Amount.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	return nil
}

// NewMsgBurn creates a message to burn tokens
func NewMsgBurn(sender string, amount sdk.Coin) *MsgBurn {
	return &MsgBurn{
		Sender: sender,
		Amount: amount,
	}
}

// NewMsgBurn creates a message to burn tokens
func NewMsgBurnFrom(sender string, amount sdk.Coin, burnFromAddress string) *MsgBurn {
	return &MsgBurn{
		Sender:          sender,
		Amount:          amount,
		BurnFromAddress: burnFromAddress,
	}
}

func (m MsgBurn) Validate(accAddrCodec address.Codec) error {
	if addr, err := accAddrCodec.StringToBytes(m.Sender); err != nil {
		return err
	} else if len(addr) == 0 {
		return ErrEmptySender
	}

	if !m.Amount.IsValid() || m.Amount.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	return nil
}

// NewMsgForceTransfer creates a transfer funds from one account to another
func NewMsgForceTransfer(sender string, amount sdk.Coin, fromAddr, toAddr string) *MsgForceTransfer {
	return &MsgForceTransfer{
		Sender:              sender,
		Amount:              amount,
		TransferFromAddress: fromAddr,
		TransferToAddress:   toAddr,
	}
}

func (m MsgForceTransfer) Validate(accAddrCodec address.Codec) error {
	if addr, err := accAddrCodec.StringToBytes(m.Sender); err != nil {
		return err
	} else if len(addr) == 0 {
		return ErrEmptySender
	}

	if addr, err := accAddrCodec.StringToBytes(m.TransferFromAddress); err != nil {
		return err
	} else if len(addr) == 0 {
		return ErrEmptyTransferFromAddress
	}

	if addr, err := accAddrCodec.StringToBytes(m.TransferToAddress); err != nil {
		return err
	} else if len(addr) == 0 {
		return ErrEmptyTransferToAddress
	}

	if !m.Amount.IsValid() || m.Amount.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	return nil
}

// NewMsgChangeAdmin creates a message to burn tokens
func NewMsgChangeAdmin(sender, denom, newAdmin string) *MsgChangeAdmin {
	return &MsgChangeAdmin{
		Sender:   sender,
		Denom:    denom,
		NewAdmin: newAdmin,
	}
}

func (m MsgChangeAdmin) Validate(accAddrCodec address.Codec) error {
	if addr, err := accAddrCodec.StringToBytes(m.Sender); err != nil {
		return err
	} else if len(addr) == 0 {
		return ErrEmptySender
	}

	if addr, err := accAddrCodec.StringToBytes(m.NewAdmin); err != nil {
		return err
	} else if len(addr) == 0 {
		return ErrEmptyNewAdmin
	}

	if _, _, err := DeconstructDenom(accAddrCodec, m.Denom); err != nil {
		return err
	}

	return nil
}

// NewMsgChangeAdmin creates a message to burn tokens
func NewMsgSetDenomMetadata(sender string, metadata banktypes.Metadata) *MsgSetDenomMetadata {
	return &MsgSetDenomMetadata{
		Sender:   sender,
		Metadata: metadata,
	}
}

func (m MsgSetDenomMetadata) Validate(accAddrCodec address.Codec) error {
	if addr, err := accAddrCodec.StringToBytes(m.Sender); err != nil {
		return err
	} else if len(addr) == 0 {
		return ErrEmptySender
	}

	if err := m.Metadata.Validate(); err != nil {
		return err
	}

	if _, _, err := DeconstructDenom(accAddrCodec, m.Metadata.Base); err != nil {
		return err
	}
	return nil
}

// NewMsgSetBeforeSendHook creates a message to set a new before send hook
func NewMsgSetBeforeSendHook(sender string, denom string, cosmwasmAddress string) *MsgSetBeforeSendHook {
	return &MsgSetBeforeSendHook{
		Sender:          sender,
		Denom:           denom,
		CosmwasmAddress: cosmwasmAddress,
	}
}

func (m MsgSetBeforeSendHook) Validate(accAddrCodec address.Codec) error {
	if addr, err := accAddrCodec.StringToBytes(m.Sender); err != nil {
		return err
	} else if len(addr) == 0 {
		return ErrEmptySender
	}

	if addr, err := accAddrCodec.StringToBytes(m.CosmwasmAddress); err != nil {
		return err
	} else if len(addr) == 0 {
		return ErrEmptySender
	}

	if _, _, err := DeconstructDenom(accAddrCodec, m.Denom); err != nil {
		return ErrInvalidDenom
	}
	return nil
}
