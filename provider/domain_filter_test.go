/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type domainFilterTest struct {
	domainFilter []string
	exclusions   []string
	domains      []string
	expected     bool
}

var domainFilterTests = []domainFilterTest{
	{
		[]string{"google.com.", "exaring.de", "inovex.de"},
		[]string{},
		[]string{"google.com", "exaring.de", "inovex.de"},
		true,
	},
	{
		[]string{"google.com.", "exaring.de", "inovex.de"},
		[]string{},
		[]string{"google.com", "exaring.de", "inovex.de"},
		true,
	},
	{
		[]string{"google.com.", "exaring.de.", "inovex.de"},
		[]string{},
		[]string{"google.com", "exaring.de", "inovex.de"},
		true,
	},
	{
		[]string{"foo.org.      "},
		[]string{},
		[]string{"foo.org"},
		true,
	},
	{
		[]string{"   foo.org"},
		[]string{},
		[]string{"foo.org"},
		true,
	},
	{
		[]string{"foo.org."},
		[]string{},
		[]string{"foo.org"},
		true,
	},
	{
		[]string{"foo.org."},
		[]string{},
		[]string{"baz.org"},
		false,
	},
	{
		[]string{"baz.foo.org."},
		[]string{},
		[]string{"foo.org"},
		false,
	},
	{
		[]string{"", "foo.org."},
		[]string{},
		[]string{"foo.org"},
		true,
	},
	{
		[]string{"", "foo.org."},
		[]string{},
		[]string{},
		true,
	},
	{
		[]string{""},
		[]string{},
		[]string{"foo.org"},
		true,
	},
	{
		[]string{""},
		[]string{},
		[]string{},
		true,
	},
	{
		[]string{" "},
		[]string{},
		[]string{},
		true,
	},
	{
		[]string{"bar.sub.example.org"},
		[]string{},
		[]string{"foo.bar.sub.example.org"},
		true,
	},
	{
		[]string{"example.org"},
		[]string{},
		[]string{"anexample.org", "test.anexample.org"},
		false,
	},
	{
		[]string{".example.org"},
		[]string{},
		[]string{"anexample.org", "test.anexample.org"},
		false,
	},
	{
		[]string{".example.org"},
		[]string{},
		[]string{"example.org"},
		false,
	},
	{
		[]string{".example.org"},
		[]string{},
		[]string{"test.example.org"},
		true,
	},
	{
		[]string{"anexample.org"},
		[]string{},
		[]string{"example.org", "test.example.org"},
		false,
	},
	{
		[]string{".org"},
		[]string{},
		[]string{"example.org", "test.example.org", "foo.test.example.org"},
		true,
	},
	{
		[]string{"example.org"},
		[]string{"api.example.org"},
		[]string{"example.org", "test.example.org", "foo.test.example.org"},
		true,
	},
	{
		[]string{"example.org"},
		[]string{"api.example.org"},
		[]string{"foo.api.example.org", "api.example.org"},
		false,
	},
	{
		[]string{"   example.org. "},
		[]string{"   .api.example.org    "},
		[]string{"foo.api.example.org", "bar.baz.api.example.org."},
		false,
	},
	{
		[]string{"example.org."},
		[]string{"api.example.org"},
		[]string{"dev-api.example.org", "qa-api.example.org"},
		true,
	},
	{
		[]string{"example.org."},
		[]string{"api.example.org"},
		[]string{"dev.api.example.org", "qa.api.example.org"},
		false,
	},
	{
		[]string{"example.org", "api.example.org"},
		[]string{"internal.api.example.org"},
		[]string{"foo.api.example.org"},
		true,
	},
	{
		[]string{"example.org", "api.example.org"},
		[]string{"internal.api.example.org"},
		[]string{"foo.internal.api.example.org"},
		false,
	},
}

func TestDomainFilterMatch(t *testing.T) {
	for i, tt := range domainFilterTests {
		if len(tt.exclusions) > 0 {
			t.Skip("NewDomainFilter() doesn't support exclusions")
		}
		domainFilter := NewDomainFilter(tt.domainFilter)
		for _, domain := range tt.domains {
			assert.Equal(t, tt.expected, domainFilter.Match(domain), "should not fail: %v in test-case #%v", domain, i)
			assert.Equal(t, tt.expected, domainFilter.Match(domain+"."), "should not fail: %v in test-case #%v", domain+".", i)
		}
	}
}

func TestDomainFilterWithExclusions(t *testing.T) {
	for i, tt := range domainFilterTests {
		domainFilter := NewDomainFilterWithExclusions(tt.domainFilter, tt.exclusions)
		for _, domain := range tt.domains {
			assert.Equal(t, tt.expected, domainFilter.Match(domain), "should not fail: %v in test-case #%v", domain, i)
			assert.Equal(t, tt.expected, domainFilter.Match(domain+"."), "should not fail: %v in test-case #%v", domain+".", i)
		}
	}
}

func TestDomainFilterMatchWithEmptyFilter(t *testing.T) {
	for _, tt := range domainFilterTests {
		domainFilter := DomainFilter{}
		for i, domain := range tt.domains {
			assert.True(t, domainFilter.Match(domain), "should not fail: %v in test-case #%v", domain, i)
			assert.True(t, domainFilter.Match(domain+"."), "should not fail: %v in test-case #%v", domain+".", i)
		}
	}
}
