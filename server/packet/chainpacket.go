package packet

// TXPacket net packet for TX
type TXPacket struct{
	AddFrom     string
	From 		string
	To 			string
	Money 		int
}

type WalletListPacket struct{
	Addresses []string
}
