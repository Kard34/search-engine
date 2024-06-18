package iq_wordbreak

type CharFlag int

const (
	
	UpUpperLevel int= 0 
	UpperLevel  int=1  
	MiddleLevel int= 2
	LowerLevel int= 3
	LevelMask int= 3
	// 'Type
	Consonant int= 4
	LeadVowel int= 8
	FollowVowel int= 16
	UpperVowel int= 32
	LowerVowel int= 64
	Tone int= 128
	Special int= 256
	Combine int= 512
	Stone int= 1024
	ThaiNumber int= 2048
	Dot int= 4096
	English int= 8192

	Alien int= 65536

	RearVowel int =   FollowVowel | UpperVowel | LowerVowel
	Vowel int= LeadVowel | RearVowel
	UnLeadable int= Tone | Special | RearVowel
	ThaiAlpha int= Consonant | Vowel | Tone | Special
	BreakAlpha int= ThaiAlpha | Dot
	Alpha int= BreakAlpha | English
)
