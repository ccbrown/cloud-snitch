package store

type PrimaryIndex struct {
	HashKey  []byte `dynamodbav:"_hk"`
	RangeKey []byte `dynamodbav:"_rk"`
}

// All of the indexes are generically named to facilitate GSI overloading:
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/bp-gsi-overloading.html

type ByteByteIndex1 struct {
	HashKey  []byte `dynamodbav:"_bb1h,omitempty"`
	RangeKey []byte `dynamodbav:"_bb1r,omitempty"`
}

type ByteByteIndex2 struct {
	HashKey  []byte `dynamodbav:"_bb2h,omitempty"`
	RangeKey []byte `dynamodbav:"_bb2r,omitempty"`
}
