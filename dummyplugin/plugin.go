package dummyplugin

import (
	"context"

	"github.com/notaryproject/notation-go/plugin/proto"
)

type DummyPlugin struct {
}

// NewDummyPlugin validate the metadata of the plugin and return a *CLIPlugin.
func NewDummyPlugin() (*DummyPlugin, error) {
	return &DummyPlugin{}, nil
}

// DescribeKey returns the KeySpec of a key.
func (p *DummyPlugin) DescribeKey(ctx context.Context, req *proto.DescribeKeyRequest) (*proto.DescribeKeyResponse, error) {
	return &proto.DescribeKeyResponse{}, nil
}

// GenerateSignature generates the raw signature based on the request.
func (p *DummyPlugin) GenerateSignature(ctx context.Context, req *proto.GenerateSignatureRequest) (*proto.GenerateSignatureResponse, error) {
	return &proto.GenerateSignatureResponse{}, nil
}

// GenerateEnvelope generates the Envelope with signature based on the
// request.
func (p *DummyPlugin) GenerateEnvelope(ctx context.Context, req *proto.GenerateEnvelopeRequest) (*proto.GenerateEnvelopeResponse, error) {
	return &proto.GenerateEnvelopeResponse{}, nil
}

func (p *DummyPlugin) VerifySignature(ctx context.Context, req *proto.VerifySignatureRequest) (*proto.VerifySignatureResponse, error) {
	return &proto.VerifySignatureResponse{
		ProcessedAttributes: []interface{}{},
		VerificationResults: map[proto.Capability]*proto.VerificationResult{},
	}, nil
}

func (p *DummyPlugin) GetMetadata(ctx context.Context, req *proto.GetMetadataRequest) (*proto.GetMetadataResponse, error) {
	return &proto.GetMetadataResponse{
		Name:         "vola! we have new plugin",
		Description:  "Do I really need a description :O",
		URL:          "Reach out to me via url",
		Version:      "1.0.0",
		Capabilities: []proto.Capability{proto.CapabilitySignatureGenerator, proto.CapabilityTrustedIdentityVerifier, proto.CapabilityRevocationCheckVerifier},
	}, nil
}
