// Copyright 2021 The Rode Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"strings"
	"time"
)

func formatProtoTimestamp(timestamp *timestamppb.Timestamp) string {
	return timestamp.AsTime().Format(time.RFC3339Nano)
}

type policyVersionIdComponents struct {
	policyId string
	version  int
}

func parsePolicyVersionId(policyVersionId string) (id *policyVersionIdComponents, ver error) {
	parts := strings.Split(policyVersionId, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("policy version id does not match format")
	}

	policyId := parts[0]

	if _, err := uuid.ParseUUID(policyId); err != nil {
		return nil, fmt.Errorf("invalid policy id: %s", err)
	}

	version, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("policy version id does not contain a version: %s", err)
	}

	return &policyVersionIdComponents{
		policyId,
		version,
	}, nil
}
