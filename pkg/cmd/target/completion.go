/*
SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors

SPDX-License-Identifier: Apache-2.0
*/

package target

import (
	"strings"

	"github.com/gardener/gardenctl-v2/internal/util"
	commonTarget "github.com/gardener/gardenctl-v2/pkg/cmd/common/target"
)

func validArgsFunction(f util.Factory, o *Options, args []string, toComplete string) ([]string, error) {
	if len(args) == 0 {
		return []string{
			string(commonTarget.TargetKindGarden),
			string(commonTarget.TargetKindProject),
			string(commonTarget.TargetKindSeed),
			string(commonTarget.TargetKindShoot),
		}, nil
	}

	kind := commonTarget.TargetKind(strings.TrimSpace(args[0]))
	if err := commonTarget.ValidateKind(kind); err != nil {
		return nil, err
	}

	manager, err := f.Manager()
	if err != nil {
		return nil, err
	}

	// NB: this uses the DynamicTargetProvider from the root cmd and
	// is therefore aware of flags like --garden; the goal here is to
	// allow the user to type "gardenctl target --garden [tab][select] --project [tab][select] shoot [tab][select]"
	currentTarget, err := manager.CurrentTarget()
	if err != nil {
		return nil, err
	}

	ctx := f.Context()

	var result []string

	switch kind {
	case commonTarget.TargetKindGarden:
		result, err = util.GardenNames(manager)
	case commonTarget.TargetKindProject:
		result, err = util.ProjectNamesForTarget(ctx, manager, currentTarget)
	case commonTarget.TargetKindSeed:
		result, err = util.SeedNamesForTarget(ctx, manager, currentTarget)
	case commonTarget.TargetKindShoot:
		result, err = util.ShootNamesForTarget(ctx, manager, currentTarget)
	}

	return result, err
}
