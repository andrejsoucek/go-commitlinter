package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFormat(t *testing.T) {
	testCases := []struct {
		Name    string
		Message string
		want    Format
		wantErr error
	}{
		{
			Name:    "happy path",
			Message: "feat(test): samples [ISSUE-1234]",
			want: Format{
				Type:    "feat",
				Scope:   "test",
				Subject: "samples",
				Task:    "ISSUE-1234",
			},
			wantErr: nil,
		},
		{
			Name:    "happy path 2",
			Message: "feat(test): samples [ISSUE-5432]",
			want: Format{
				Type:    "feat",
				Scope:   "test",
				Subject: "samples",
				Task:    "ISSUE-5432",
			},
			wantErr: nil,
		},
		{
			Name:    "happy path 3",
			Message: "fix: typo",
			want: Format{
				Type:    "fix",
				Scope:   "",
				Subject: "typo",
				Task:    "",
			},
			wantErr: nil,
		},
		{
			Name:    "whitespace",
			Message: "feat(test): samples [ISSUE-5432] ",
			want: Format{
				Type:    "feat",
				Scope:   "test",
				Subject: "samples",
				Task:    "ISSUE-5432",
			},
			wantErr: nil,
		},
		{
			Name:    "whitespace 2",
			Message: "feat(test):    samples [ISSUE-5432]",
			want: Format{
				Type:    "feat",
				Scope:   "test",
				Subject: "samples",
				Task:    "ISSUE-5432",
			},
			wantErr: nil,
		},
		{
			Name:    "whitespace 3",
			Message: "feat(test): samples      [ISSUE-5432]",
			want: Format{
				Type:    "feat",
				Scope:   "test",
				Subject: "samples",
				Task:    "ISSUE-5432",
			},
			wantErr: nil,
		},
		{
			Name:    "invalid format",
			Message: "chore(test):samples",
			wantErr: ErrFormat,
		},
		{
			Name:    "scope empty",
			Message: "docs: global",
			want: Format{
				Type:    "docs",
				Scope:   "",
				Subject: "global",
				Task:    "",
			},
			wantErr: nil,
		},
		{
			Name:    "scope empty 2",
			Message: "perf(): global",
			wantErr: ErrScope,
		},
		{
			Name:    "type empty",
			Message: "(test): test",
			wantErr: ErrFormat,
		},
		{
			Name:    "subject empty 1",
			Message: "ref(test):",
			wantErr: ErrFormat,
		},
		{
			Name:    "subject empty 2",
			Message: "refactor(test):   ",
			wantErr: ErrFormat,
		},
		{
			Name:    "subject empty 3",
			Message: "style(test):        		 ",
			wantErr: ErrFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			f, err := NewFormat(tc.Message)
			if tc.wantErr == nil {
				assert.NoError(t, err)
			}
			if err != nil {
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				return
			}

			assert.Equal(t, tc.want, f)
		})
	}
}

func TestVerify(t *testing.T) {
	testCases := []struct {
		Name    string
		Message string
		want    Format
		wantErr error
	}{
		{
			Name:    "happy path",
			Message: "feat(api): samples [PROJECT-45]",
			want: Format{
				Type:    "feat",
				Scope:   "api",
				Subject: "samples",
				Task:    "PROJECT-45",
			},
			wantErr: nil,
		},
		{
			Name:    "invalid type",
			Message: "invalid(api): test",
			wantErr: ErrType,
		},
		{
			Name:    "invalid scope",
			Message: "feat(invalid): test",
			wantErr: ErrScope,
		},
		{
			Name:    "invalid style",
			Message: "Fix(client): test",
			wantErr: ErrStyle,
		},
		{
			Name:    "invalid subject",
			Message: "feat(api): Add hoge",
			wantErr: ErrSubject,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			f, err := NewFormat(tc.Message)
			if err != nil {
				assert.NoError(t, err)
			}
			c, _ := NewConfig("rule-sample.yaml")

			err = f.Verify(c)
			if tc.wantErr == nil {
				assert.NoError(t, err)
			}
			if err != nil {
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				return
			}

			assert.Equal(t, tc.want, f)
		})
	}
}
