package shog

type Key int

const (
	Null Key = iota
	StartOfHeading
	StartOfText
	EndOfText
	EndOfTransmission
	Enquiry
	Acknoledge
	Bell
	Backspace
	CharacterTabulation
	LineFeed
	LineTabulation
	FormFeed
	CarriageReturn
	ShiftOut
	ShiftIn
	DataLinkEscape
	DeviceControlOne
	DeviceControlTwo
	DeviceControlTree
	DeviceControlFour
	NegativeAcknowledge
	SynchronousIdle
	EndOfTransmissionBlock
	Cancel
	EndOfMedium
	Subtitute
	Escape
	InformamtionSeparatorFour
	InformamtionSeparatorThree
	InformamtionSeparatorTwo
	InformamtionSeparatorOne
)
