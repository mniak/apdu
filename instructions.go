package apdu

type Instruction byte

// Table 4.1 — Commands in the alphabetic order
const (
	Instruction44_ActivateFile                   Instruction = 0x44
	InstructionE2_AppendRecord                   Instruction = 0xE2
	Instruction24_ChangeReferenceData            Instruction = 0x24
	InstructionE0_CreateFile                     Instruction = 0xE0
	Instruction04_DeactivateFile                 Instruction = 0x04
	InstructionE4_DeleteFile                     Instruction = 0xE4
	Instruction26_DisableVerificationRequirement Instruction = 0x26
	Instruction28_EnableVerificationRequirement  Instruction = 0x28
	InstructionC2_Envelope                       Instruction = 0xC2
	InstructionC3_Envelope                       Instruction = 0xC3
	Instruction0E_EraseBinary                    Instruction = 0x0E
	Instruction0F_EraseBinary                    Instruction = 0x0F
	Instruction0C_EraseRecords                   Instruction = 0x0C
	Instruction82_ExternalMutualAuthenticate     Instruction = 0x82
	Instruction86_GeneralAuthenticate            Instruction = 0x86
	Instruction87_GeneralAuthenticate            Instruction = 0x87
	Instruction46_GenerateAsymmetricKeyPair      Instruction = 0x46
	Instruction84_GetChallenge                   Instruction = 0x84
	InstructionCA_GetData                        Instruction = 0xCA
	InstructionCB_GetData                        Instruction = 0xCB
	InstructionC0_GetResponse                    Instruction = 0xC0
	Instruction88_InternalAuthenticate           Instruction = 0x88
	Instruction70_ManageChannel                  Instruction = 0x70
	Instruction22_ManageSecurityEnvironment      Instruction = 0x22
	Instruction10_PerformScqlOperation           Instruction = 0x10
	Instruction2A_PerformSecurityOperation       Instruction = 0x2A
	Instruction12_PerformTransactionOperation    Instruction = 0x12
	Instruction14_PerformUserOperation           Instruction = 0x14
	InstructionDA_PutData                        Instruction = 0xDA
	InstructionDB_PutData                        Instruction = 0xDB
	InstructionB0_ReadBinary                     Instruction = 0xB0
	InstructionB1_ReadBinary                     Instruction = 0xB1
	InstructionB2_ReadRecords                    Instruction = 0xB2
	InstructionB3_ReadRecords                    Instruction = 0xB3
	Instruction2C_ResetRetryCounter              Instruction = 0x2C
	InstructionA0_SearchBinary                   Instruction = 0xA0
	InstructionA1_SearchBinary                   Instruction = 0xA1
	InstructionA2_SearchRecord                   Instruction = 0xA2
	InstructionA4_Select                         Instruction = 0xA4
	InstructionFE_TerminateCardUsage             Instruction = 0xFE
	InstructionE6_TerminateDf                    Instruction = 0xE6
	InstructionE8_TerminateEf                    Instruction = 0xE8
	InstructionD6_UpdateBinary                   Instruction = 0xD6
	InstructionD7_UpdateBinary                   Instruction = 0xD7
	InstructionDC_UpdateRecord                   Instruction = 0xDC
	InstructionDD_UpdateRecord                   Instruction = 0xDD
	Instruction20_Verify                         Instruction = 0x20
	Instruction21_Verify                         Instruction = 0x21
	InstructionD0_WriteBinary                    Instruction = 0xD0
	InstructionD1_WriteBinary                    Instruction = 0xD1
	InstructionD2_WriteRecord                    Instruction = 0xD2
)
