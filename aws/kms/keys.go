package kms

import (
	"fmt"
	"strings"

	awskms "github.com/aws/aws-sdk-go/service/kms"
	"github.com/genevieve/leftovers/aws/common"
)

type keysClient interface {
	ListKeys(*awskms.ListKeysInput) (*awskms.ListKeysOutput, error)
	DescribeKey(*awskms.DescribeKeyInput) (*awskms.DescribeKeyOutput, error)
	ListResourceTags(*awskms.ListResourceTagsInput) (*awskms.ListResourceTagsOutput, error)
	DisableKey(*awskms.DisableKeyInput) (*awskms.DisableKeyOutput, error)
	ScheduleKeyDeletion(*awskms.ScheduleKeyDeletionInput) (*awskms.ScheduleKeyDeletionOutput, error)
}

type Keys struct {
	client keysClient
	logger logger
}

func NewKeys(client keysClient, logger logger) Keys {
	return Keys{
		client: client,
		logger: logger,
	}
}

func (k Keys) ListOnly(filter string) ([]common.Deletable, error) {
	return k.get(filter)
}

func (k Keys) List(filter string) ([]common.Deletable, error) {
	resources, err := k.get(filter)
	if err != nil {
		return nil, err
	}

	var delete []common.Deletable
	for _, r := range resources {
		proceed := k.logger.PromptWithDetails(r.Type(), r.Name())
		if !proceed {
			continue
		}

		delete = append(delete, r)
	}

	return delete, nil
}

func (k Keys) get(filter string) ([]common.Deletable, error) {
	keys, err := k.client.ListKeys(&awskms.ListKeysInput{})
	if err != nil {
		return nil, fmt.Errorf("Listing KMS Keys: %s", err)
	}

	var resources []common.Deletable
	for _, key := range keys.Keys {
		metadata, err := k.client.DescribeKey(&awskms.DescribeKeyInput{KeyId: key.KeyId})
		if err != nil {
			k.logger.Printf("ERROR describing key %s: %s", *key.KeyId, err)
		}

		if alreadyDeleted(metadata.KeyMetadata) {
			continue
		}

		tags, err := k.client.ListResourceTags(&awskms.ListResourceTagsInput{KeyId: key.KeyId})
		if err != nil {
			k.logger.Printf("ERROR listing tags for key %s: %s", *key.KeyId, err)
		}

		resource := NewKey(k.client, key.KeyId, metadata.KeyMetadata, tags.Tags)

		if !strings.Contains(resource.Name(), filter) {
			continue
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

func alreadyDeleted(metadata *awskms.KeyMetadata) bool {
	if metadata == nil ||
		*metadata.KeyState == awskms.KeyStatePendingDeletion ||
		*metadata.KeyState == awskms.KeyStateDisabled {
		return true
	}

	return false
}
