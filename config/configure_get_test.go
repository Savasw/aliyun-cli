/*
 * Copyright (C) 1999-2019 Alibaba Group Holding Limited
 */
package config

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aliyun/aliyun-cli/cli"
)

func TestDoConfigureGet(t *testing.T) {
	w := new(bytes.Buffer)
	ctx := cli.NewCommandContext(w)
	AddFlags(ctx.Flags())
	originhook := hookLoadConfiguration
	defer func() {
		hookLoadConfiguration = originhook
	}()
	//testcase 1
	hookLoadConfiguration = func(fn func(w io.Writer) (Configuration, error)) func(w io.Writer) (Configuration, error) {
		return func(w io.Writer) (Configuration, error) {
			return Configuration{}, errors.New("error")
		}
	}
	doConfigureGet(ctx, []string{})
	assert.Equal(t, "\x1b[1;31mload configuration failed error\x1b[0m\n", w.String())

	//testcase 2
	hookLoadConfiguration = func(fn func(w io.Writer) (Configuration, error)) func(w io.Writer) (Configuration, error) {
		return func(w io.Writer) (Configuration, error) {
			return Configuration{CurrentProfile: "default", Profiles: []Profile{Profile{Name: "default", Mode: AK, AccessKeyId: "default_aliyun_access_key_id", AccessKeySecret: "default_aliyun_access_key_secret", OutputFormat: "json"}, Profile{Name: "aaa", Mode: AK, AccessKeyId: "sdf", AccessKeySecret: "ddf", OutputFormat: "json"}}}, nil
		}
	}
	w.Reset()
	ctx.Flags().Flags()[1].SetAssigned(true)
	doConfigureGet(ctx, []string{})
	assert.Equal(t, "\x1b[1;31mprofile  not found!\x1b[0m\n", w.String())

	//testcase 3

	hookLoadConfiguration = func(fn func(w io.Writer) (Configuration, error)) func(w io.Writer) (Configuration, error) {
		return func(w io.Writer) (Configuration, error) {
			return Configuration{CurrentProfile: "default", Profiles: []Profile{Profile{Name: "default", Mode: AK, AccessKeyId: "default_aliyun_access_key_id", AccessKeySecret: "default_aliyun_access_key_secret", OutputFormat: "json"}, Profile{Name: "aaa", Mode: AK, AccessKeyId: "sdf", AccessKeySecret: "ddf", OutputFormat: "json"}}}, nil
		}
	}
	w.Reset()
	ctx.Flags().Flags()[1].SetAssigned(false)
	doConfigureGet(ctx, []string{"profile", "mode", "access-key-id", "access-key-secret", "sts-token", "ram-role-name", "ram-role-arn", "role-session-name", "private-key", "key-pair-name", "region", "language"})
	assert.Equal(t, "profile=default\nmode=AK\naccess-key-id=*************************_id\naccess-key-secret=*****************************ret\nsts-token=\nram-role-name=\nram-role-arn=\nrole-session-name=\nprivate-key=\nkey-pair-name=\nlanguage=\n\n", w.String())
}
