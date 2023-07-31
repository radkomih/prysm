package beacon

type Phase0ProduceBlockV3Response struct {
	Version                 string       `json:"version" validate:"required"`
	ExecutionPayloadBlinded bool         `json:"execution_payload_blinded" validate:"required"`
	ExeuctionPayloadValue   string       `json:"exeuction_payload_value" validate:"required"`
	Data                    *BeaconBlock `json:"data" validate:"required"`
}

type AltairProduceBlockV3Response struct {
	Version                 string             `json:"version" validate:"required"`
	ExecutionPayloadBlinded bool               `json:"execution_payload_blinded" validate:"required"`
	ExeuctionPayloadValue   string             `json:"exeuction_payload_value" validate:"required"`
	Data                    *BeaconBlockAltair `json:"data" validate:"required"`
}

type BellatrixProduceBlockV3Response struct {
	Version                 string                `json:"version" validate:"required"`
	ExecutionPayloadBlinded bool                  `json:"execution_payload_blinded" validate:"required"`
	ExeuctionPayloadValue   string                `json:"exeuction_payload_value" validate:"required"`
	Data                    *BeaconBlockBellatrix `json:"data" validate:"required"`
}

type BlindedBellatrixProduceBlockV3Response struct {
	Version                 string                       `json:"version" validate:"required"`
	ExecutionPayloadBlinded bool                         `json:"execution_payload_blinded" validate:"required"`
	ExeuctionPayloadValue   string                       `json:"exeuction_payload_value" validate:"required"`
	Data                    *BlindedBeaconBlockBellatrix `json:"data" validate:"required"`
}

type CapellaProduceBlockV3Response struct {
	Version                 string              `json:"version" validate:"required"`
	ExecutionPayloadBlinded bool                `json:"execution_payload_blinded" validate:"required"`
	ExeuctionPayloadValue   string              `json:"exeuction_payload_value" validate:"required"`
	Data                    *BeaconBlockCapella `json:"data" validate:"required"`
}

type BlindedCapellaProduceBlockV3Response struct {
	Version                 string                     `json:"version" validate:"required"`
	ExecutionPayloadBlinded bool                       `json:"execution_payload_blinded" validate:"required"`
	ExeuctionPayloadValue   string                     `json:"exeuction_payload_value" validate:"required"`
	Data                    *BlindedBeaconBlockCapella `json:"data" validate:"required"`
}

type DenebProduceBlockV3Response struct {
	Version                 string                    `json:"version" validate:"required"`
	ExecutionPayloadBlinded bool                      `json:"execution_payload_blinded" validate:"required"`
	ExeuctionPayloadValue   string                    `json:"exeuction_payload_value" validate:"required"`
	Data                    *BeaconBlockContentsDeneb `json:"data" validate:"required"`
}

type BlindedDenebProduceBlockV3Response struct {
	Version                 string                           `json:"version" validate:"required"`
	ExecutionPayloadBlinded bool                             `json:"execution_payload_blinded" validate:"required"`
	ExeuctionPayloadValue   string                           `json:"exeuction_payload_value" validate:"required"`
	Data                    *BlindedBeaconBlockContentsDeneb `json:"data" validate:"required"`
}
